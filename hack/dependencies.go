// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/blang/semver"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/yaml"
)

var (
	cfgFile      string
	dependencies []*dependency
	client       *http.Client
)

type dependency struct {
	Name           string                        `json:"name"`
	Version        string                        `json:"version"`
	Pattern        string                        `json:"pattern"`
	AutoUpdate     *autoUpdate                   `json:"autoupdate,omitempty"`
	Dev            bool                          `json:"dev"`
	Checksums      map[string]map[string]*string `json:"checksums"`
	TarballSubpath *string                       `json:"tarballSubpath,omitempty"`
}

type autoUpdate struct {
	Github    string `json:"github"`
	Checksums struct {
		File         *string `json:"file,omitempty"`
		ReleaseNotes *bool   `json:"releaseNotes,omitempty"`
	} `json:"checksums,omitempty"`
}

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	Body string `json:"body"`
}

type fields struct {
	Name,
	Version,
	OS,
	Arch string
}

type installCommand struct {
	*cobra.Command
	arch, os, destDir string
	dev               bool
}

func newInstallCommand() *cobra.Command {
	c := &installCommand{}
	command := &cobra.Command{
		Use:   "install",
		Short: "install dependencies",
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.Install()
		},
	}
	command.Flags().StringVar(&c.os, "os", runtime.GOOS, "os to install dependencies for")
	command.Flags().StringVar(&c.arch, "arch", runtime.GOARCH, "arch to install dependencies for")
	command.Flags().StringVarP(&c.destDir, "destination", "d", "", "directory to which binaries will be installed (required)")
	command.Flags().BoolVar(&c.dev, "dev", false, "only install dev dependencies (default false)")
	command.MarkFlagRequired("destination")
	return command
}

func (c *installCommand) Install() error {
	g := errgroup.Group{}
	for _, dep := range dependencies {
		if c.dev && !dep.Dev {
			continue
		}
		dep := dep // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			if err := c.downloadAndVerify(dep); err != nil {
				return fmt.Errorf("downloading %s: %w", dep.Name, err)
			}
			return nil
		})
	}
	return g.Wait()
}

func (c *installCommand) downloadAndVerify(dep *dependency) error {
	checksum := dep.Checksums[c.os][c.arch]
	if checksum == nil {
		return fmt.Errorf("%s not supported on platform %s/%s", dep.Name, c.os, c.arch)
	}

	urlTpl, err := template.New(dep.Name).Parse(dep.Pattern)
	if err != nil {
		return err
	}
	var url bytes.Buffer
	fields := fields{Name: dep.Name, Version: dep.Version, OS: c.os, Arch: c.arch}
	if err := urlTpl.Execute(&url, fields); err != nil {
		return err
	}
	resp, err := client.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code retrieving url: %s: %s", url.String(), resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("%s downloaded", dep.Name)

	hasher := sha256.New()
	_, err = hasher.Write(content)
	if err != nil {
		return err
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if actual != *checksum {
		return fmt.Errorf("%s: wrong checksum, expected: %s, got: %s", dep.Name, *checksum, actual)
	}
	log.Printf("%s validated", dep.Name)

	if dep.TarballSubpath != nil {
		tarballTpl, err := template.New(dep.Name).Parse(*dep.TarballSubpath)
		if err != nil {
			return err
		}
		var tarballSubpath bytes.Buffer
		if err := tarballTpl.Execute(&tarballSubpath, fields); err != nil {
			return err
		}
		gzReader, err := gzip.NewReader(bytes.NewReader(content))
		if err != nil {
			return err
		}
		tgzReader := tar.NewReader(gzReader)
		for {
			hdr, err := tgzReader.Next()
			if err == io.EOF {
				return fmt.Errorf("tarball subpath not found: %s", tarballSubpath.String())
			}
			if hdr.Name == tarballSubpath.String() {
				out, err := os.Create(path.Join(c.destDir, dep.Name))
				if err != nil {
					return err
				}
				_, err = io.Copy(out, tgzReader)
				if err != nil {
					return err
				}
				if err := out.Chmod(0777); err != nil {
					return err
				}
				log.Printf("%s extracted", dep.Name)
				break
			}
		}
	} else {
		if err = os.WriteFile(path.Join(c.destDir, dep.Name), content, 0777); err != nil {
			return err
		}
	}
	return nil
}

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update dependencies",
		PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
			depsYAML, err := yaml.Marshal(&dependencies)
			if err != nil {
				return err
			}
			err = os.WriteFile(cfgFile, depsYAML, 0644)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			g := errgroup.Group{}
			for _, dep := range dependencies {
				if dep.AutoUpdate == nil {
					continue
				}
				dep := dep // https://golang.org/doc/faq#closures_and_goroutines
				g.Go(func() error {
					if err := update(dep); err != nil {
						return fmt.Errorf("updating %s: %w", dep.Name, err)
					}
					return nil
				})
			}
			return g.Wait()
		},
	}
}

