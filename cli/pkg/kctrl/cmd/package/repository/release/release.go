package release

import (
	"fmt"
	"os"
	goexec "os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/release/build"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/yaml"

	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/app/init/common"
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
	DefaultVersion        = "0.0.0-%d"
	DefaultPkgRepoName    = "basic.carvel.dev"
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
		Short: "Build and create a package repository",
		RunE:  func(_ *cobra.Command, args []string) error { return o.Run() },
	}

	cmd.Flags().StringVarP(&o.pkgRepoVersion, "version", "v", "", "Version to be released")
	cmd.Flags().StringVar(&o.chdir, "chdir", "", "Working directory with repo which needs to be bundles")
	cmd.Flags().StringVar(&o.outputLocation, "copy-to", "", "Output location for pkgrepo-build.yml")
	cmd.Flags().BoolVar(&o.debug, "debug", false, "Include debug output")

	return cmd
}

func (o *ReleaseOptions) Run() error {
	o.ui.PrintHeaderText("\nPrerequisites")
	o.ui.PrintInformationalText("Welcome! Before we start on creating the package repository please ensure the following prerequisites are met:\n 1. You have the `packages` directory which contains Package and PackageMetadata files and should be available on the root directory (or use -–chdir to change root directory).\n 2. You have access to an OCI registry, and you have authenticated locally so that you can push images. e.g. docker login\n")

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

	o.ui.PrintHeaderText("\nBasic Information(Step 1/3)")
	o.ui.PrintInformationalText("A package repository name is the name with which it will be referenced while deploying on the cluster.")
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

	o.ui.PrintInformationalText("To create package repository, we will create an imgpkg bundle first. imgpkg is a Carvel tool that allows users to package, distribute, and relocate a set of files as one OCI artifact: a bundle.")
	o.ui.PrintInformationalText("\nA package repository bundle is an imgpkg bundle that holds PackageMetadata and Package CRs.\n")

	o.ui.PrintHeaderText("\nRegistry URL(Step 2/3)")
	defaultRegistryURL := pkgRepoBuild.Spec.Export.ImgpkgBundle.Image
	textOpts = ui.TextOpts{
		Label:        "Enter the registry url to push the package repository bundle",
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

	o.ui.PrintInformationalText("Let's use `kbld` to create immutable image reference. Kbld scans all the files in bundle configuration for any references of images and creates a mapping of image tags to a URL with sha256 digest.\n")
	o.ui.PrintActionableText("Lock image references using Kbld")

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
			BundlePath:        fmt.Sprintf("%s:build-%d", pkgRepoBuild.Spec.Export.ImgpkgBundle.Image, time.Now().Unix()),
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

	artifactWriter := NewArtefactWriter(pkgRepoName, wd)
	err = artifactWriter.WritePackageRepositoryFile(bundleURL)
	if err != nil {
		return err
	}
	o.ui.PrintInformationalText("Successfully created package-repository.yml\n")
	o.ui.PrintInformationalText("\n**Next steps**\nPackage Repository can be consumed in following ways: \n1. Use kctrl to deploy this package repository on the Kubernetes cluster by\nrunning `kctrl package repository add -r demo-pkg-repo --url <PKG_REPO_BUNDLE_URL>` \n2. Alternatively, use 'kubectl apply -f package-repository.yml' to apply PackageRepository CR directly to your cluster.\n")
	return nil
}

func (o *ReleaseOptions) getPackageRepositoryBuild(pkgRepoBuildFilePath string) (*build.PackageRepoBuild, error) {
	var packageRepoBuild *build.PackageRepoBuild
	exists, err := common.IsFileExists(pkgRepoBuildFilePath)
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
