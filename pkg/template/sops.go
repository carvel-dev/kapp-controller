// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/memdir"
	"golang.org/x/crypto/openpgp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Sops executes sops tool to decrypt configuration.
// Currently supports age and gpg as its decryption backends.
type Sops struct {
	opts       v1alpha1.AppTemplateSops
	appContext AppContext
	coreClient kubernetes.Interface
	cmdRunner  exec.CmdRunner
}

var _ Template = &Sops{}

// NewSops returns sops.
func NewSops(opts v1alpha1.AppTemplateSops, appContext AppContext,
	coreClient kubernetes.Interface, cmdRunner exec.CmdRunner) *Sops {

	return &Sops{opts, appContext, coreClient, cmdRunner}
}

func (t *Sops) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return t.decryptDir(dirPath, nil), false
}

func (t *Sops) TemplateStream(input io.Reader, dirPath string) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating data is not supported"))
}

func (t *Sops) decryptDir(dirPath string, input io.Reader) exec.CmdRunResult {
	result := exec.CmdRunResult{}

	config, err := t.makeConfig()
	if err != nil {
		result.AttachErrorf("Building config: %s", err)
		return result
	}

	defer config.CryptoHomeDir.Remove()

	var args, env []string
	// Be explicit about the config path to avoid sops searching for it
	args = []string{"--config=" + config.ConfigPath}

	switch {
	case t.opts.PGP != nil:
		env = []string{"GNUPGHOME=" + config.CryptoHomeDir.Path()}
	case t.opts.Age != nil:
		env = []string{"SOPS_AGE_KEY_FILE=" + filepath.Join(config.CryptoHomeDir.Path(), "key.txt")}
	default:
		result.AttachErrorf("%s", fmt.Errorf("Unsupported SOPS strategy"))
		return result
	}

	err = t.decryptPathsWithinDir(dirPath, args, env)
	if err != nil {
		result.AttachErrorf("Decrypting dir: %s", err)
	}

	return result
}