func update(dep *dependency) error {
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", dep.AutoUpdate.Github))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code from server: %s", resp.Status)
	}
	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}
	current, err := semver.ParseTolerant(dep.Version)
	if err != nil {
		return fmt.Errorf("err parsing semver from version %q: %w", dep.Version, err)
	}
	latest, err := semver.ParseTolerant(release.TagName)
	if err != nil {
		return fmt.Errorf("err parsing semver from version %q: %w", release.TagName, err)
	}
	if latest.LTE(current) {
		log.Printf("%s is already at latest version %s", dep.Name, latest.String())
		return nil
	}
	log.Printf("Updating %s to %s", dep.Name, latest.String())
	dep.Version = "v" + latest.String()

	return updateChecksums(dep, &release)
}

func newSyncChecksums() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-checksums",
		Short: "sync checksums for the current version",
		RunE: func(_ *cobra.Command, _ []string) error {
			g := errgroup.Group{}
			for _, dep := range dependencies {
				if dep.AutoUpdate == nil {
					continue
				}
				dep := dep // https://golang.org/doc/faq#closures_and_goroutines
				g.Go(func() error {
					release, err := getRelease(dep)
					if err != nil {
						return fmt.Errorf("getting release %s: %w", dep.Name, err)
					}
					if err := updateChecksums(dep, release); err != nil {
						return fmt.Errorf("updating %s: %w", dep.Name, err)
					}
					return nil
				})
			}
			return g.Wait()
		},
	}
}
func getRelease(dep *dependency) (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", dep.AutoUpdate.Github, dep.Version)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code from server fetching %s: %s", url, resp.Status)
	}
	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func updateChecksums(dep *dependency, release *githubRelease) error {
	urlTpl, err := template.New(dep.Name).Parse(dep.Pattern)
	if err != nil {
		return err
	}

	checksums := map[string]string{}
	if dep.AutoUpdate.Checksums.File != nil {
		var checksumFile []byte
		for _, asset := range release.Assets {
			if asset.Name == *dep.AutoUpdate.Checksums.File {
				resp, err := client.Get(asset.BrowserDownloadURL)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				bs, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				checksumFile = bs
				break
			}
		}
		if len(checksumFile) == 0 {
			return fmt.Errorf("did not find checksums file in release assets")
		}
		scanner := bufio.NewScanner(bytes.NewReader(checksumFile))
		for scanner.Scan() {
			split := strings.Fields(scanner.Text())
			if len(split) != 2 {
				return fmt.Errorf("expected checksums file to be in format 'sha256 filename'")
			}
			filename := strings.TrimLeft(split[1], "./")
			checksums[filename] = split[0]
		}
	}
	if dep.AutoUpdate.Checksums.ReleaseNotes != nil && *dep.AutoUpdate.Checksums.ReleaseNotes {
		regex := regexp.MustCompile(`(?m)([a-z0-9]{64})[\s./]+([a-zA-Z0-9-.]*)`)
		allMatches := regex.FindAllStringSubmatch(release.Body, -1)
		for _, matches := range allMatches {
			if len(matches) != 3 {
				log.Fatalf("weird: %v", matches)
			}
			checksums[matches[2]] = matches[1]
		}
	}

	for os, arches := range dep.Checksums {
		for arch := range arches {
			var url bytes.Buffer
			if err := urlTpl.Execute(&url, fields{Name: dep.Name, Version: dep.Version, OS: os, Arch: arch}); err != nil {
				log.Fatal(err)
			}
			filename := path.Base(url.String())
			checksum, ok := checksums[filename]
			if !ok {
				return fmt.Errorf("no checksum found for filename %s", filename)
			}
			arches[arch] = &checksum
		}
	}
	return nil
}

type githubTransport string

func newGithubTransport() http.RoundTripper {
	token := os.Getenv("GITHUB_TOKEN")
	return githubTransport(token)
}

func (t githubTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", t))
	return http.DefaultTransport.RoundTrip(r)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client = &http.Client{
		Transport: newGithubTransport(),
	}
	var rootCmd = &cobra.Command{Use: "dependencies"}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "hack/dependencies.yml", "config file")
	cobra.OnInitialize(func() {
		depsYAML, err := os.ReadFile(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(depsYAML, &dependencies); err != nil {
			log.Fatal(err)
		}
	})

	update := newUpdateCommand()
	update.AddCommand(newSyncChecksums())
	rootCmd.AddCommand(newInstallCommand(), update)
	rootCmd.Execute()
}
