package release

import (
	"fmt"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"
	"time"

	appInit "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/release/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	cmdcore "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/core"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/logger"
)

type ReleaseOptions struct {
	ui          cmdcore.AuthoringUI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	pkgRepoVersion string
	chdir          string
	outputLocation string
	debug          bool
}

const (
	PkgRepoBuildFileName  = "pkgrepo-build.yml"
	PkgRepositoryFileName = "package-repository.yml"
	DefaultVersion        = "0.0.0-build.%d"
	DefaultPkgRepoName    = "sample-repo.carvel.dev"
	PackagesDirectory     = "packages"

	LockOutputFolder = ".imgpkg"
	LockOutputFile   = "images.yml"
)

func NewReleaseOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *ReleaseOptions {
	return &ReleaseOptions{ui: cmdcore.NewAuthoringUIImpl(ui), depsFactory: depsFactory, logger: logger}
}

func NewReleaseCmd(o *ReleaseOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Build and create a package repository (experimental)",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run() },
		Annotations: map[string]string{
			cmdcore.PackageAuthoringCommandsHelpGroup.Key: cmdcore.PackageAuthoringCommandsHelpGroup.Value,
		},
	}

	cmd.Flags().StringVarP(&o.pkgRepoVersion, "version", "v", "", "Version to be released")
	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Location of the working directory")
	cmd.Flags().StringVar(&o.outputLocation, "copy-to", "", "Output location for pkgrepo-build.yml")
	cmd.Flags().BoolVar(&o.debug, "debug", false, "Include debug output")

	return cmd
}

func (o *ReleaseOptions) Run() error {
	o.ui.PrintHeaderText("\nPrerequisites")
	o.ui.PrintInformationalText("1. A `packages` directory containing Package and PackageMetadata files should be present in the working directory.\n2. The host must be authorized to push images to a registry (can be set up by running `docker login`)\n")

	if o.pkgRepoVersion == "" {
		o.pkgRepoVersion = fmt.Sprintf(DefaultVersion, time.Now().Unix())
	}

	if o.chdir != "" {
		err := os.Chdir(o.chdir)
		if err != nil {
			return err
		}
	}

	pkgRepoBuild, err := o.getPackageRepositoryBuild(PkgRepoBuildFileName)
	if err != nil {
		return err
	}

	o.ui.PrintHeaderText("\nBasic Information")
	pkgRepoName := pkgRepoBuild.Name
	if pkgRepoName == "" {
		pkgRepoName = DefaultPkgRepoName
	}
	textOpts := ui.TextOpts{
		Label:        "Enter the package repository name",
		Default:      pkgRepoName,
		ValidateFunc: nil,
	}
	pkgRepoName, err = o.ui.AskForText(textOpts)
	if err != nil {
		return err
	}
	pkgRepoBuild.Name = pkgRepoName

	o.ui.PrintHeaderText("Registry URL")
	o.ui.PrintInformationalText("The bundle created needs to be pushed to an OCI registry." +
		" Registry URL format: <REGISTRY_URL/REPOSITORY_NAME> e.g. index.docker.io/k8slt/sample-bundle")
	defaultRegistryURL := pkgRepoBuild.Spec.Export.ImgpkgBundle.Image
	textOpts = ui.TextOpts{
		Label:        "Enter the registry url",
		Default:      defaultRegistryURL,
		ValidateFunc: nil,
	}
	registryURL, err := o.ui.AskForText(textOpts)
	if err != nil {
		return err
	}

	pkgRepoBuild.Spec.Export.ImgpkgBundle.Image = strings.TrimSpace(registryURL)
	err = pkgRepoBuild.WriteToFile()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	o.ui.PrintInformationalText("kbld ensures that all image references are resolved to an immutable reference.")
	o.ui.PrintActionableText("Lock image references using kbld")

	packagesFolderPath := filepath.Join(wd, PackagesDirectory)
	_, err = os.Stat(packagesFolderPath)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Expected to find `packages` directory in the root")
	}

	// Make lock output directory if it does not exist
	tmpImgpkgFolder := filepath.Join(wd, LockOutputFolder)
	_, err = os.Stat(tmpImgpkgFolder)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(tmpImgpkgFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	defer os.RemoveAll(filepath.Join(wd, LockOutputFolder))

	tempImgpkgLockPath := filepath.Join(LockOutputFolder, LockOutputFile)

	//running kbld
	kbldCmdRunner := exec.NewPlainCmdRunner()
	cmd := goexec.Command("kbld", "-f", PackagesDirectory, "--imgpkg-lock-output", tempImgpkgLockPath)
	err = kbldCmdRunner.Run(cmd)
	if err != nil {
		return err
	}
	o.ui.PrintCmdExecutionOutput(fmt.Sprintf("\n$ %s", strings.Join(cmd.Args, " ")))

	var bundleURL string

	switch {
	case pkgRepoBuild.Spec.Export.ImgpkgBundle != nil:
		imgpkgRunner := ImgpkgRunner{
			BundlePath:        fmt.Sprintf("%s:%s", pkgRepoBuild.Spec.Export.ImgpkgBundle.Image, o.pkgRepoVersion),
			Paths:             []string{"packages"},
			UseKbldImagesLock: true,
			ImgLockFilepath:   tempImgpkgLockPath,
			UI:                o.ui,
		}
		bundleURL, err = imgpkgRunner.Run()
		if err != nil {
			return err
		}
	}

	artifactWriter := NewArtifactWriter(pkgRepoName, wd)
	err = artifactWriter.WritePackageRepositoryFile(bundleURL)
	if err != nil {
		return err
	}
	o.ui.PrintInformationalText("Successfully created package-repository.yml\n")
	o.ui.PrintHeaderText("\nNext steps")
	o.ui.PrintInformationalText("1. Add the package repository to the Kubernetes cluster by running `kctrl package repository add -r <REPO_NAME> --url <PKG_REPO_BUNDLE_URL>`\n2. Alternatively, use `kapp` or `kubectl` to apply `package-repository.yml` to your cluster.\n")
	return nil
}

func (o *ReleaseOptions) getPackageRepositoryBuild(pkgRepoBuildFilePath string) (*build.PackageRepoBuild, error) {
	var packageRepoBuild *build.PackageRepoBuild
	exists, err := appInit.IsFileExists(pkgRepoBuildFilePath)
	if err != nil {
		return nil, err
	}

	if exists {
		packageRepoBuild, err = o.newPackageRepoBuildFromFile(pkgRepoBuildFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		packageRepoBuild = &build.PackageRepoBuild{
			TypeMeta: metav1.TypeMeta{
				Kind:       build.PkgRepoBuildKind,
				APIVersion: build.PkgRepoBuildAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(time.Now()),
			},
			Spec: build.PackageRepoBuildSpec{
				Export: &build.PackageRepoBuildExport{
					ImgpkgBundle: &build.PackageRepoBuildExportImgpkgBundle{},
				},
			},
		}
	}
	return packageRepoBuild, nil
}

func (o *ReleaseOptions) newPackageRepoBuildFromFile(filePath string) (*build.PackageRepoBuild, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var packageRepoBuild build.PackageRepoBuild
	err = yaml.Unmarshal(content, &packageRepoBuild)
	if err != nil {
		return nil, err
	}
	return &packageRepoBuild, nil
}