func (t *Sops) decryptPathsWithinDir(dirPath string, args, env []string) error {
	var selectedDirPaths []string

	if len(t.opts.Paths) == 0 {
		selectedDirPaths = append(selectedDirPaths, dirPath)
	} else {
		for _, path := range t.opts.Paths {
			checkedPath, err := memdir.ScopedPath(dirPath, path)
			if err != nil {
				return fmt.Errorf("Checking path: %s", err)
			}

			info, err := os.Stat(checkedPath)
			if err != nil {
				return err
			}

			isDir, err := t.checkDirOrDecryptSopsFile(info, checkedPath, args, env)
			if err != nil {
				return err
			} else if isDir {
				selectedDirPaths = append(selectedDirPaths, checkedPath)
			}
		}
	}

	for _, dirPath := range selectedDirPaths {
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			_, err = t.checkDirOrDecryptSopsFile(info, path, args, env)
			return err
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Sops) checkDirOrDecryptSopsFile(info os.FileInfo, path string, args, env []string) (bool, error) {
	if info.IsDir() {
		return true, nil
	}
	matched, newPath := t.isSopsFile(path)
	if matched {
		err := t.decryptSopsFile(path, newPath, args, env)
		if err != nil {
			return false, fmt.Errorf("Decrypting file '%s': %s", path, err)
		}
	}
	return false, nil
}

func (t *Sops) decryptSopsFile(path, newPath string, args, env []string) error {
	decryptArgs := append(append([]string{}, args...), "--decrypt", path)

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("sops", decryptArgs...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := t.cmdRunner.Run(cmd)
	if err != nil {
		return fmt.Errorf("Running sops: %s, %v", stderrBs.String(), err)
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("Removing encrypted file: %s", err)
	}

	contentsBs, err := t.shapeDecryptedContents(stdoutBs.Bytes())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(newPath, contentsBs, 0600)
	if err != nil {
		return fmt.Errorf("Writing decrypted file: %s", err)
	}

	return nil
}

var (
	sopsExts = map[string]string{
		".sops.yml":  ".yml",
		".sops.yaml": ".yaml",
	}
)

func (*Sops) isSopsFile(path string) (bool, string) {
	for ext, replExt := range sopsExts {
		if strings.HasSuffix(path, ext) {
			return true, strings.TrimSuffix(path, ext) + replExt
		}
	}
	return false, ""
}

func (*Sops) shapeDecryptedContents(contentsBs []byte) ([]byte, error) {
	// TODO we currently do not support any kind of enveloping
	// which might be needed for cases like ytt data values

	// const (
	// 	dataKey = "sops.k14s.io/data"
	// )

	// var contents map[string]interface{}

	// err := yaml.Unmarshal(contentsBs, &contents)
	// if err != nil {
	// 	return nil, fmt.Errorf("Unmarshaling decrypted file as YAML: %s", err)
	// }

	// if dataVal, found := contents[dataKey]; found {
	// 	dataStr, ok := dataVal.(string)
	// 	if !ok {
	// 		return nil, fmt.Errorf("Expected key '%s' value to be a string", dataKey)
	// 	}
	// 	contentsBs = []byte(dataStr)
	// }

	return contentsBs, nil
}

type sopsCryptoStrategy int8 // just a hunch that we'll never have a 257th encryption provider but if you're still using this in the year 2666 and upgrading to an int16, sorry about this and also all the CO2, single-use plastic, and whatever else we did wrong
const (
	pgp sopsCryptoStrategy = iota
	age
)

type sopsConfig struct {
	CryptoHomeDir *memdir.TmpDir
	ConfigPath    string
	Strategy      sopsCryptoStrategy
}

func (sc *sopsConfig) CreateConfigPath() error {
	sc.ConfigPath = filepath.Join(sc.CryptoHomeDir.Path(), ".sops.yml")

	err := ioutil.WriteFile(sc.ConfigPath, []byte("{}"), 0600)
	if err != nil {
		return fmt.Errorf("Generating config file: %s", err)
	}
	return nil
}

// extractKeysFromSecretRefContents interprets the secretContents according to the encryption strategy configured in sc.
func (sc *sopsConfig) ExtractKeysFromSecretRefContents(secretContents string) error {
	switch sc.Strategy {
	case pgp:
		err := gpgKeyring{secretContents}.Write(sc.CryptoHomeDir.Path())
		if err != nil {
			return fmt.Errorf("Generating secring.gpg: %s", err)
		}
	case age:
		err := ioutil.WriteFile(filepath.Join(sc.CryptoHomeDir.Path(), "key.txt"), []byte(secretContents), 0600)
		if err != nil {
			return fmt.Errorf("Creating key.txt file: %s", err)
		}
	default:
		return fmt.Errorf("Unrecognized sops encryption strategy %d", sc.Strategy)
	}
	return nil
}

func (t *Sops) makeConfig() (sopsConfig, error) {
	cryptoHomeDir := memdir.NewTmpDir("template-sops-config")
	config := sopsConfig{cryptoHomeDir, "", 0}

	err := cryptoHomeDir.Create()
	if err != nil {
		return config, err
	}

	var secretRef *v1alpha1.AppTemplateSopsPrivateKeysSecretRef
	if t.opts.PGP != nil {
		secretRef = t.opts.PGP.PrivateKeysSecretRef
		config.Strategy = pgp
	} else if t.opts.Age != nil {
		secretRef = t.opts.Age.PrivateKeysSecretRef
		config.Strategy = age
	}

	secretContents, err := t.getSecretContents(*secretRef)
	if err != nil {
		cryptoHomeDir.Remove()
		return config, fmt.Errorf("Getting private keys secret: %s", err)
	}

	err = config.ExtractKeysFromSecretRefContents(secretContents)
	if err != nil {
		cryptoHomeDir.Remove()
		return config, err
	}

	err = config.CreateConfigPath()
	if err != nil {
		cryptoHomeDir.Remove()
		return config, err
	}

	return config, nil
}

func (t *Sops) getSecretContents(secretRef v1alpha1.AppTemplateSopsPrivateKeysSecretRef) (string, error) {
	secret, err := t.coreClient.CoreV1().Secrets(t.appContext.Namespace).Get(context.Background(), secretRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	var result string

	for _, val := range secret.Data {
		result += string(val) + "\n"
	}

	return result, nil
}

type gpgKeyring struct {
	contents string
}

func (k gpgKeyring) Write(dirPath string) error {
	// TODO currently this only reads a single key
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(k.contents))
	if err != nil {
		return fmt.Errorf("Reading private keys: %s", err)
	}

	if len(entityList) < 1 {
		return fmt.Errorf("Expected to find at least one private key")
	}

	file, err := os.OpenFile(filepath.Join(dirPath, "secring.gpg"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Creating secring.gpg file: %s", err)
	}

	// Ignore dup close error in happy path
	defer file.Close()

	for _, entity := range entityList {
		if entity == nil {
			return fmt.Errorf("Expected to find private key within entity")
		}

		err = entity.PrivateKey.Serialize(file)
		if err != nil {
			return fmt.Errorf("Serializing pk: %s", err)
		}

		for _, ident := range entity.Identities {
			err = ident.UserId.Serialize(file)
			if err != nil {
				return fmt.Errorf("Serializing ident user id: %s", err)
			}
			err = ident.SelfSignature.Serialize(file)
			if err != nil {
				return fmt.Errorf("Serializing ident self sig: %s", err)
			}
		}
		for _, subkey := range entity.Subkeys {
			err = subkey.PrivateKey.Serialize(file)
			if err != nil {
				return fmt.Errorf("Serializing subkey pk: %s", err)
			}
			err = subkey.Sig.Serialize(file)
			if err != nil {
				return fmt.Errorf("Serializing subkey sig: %s", err)
			}
		}
	}

	// Make sure to close successfully to make sure contents are flushed
	err = file.Close()
	if err != nil {
		return fmt.Errorf("Closing secring.gpg: %s", err)
	}

	return nil
}
