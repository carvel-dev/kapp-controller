// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
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
	dependencies depslice
	client       *http.Client
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	Body string `json:"body"`
}

type platform struct {
	OS,
	Arch string
}
type templateArgs struct {
	platform
	*dependency
}

type depslice []*dependency

// iterates through each dependency, calling `f` on a Goroutine per dependency and aggregating errors
func (d depslice) each(f func(context.Context, *dependency) error) error {
	g, ctx := errgroup.WithContext(context.Background())
	for _, dep := range d {
		dep := dep // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			return f(ctx, dep)
		})
	}
	return g.Wait()
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
		Run: func(_ *cobra.Command, _ []string) {
			if err := c.Install(); err != nil {
				log.Fatal(err)
			}
		},
	}
	command.Flags().StringVar(&c.os, "os", runtime.GOOS, "os to install dependencies for")
	command.Flags().StringVar(&c.arch, "arch", runtime.GOARCH, "arch to install dependencies for")
	command.Flags().StringVarP(&c.destDir, "destination", "d", "", "directory to which binaries will be installed (required)")
	command.Flags().BoolVar(&c.dev, "dev", false, "only install dev dependencies (default false)")
	command.MarkFlagRequired("destination")
	return command
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
		Run: func(_ *cobra.Command, _ []string) {
			err := dependencies.each(func(ctx context.Context, dep *dependency) error {
				if err := dep.update(ctx); err != nil {
					return fmt.Errorf("updating %s: %w", dep.Name, err)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}

func newSyncChecksums() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-checksums",
		Short: "sync checksums for the current version",
		Run: func(_ *cobra.Command, _ []string) {
			err := dependencies.each(func(ctx context.Context, dep *dependency) error {
				release, err := dep.getRelease(ctx)
				if err != nil {
					return fmt.Errorf("getting release %s: %w", dep.Name, err)
				}
				dep.Version = release.TagName
				if err := dep.updateChecksums(ctx); err != nil {
					return fmt.Errorf("updating %s: %w", dep.Name, err)
				}
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}

// dowlnoads, verifies checksum, and (optionally) extracts the binary from a tarball
func (c *installCommand) Install() error {
	return dependencies.each(func(ctx context.Context, dep *dependency) error {
		if c.dev && !dep.Dev {
			return nil
		}
		if err := c.downloadAndVerify(ctx, dep); err != nil {
			return fmt.Errorf("downloading %s: %w", dep.Name, err)
		}
		return nil
	})
}

func (c *installCommand) downloadAndVerify(ctx context.Context, dep *dependency) error {
	arches, ok := dep.Checksums[c.os]
	if !ok {
		return fmt.Errorf("not supported on os %s", c.os)
	}
	checksum, ok := arches[c.arch]
	if !ok {
		return fmt.Errorf("not supported on platform %s/%s", c.os, c.arch)
	}

	opts := platform{OS: c.os, Arch: c.arch}
	blob, err := dep.download(ctx, opts)
	if err != nil {
		return err
	}
	log.Printf("%s downloaded", dep.Name)

	newChecksum, err := blob.Checksum()
	if err != nil {
		return err
	}

	if checksum != newChecksum {
		return fmt.Errorf("wrong checksum, expected: %s, got: %s", checksum, newChecksum)
	}

	log.Printf("%s validated", dep.Name)

	file, err := blob.Binary()
	if err != nil {
		return err
	}
	filepath := path.Join(c.destDir, dep.Name)
	dest, err := os.Create(filepath)
	if err != nil {
		return err
	}
	if err := os.Chmod(filepath, 0777); err != nil {
		return err
	}
	if _, err := io.Copy(dest, file); err != nil {
		return err
	}
	return nil
}

type checksum string

// returns a hex-encoded sha256 hash of the given bytes
func newChecksumFromBytes(bs []byte) (checksum, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, bytes.NewReader(bs))
	if err != nil {
		return "", err
	}
	return checksum(hex.EncodeToString(hasher.Sum(nil))), nil
}

type blob interface {
	//  Returns a checksum of the blob
	Checksum() (checksum, error)
	// Returns the executable binary
	Binary() (io.Reader, error)
}

// tarballBlob wraps a .tgz file to implement the blob interface
type tarballBlob struct {
	content []byte
	subpath string
}

func (t *tarballBlob) Checksum() (checksum, error) {
	return newChecksumFromBytes(t.content)
}

func (t *tarballBlob) Binary() (io.Reader, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(t.content))
	if err != nil {
		return nil, err
	}
	tgzReader := tar.NewReader(gzReader)
	for {
		hdr, err := tgzReader.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("tarball subpath not found: %s", t.subpath)
		}
		if hdr.Name == t.subpath {
			return tgzReader, nil
		}
	}
}

// fileBlob wraps an executable file to implement the blob interface
type fileBlob []byte

func (f fileBlob) Checksum() (checksum, error) {
	return newChecksumFromBytes(f)
}

func (f fileBlob) Binary() (io.Reader, error) {
	return bytes.NewReader(f), nil
}

type dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// Template for the download URL of the artifact
	URLTemplate string `json:"urlTemplate"`
	// Whether this should be installed by the --dev flag
	Dev bool `json:"dev"`
	// Github repo (org/name)
	Repo string `json:"repo,omitempty"`
	// map["linux"]map["amd64"] => "sha256"
	Checksums      map[string]map[string]checksum `json:"checksums"`
	TarballSubpath *string                        `json:"tarballSubpath,omitempty"`
}

// downloads the dependency by its url template
func (d *dependency) download(ctx context.Context, opts platform) (blob, error) {
	templateArgs := templateArgs{opts, d}
	urlTpl, err := template.New(d.Name).Parse(d.URLTemplate)
	if err != nil {
		return nil, err
	}
	var url bytes.Buffer
	if err := urlTpl.Execute(&url, templateArgs); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code retrieving url: %s: %s", url.String(), resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if d.TarballSubpath != nil {
		tarballTpl, err := template.New(d.Name).Parse(*d.TarballSubpath)
		if err != nil {
			return nil, err
		}
		var tarballSubpath bytes.Buffer
		if err := tarballTpl.Execute(&tarballSubpath, templateArgs); err != nil {
			return nil, err
		}
		return &tarballBlob{
			subpath: tarballSubpath.String(),
			content: content,
		}, nil
	}
	return fileBlob(content), nil
}

