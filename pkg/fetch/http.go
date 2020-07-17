package fetch

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HTTP struct {
	opts       v1alpha1.AppFetchHTTP
	nsName     string
	coreClient kubernetes.Interface
}

func NewHTTP(opts v1alpha1.AppFetchHTTP, nsName string, coreClient kubernetes.Interface) *HTTP {
	return &HTTP{opts, nsName, coreClient}
}

func (t *HTTP) Retrieve(dstPath string) error {
	if len(t.opts.URL) == 0 {
		return fmt.Errorf("Expected non-empty URL")
	}

	tmpFile, err := ioutil.TempFile("", "kapp-controller-files-http")
	if err != nil {
		return err
	}

	defer os.Remove(tmpFile.Name())

	err = t.downloadFileAndChecksum(tmpFile)
	if err != nil {
		return fmt.Errorf("Downloading URL: %s", err)
	}

	contentExtractorFuncs := []func(string, string) (bool, error){
		t.tryZip,
		t.tryTgz,
		t.tryTar,
	}

	for _, f := range contentExtractorFuncs {
		final, err := f(tmpFile.Name(), dstPath)
		if final {
			return err
		}
	}

	return t.tryPlain(tmpFile.Name(), dstPath)
}

func (t *HTTP) downloadFile(dst io.Writer) error {
	req, err := http.NewRequest("GET", t.opts.URL, nil)
	if err != nil {
		return fmt.Errorf("Building request: %s", err)
	}

	err = t.addAuth(req)
	if err != nil {
		return fmt.Errorf("Adding auth to request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Initiating URL download: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected 200 OK, but was '%s'", resp.Status)
	}

	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return fmt.Errorf("Writing downloaded content: %s", err)
	}

	return nil
}

func (t *HTTP) downloadFileAndChecksum(dst io.Writer) error {
	var digestName, expectedDigestVal string
	var digestDst hash.Hash

	switch {
	case len(t.opts.SHA256) > 0:
		digestName = "sha256"
		digestDst = sha256.New()
		expectedDigestVal = t.opts.SHA256
		dst = io.MultiWriter(dst, digestDst)
	}

	err := t.downloadFile(dst)
	if err != nil {
		return err
	}

	if len(expectedDigestVal) > 0 {
		actualDigestVal := fmt.Sprintf("%x", digestDst.Sum(nil))

		if expectedDigestVal != actualDigestVal {
			errMsg := "Expected digest to match '%s:%s', but was '%s:%s'"
			return fmt.Errorf(errMsg, digestName, expectedDigestVal, digestName, actualDigestVal)
		}
	}

	return nil
}

func (t *HTTP) tryZip(path, dstPath string) (bool, error) {
	zipArchive, err := zip.OpenReader(path)
	if err != nil {
		return false, fmt.Errorf("Opening zip archive: %s", err)
	}

	defer zipArchive.Close()

	for _, f := range zipArchive.File {
		if strings.HasSuffix(f.Name, "/") {
			// TODO should we make empty directories?
			continue
		}

		srcZipFile, err := f.Open()
		if err != nil {
			return false, fmt.Errorf("Opening zip file: %s", err)
		}

		err = t.writeIntoFileAndClose(srcZipFile, dstPath, f.Name)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (t *HTTP) tryTgz(path, dstPath string) (bool, error) {
	return t.tryTarWithGzip(path, dstPath, true)
}

func (t *HTTP) tryTar(path, dstPath string) (bool, error) {
	return t.tryTarWithGzip(path, dstPath, false)
}

func (t *HTTP) tryTarWithGzip(path, dstPath string, gzipped bool) (bool, error) {
	plainFile, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("Opening archive: %s", err)
	}

	defer plainFile.Close()

	var fileReader io.Reader

	if gzipped {
		gzipFile, err := gzip.NewReader(plainFile)
		if err != nil {
			return false, fmt.Errorf("Opening gzip archive: %s", err)
		}
		fileReader = gzipFile
	} else {
		fileReader = plainFile
	}

	tarReader := tar.NewReader(fileReader)
	readEntries := false

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return readEntries, fmt.Errorf("Reading next tar header: %s", err)
		}

		readEntries = true

		switch header.Typeflag {
		case tar.TypeDir:
			// TODO should we make empty directories?
			continue

		case tar.TypeReg:
			err = t.writeIntoFile(tarReader, dstPath, header.Name)
			if err != nil {
				return true, err
			}

		default:
			return false, fmt.Errorf("Unknown file '%s' (%d)", header.Name, header.Typeflag)
		}
	}

	return true, nil
}

func (t *HTTP) tryPlain(path, dstPath string) error {
	parsedURL, err := gourl.Parse(t.opts.URL)
	if err != nil {
		return fmt.Errorf("Parsing URL: %s", err)
	}

	pathSegs := strings.Split(parsedURL.Path, "/")

	fileName := pathSegs[len(pathSegs)-1]
	if len(fileName) == 0 {
		fileName = "content"
	}

	srcFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Opening file %s: %s", path, err)
	}

	defer srcFile.Close()

	// Cannot just move since it may be on a different device
	return t.writeIntoFile(srcFile, dstPath, fileName)
}

func (t *HTTP) writeIntoFile(srcFile io.Reader, dstPath, additionalPath string) error {
	dstFilePath, err := memdir.ScopedPath(dstPath, additionalPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(dstFilePath), 0700)
	if err != nil {
		return fmt.Errorf("Making intermediate dir: %s", err)
	}

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("Creating dst file: %s", err)
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("Copying into dst file: %s", err)
	}

	return nil
}

func (t *HTTP) writeIntoFileAndClose(srcFile io.ReadCloser, dstPath, additionalPath string) error {
	defer srcFile.Close()
	return t.writeIntoFile(srcFile, dstPath, additionalPath)
}

func (t *HTTP) addAuth(req *http.Request) error {
	if t.opts.SecretRef == nil {
		return nil
	}

	secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(t.opts.SecretRef.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for name, _ := range secret.Data {
		switch name {
		case corev1.BasicAuthUsernameKey:
		case corev1.BasicAuthPasswordKey:
		default:
			return fmt.Errorf("Unknown secret field '%s' in secret '%s'", name, secret.Name)
		}
	}

	if _, found := secret.Data[corev1.BasicAuthUsernameKey]; found {
		req.SetBasicAuth(string(secret.Data[corev1.BasicAuthUsernameKey]),
			string(secret.Data[corev1.BasicAuthPasswordKey]))
	}

	return nil
}
