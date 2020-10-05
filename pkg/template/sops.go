// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/exec"
	"github.com/k14s/kapp-controller/pkg/memdir"
	"golang.org/x/crypto/openpgp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Sops struct {
	opts        v1alpha1.AppTemplateSops
	genericOpts GenericOpts
	coreClient  kubernetes.Interface
}

var _ Template = &Sops{}

func NewSops(opts v1alpha1.AppTemplateSops,
	genericOpts GenericOpts, coreClient kubernetes.Interface) *Sops {
	return &Sops{opts, genericOpts, coreClient}
}

func (t *Sops) TemplateDir(dirPath string) (exec.CmdRunResult, bool) {
	return t.decryptDir(dirPath, nil), false
}

func (t *Sops) TemplateStream(input io.Reader) exec.CmdRunResult {
	return exec.NewCmdRunResultWithErr(fmt.Errorf("Templating data is not supported"))
}

func (t *Sops) decryptDir(dirPath string, input io.Reader) exec.CmdRunResult {
	result := exec.CmdRunResult{}

	args := []string{}
	env := []string{}

	switch {
	case t.opts.PGP != nil:
		gpgHomeDir, err := t.gpgHomeWithKeyRing()
		if err != nil {
			result.AttachErrorf("Building PGP key ring: %s", err)
			return result
		}

		defer gpgHomeDir.Remove()

		args = []string{} // no additional args
		env = append(env, "GNUPGHOME="+gpgHomeDir.Path())

	default:
		result.AttachErrorf("%s", fmt.Errorf("Unsupported SOPS strategy"))
		return result
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		err = t.decryptSopsFile(path, info, args, env)
		if err != nil {
			return fmt.Errorf("Decrypting file '%s': %s", path, err)
		}
		return nil
	})

	result.AttachErrorf("Decrypting dir: %s", err)

	return result
}

func (t *Sops) decryptSopsFile(path string, info os.FileInfo, args []string, env []string) error {
	// Skip non-sops files
	if info.IsDir() {
		return nil
	}

	cont, newPath := t.isSopsFile(path)
	if !cont {
		return nil
	}

	decryptArgs := []string{}
	decryptArgs = append(decryptArgs, args...)
	decryptArgs = append(decryptArgs, "--decrypt", path)

	var stdoutBs, stderrBs bytes.Buffer

	cmd := goexec.Command("sops", decryptArgs...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Running sops: %s", err)
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

func (t *Sops) gpgHomeWithKeyRing() (*memdir.TmpDir, error) {
	if t.opts.PGP.PrivateKeysSecretRef == nil {
		return nil, fmt.Errorf("Expected to have private keys secret ref specified")
	}

	privateKeys, err := t.getFromSecret(*t.opts.PGP.PrivateKeysSecretRef)
	if err != nil {
		return nil, fmt.Errorf("Getting private keys secret: %s", err)
	}

	gpgHomeDir := memdir.NewTmpDir("template-sops-gpghome")

	err = gpgHomeDir.Create()
	if err != nil {
		return nil, err
	}

	err = gpgKeyring{privateKeys}.Write(gpgHomeDir.Path())
	if err != nil {
		return nil, fmt.Errorf("Generating secring.gpg: %s", err)
	}

	return gpgHomeDir, nil
}

func (t *Sops) getFromSecret(secretRef v1alpha1.AppTemplateSopsPGPPrivateKeysSecretRef) (string, error) {
	secret, err := t.coreClient.CoreV1().Secrets(t.genericOpts.Namespace).Get(secretRef.Name, metav1.GetOptions{})
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