// updates the version of the dependency to the one GitHub considers latest and updates checksums
func (d *dependency) update(ctx context.Context) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", d.Repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
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

	// ParseTolerant allows the "v" prefix in the version string
	current, err := semver.ParseTolerant(d.Version)
	if err != nil {
		return fmt.Errorf("err parsing semver from version %q: %w", d.Version, err)
	}
	latest, err := semver.ParseTolerant(release.TagName)
	if err != nil {
		return fmt.Errorf("err parsing semver from version %q: %w", release.TagName, err)
	}

	// short-circuit if we're already at latest
	if latest.LTE(current) {
		log.Printf("%s is already at latest version %s", d.Name, latest.String())
		return nil
	}
	log.Printf("Updating %s to %s", d.Name, latest.String())

	// add the "v" prefix back
	d.Version = "v" + latest.String()

	return d.updateChecksums(ctx)
}

// get the release specified by the current version of this dependency
func (d *dependency) getRelease(ctx context.Context) (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", d.Repo, d.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
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

// update the checksums of this dependency to those in the github release
func (d *dependency) updateChecksums(ctx context.Context) error {
	for os, arches := range d.Checksums {
		for arch := range arches {
			blob, err := d.download(ctx, platform{OS: os, Arch: arch})
			if err != nil {
				return err
			}
			checksum, err := blob.Checksum()
			if err != nil {
				return err
			}
			d.Checksums[os][arch] = checksum
			log.Printf("%s %s %s/%s checksum synced", d.Name, d.Version, os, arch)
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
	if t != "" && strings.Contains(r.URL.Host, "github.com") {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", t))
	}
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		respBody, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, fmt.Errorf("error dumping response: status %s", resp.Status)
		}
		log.Printf("%s:\n%s", r.URL, respBody)
		return nil, fmt.Errorf("non-200 status code from server: %s", resp.Status)
	}
	return resp, err
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client = &http.Client{
		// try to use a Github API token from the environment to avoid rate-limiting
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
