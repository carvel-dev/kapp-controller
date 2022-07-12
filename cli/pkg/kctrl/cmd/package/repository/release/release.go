package release

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/cmd/package/repository/release/build"
	cmdlocal "github.com/vmware-tanzu/carvel-kapp-controller/cli/pkg/kctrl/local"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
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
	PkgRepoBuildFileName = "pkgrepo-build.yml"
	PkgRepoLocation      = "bundle"
	defaultVersion       = "0.0.0-%d"
	lockOutputFolder     = ".imgpkg"
	defaultPkgRepoName   = "basic.carvel.dev"
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
		o.pkgRepoVersion = fmt.Sprintf(defaultVersion, time.Now().Unix())
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
		pkgRepoName = defaultPkgRepoName
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
	o.ui.PrintActionableText("Lock image references using Kbld and pushing imgpkg bundle to registry")

	// In-memory app for building and pushing images
	builderApp := kcv1alpha1.App{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kctrl-builder",
			Namespace: "in-memory",
			Annotations: map[string]string{
				"kctrl.carvel.dev/local-fetch-0": ".",
			},
		},
		Spec: kcv1alpha1.AppSpec{
			Fetch: []kcv1alpha1.AppFetch{
				{
					Git: &kcv1alpha1.AppFetchGit{},
				},
			},
			Template: []kcv1alpha1.AppTemplate{
				{
					Ytt: &kcv1alpha1.AppTemplateYtt{
						Paths: []string{"packages"},
					},
				}, {
					Kbld: &kcv1alpha1.AppTemplateKbld{
						Paths: []string{},
					},
				},
			},
		},
	}

	buildConfigs := cmdlocal.Configs{
		Apps: []kcv1alpha1.App{builderApp},
	}

	// Create temporary directory for imgpkg lock file
	err = os.Mkdir(filepath.Join(wd, lockOutputFolder), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.RemoveAll(filepath.Join(wd, lockOutputFolder))

	imgpkgLockPath := filepath.Join(wd, lockOutputFolder, "images.yml")
	cmdRunner := NewReleaseCmdRunner(os.Stdout, o.debug, imgpkgLockPath)
	reconciler := cmdlocal.NewReconciler(o.depsFactory, cmdRunner, o.logger)

	err = reconciler.Reconcile(buildConfigs, cmdlocal.ReconcileOpts{
		Local:     true,
		KbldBuild: true,
	})

	if err != nil {
		return err
	}

	var imgpkgBundleURL string

	switch {
	case pkgRepoBuild.Spec.Export.ImgpkgBundle != nil:
		imgpkgOutput, err := ImgpkgRunner{
			Image:             pkgRepoBuild.Spec.Export.ImgpkgBundle.Image,
			Version:           o.pkgRepoVersion,
			Paths:             []string{"packages"},
			UseKbldImagesLock: true,
			ImgLockFilepath:   imgpkgLockPath,
		}.Run()
		if err != nil {
			return err
		}
		imgpkgBundleURL, err = o.imgpkgBundleURLFromStdout(imgpkgOutput)
		if err != nil {
			return err
		}
	}

	o.ui.PrintHeaderText("Output (Step 3/3)")
	artifactWriter := NewArtefactWriter(pkgRepoName, wd)
	err = artifactWriter.WritePackageRepositoryFile(imgpkgBundleURL)
	if err != nil {
		return err
	}
	o.ui.PrintInformationalText("Successfully created package-repository.yml\n")
	o.ui.PrintInformationalText(fmt.Sprintf("\n**Next steps**\nPackage Repository can be consumed in following ways: \n1. Use kctrl to deploy this package repository on the Kubernetes cluster by\nrunning `kctrl package repository add -r demo-pkg-repo --url %s` \n2. Alternatively, use 'kubectl apply -f package-repository.yml' to apply PackageRepository CR directly to your cluster.\n", strings.Split(imgpkgBundleURL, "@")[0]+":"+o.pkgRepoVersion))
	return nil
}

func (o *ReleaseOptions) imgpkgBundleURLFromStdout(imgpkgStdout string) (string, error) {
	lines := strings.Split(imgpkgStdout, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pushed") {
			line = strings.TrimPrefix(line, "Pushed")
			line = strings.Replace(line, "'", "", -1)
			line = strings.Replace(line, " ", "", -1)
			return line, nil
		}
	}
	return "", fmt.Errorf("Could not get imgpkg bundle location")
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
