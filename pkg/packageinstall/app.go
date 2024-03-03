// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	ManuallyControlledAnnKey = "ext.packaging.carvel.dev/manually-controlled"

	HelmTemplateOverlayNameKey      = "ext.packaging.carvel.dev/helm-template-name"
	HelmTemplateOverlayNameSpaceKey = "ext.packaging.carvel.dev/helm-template-namespace"

	// Resulting secret names are sorted deterministically by suffix
	ExtYttPathsFromSecretNameAnnKey  = "ext.packaging.carvel.dev/ytt-paths-from-secret-name"
	ExtHelmPathsFromSecretNameAnnKey = "ext.packaging.carvel.dev/helm-template-values-from-secret-name"

	// ExtYttDataValuesOverlaysAnnKey if set, adds the pkgi's values secrets as overlays/paths, not as values, to the app
	ExtYttDataValuesOverlaysAnnKey = "ext.packaging.carvel.dev/ytt-data-values-overlays"

	ExtFetchSecretNameAnnKeyFmt = "ext.packaging.carvel.dev/fetch-%d-secret-name"
)

// NewApp creates a new instance of v1alpha1.App based on the provided parameters.
// It takes an existingApp, pkgInstall, pkgVersion, and opts, and returns the newly created App.
func NewApp(existingApp *v1alpha1.App, pkgInstall *pkgingv1alpha1.PackageInstall, pkgVersion datapkgingv1alpha1.Package, opts Opts) (*v1alpha1.App, error) {
	desiredApp := existingApp.DeepCopy()

	if _, found := existingApp.Annotations[ManuallyControlledAnnKey]; found {
		// Skip all updates to App CR if annotation is present
		return desiredApp, nil
	}

	desiredApp.Name = pkgInstall.Name
	desiredApp.Namespace = pkgInstall.Namespace

	if desiredApp.Annotations == nil {
		desiredApp.Annotations = map[string]string{}
	}
	desiredApp.Annotations["packaging.carvel.dev/package-ref-name"] = pkgVersion.Spec.RefName
	desiredApp.Annotations["packaging.carvel.dev/package-version"] = pkgVersion.Spec.Version

	desiredApp.Spec = *pkgVersion.Spec.Template.Spec
	desiredApp.Spec.ServiceAccountName = pkgInstall.Spec.ServiceAccountName
	if pkgInstall.Spec.SyncPeriod == nil {
		desiredApp.Spec.SyncPeriod = &metav1.Duration{Duration: opts.DefaultSyncPeriod}
	} else {
		desiredApp.Spec.SyncPeriod = pkgInstall.Spec.SyncPeriod
	}
	desiredApp.Spec.NoopDelete = pkgInstall.Spec.NoopDelete
	desiredApp.Spec.Paused = pkgInstall.Spec.Paused
	desiredApp.Spec.Canceled = pkgInstall.Spec.Canceled
	desiredApp.Spec.Cluster = pkgInstall.Spec.Cluster
	desiredApp.Spec.DefaultNamespace = pkgInstall.Spec.DefaultNamespace

	err := controllerutil.SetControllerReference(pkgInstall, desiredApp, scheme.Scheme)
	if err != nil {
		return &v1alpha1.App{}, err
	}

	for i, fetchStep := range desiredApp.Spec.Fetch {
		annKey := fmt.Sprintf(ExtFetchSecretNameAnnKeyFmt, i)

		secretName, found := pkgInstall.Annotations[annKey]
		if !found {
			continue
		}

		secretRef := &kcv1alpha1.AppFetchLocalRef{Name: secretName}
		switch {
		case fetchStep.Inline != nil:
			// do nothing
		case fetchStep.Image != nil:
			desiredApp.Spec.Fetch[i].Image.SecretRef = secretRef
		case fetchStep.HTTP != nil:
			desiredApp.Spec.Fetch[i].HTTP.SecretRef = secretRef
		case fetchStep.Git != nil:
			desiredApp.Spec.Fetch[i].Git.SecretRef = secretRef
		case fetchStep.HelmChart != nil:
			if desiredApp.Spec.Fetch[i].HelmChart.Repository != nil {
				desiredApp.Spec.Fetch[i].HelmChart.Repository.SecretRef = secretRef
			}
		case fetchStep.ImgpkgBundle != nil:
			desiredApp.Spec.Fetch[i].ImgpkgBundle.SecretRef = secretRef
		default:
			// do nothing
		}
	}

	templatesPatcher := templateStepsPatcher{
		yttPatcher: &yttStepPatcher{
			addValuesAsInlinePaths: pkgiHasAnnotation(pkgInstall, ExtYttDataValuesOverlaysAnnKey),
			additionalPaths:        secretNamesFromAnn(pkgInstall, ExtYttPathsFromSecretNameAnnKey),
		},
		helmPatcher: &helmStepPatcher{
			additionalPaths: secretNamesFromAnn(pkgInstall, ExtHelmPathsFromSecretNameAnnKey),
			name:            pkgiAnnotationValue(pkgInstall, HelmTemplateOverlayNameKey),
			namespace:       pkgiAnnotationValue(pkgInstall, HelmTemplateOverlayNameSpaceKey),
		},
		cuePatcher: &cueStepPatcher{},

		templateSteps: desiredApp.Spec.Template,
		values:        pkgInstall.Spec.Values,
	}

	if err := templatesPatcher.patch(); err != nil {
		return &v1alpha1.App{}, err
	}

	return desiredApp, nil
}

type stepClass string

const (
	// anything that can take values
	stepClassValueable stepClass = "valueable"
	// only helm template steps
	stepClassHelm stepClass = "helm"
	// only ytt template steps
	stepClassYtt stepClass = "ytt"
	// only cue template steps
	stepClassCue stepClass = "cue"
)

type yttStepPatcher struct {
	addValuesAsInlinePaths bool     // TODO: support multiple ytt steps
	additionalPaths        []string // TODO: support multiple ytt steps
}

func (yp *yttStepPatcher) addValues(yttStep *kcv1alpha1.AppTemplateYtt, value pkgingv1alpha1.PackageInstallValues) {
	if yp.addValuesAsInlinePaths {
		addSecretAsInlinePath(&yttStep.Inline, value.SecretRef.Name)
	} else {
		addSecretAsValueSource(&yttStep.ValuesFrom, value.SecretRef.Name)
	}
}

func (yp *yttStepPatcher) addPaths(yttStep *kcv1alpha1.AppTemplateYtt) {
	for _, secretName := range yp.additionalPaths {
		addSecretAsInlinePath(&yttStep.Inline, secretName)
	}
}

type helmStepPatcher struct {
	additionalPaths []string
	name            string // TODO: support multiple helm steps
	namespace       string // TODO: support multiple helm steps
}

func (hp *helmStepPatcher) addValues(helmStep *kcv1alpha1.AppTemplateHelmTemplate, value pkgingv1alpha1.PackageInstallValues) {
	addSecretAsValueSource(&helmStep.ValuesFrom, value.SecretRef.Name)
}

func (hp *helmStepPatcher) addPaths(helmStep *kcv1alpha1.AppTemplateHelmTemplate) {
	for _, secretName := range hp.additionalPaths {
		addSecretAsValueSource(&helmStep.ValuesFrom, secretName)
	}
}

func (hp *helmStepPatcher) setNameAndNamespace(helmStep *kcv1alpha1.AppTemplateHelmTemplate) {
	if hp.name != "" {
		helmStep.Name = hp.name
	}
	if hp.namespace != "" {
		helmStep.Namespace = hp.namespace
	}
}

type cueStepPatcher struct{}

func (cp *cueStepPatcher) addValues(cueStep *kcv1alpha1.AppTemplateCue, value pkgingv1alpha1.PackageInstallValues) {
	addSecretAsValueSource(&cueStep.ValuesFrom, value.SecretRef.Name)
}

type templateStepsPatcher struct {
	templateSteps []kcv1alpha1.AppTemplate
	values        []pkgingv1alpha1.PackageInstallValues

	yttPatcher  *yttStepPatcher
	helmPatcher *helmStepPatcher
	cuePatcher  *cueStepPatcher

	classifiedSteps [][]stepClass
	once            sync.Once
}

func (p *templateStepsPatcher) classifySteps() {
	p.classifiedSteps = make([][]stepClass, len(p.templateSteps))

	for i, step := range p.templateSteps {
		classes := []stepClass{}

		if step.HelmTemplate != nil {
			classes = append(classes, stepClassHelm, stepClassValueable)
		}
		if step.Ytt != nil {
			classes = append(classes, stepClassYtt, stepClassValueable)
		}
		if step.Cue != nil {
			classes = append(classes, stepClassCue, stepClassValueable)
		}

		p.classifiedSteps[i] = classes
	}
}

func (p *templateStepsPatcher) stepHasClass(stepIdx int, class stepClass) bool {
	p.once.Do(p.classifySteps)

	for _, stepClass := range p.classifiedSteps[stepIdx] {
		if stepClass == class {
			return true
		}
	}

	return false
}

func (p *templateStepsPatcher) getClassifiedSteps(class stepClass) []int {
	p.once.Do(p.classifySteps)

	steps := []int{}
	for i := range p.classifiedSteps {
		if p.stepHasClass(i, class) {
			steps = append(steps, i)
		}
	}

	return steps
}

func (p *templateStepsPatcher) firstOf(class stepClass) (int, bool) {
	classifiedSteps := p.getClassifiedSteps(class)
	if len(classifiedSteps) < 1 {
		return 0, false
	}

	return classifiedSteps[0], true
}

func (p *templateStepsPatcher) defaultStepIdxs(stepIdxs []int, class stepClass) ([]int, error) {
	if len(stepIdxs) > 0 {
		return stepIdxs, nil
	}

	first, ok := p.firstOf(class)
	if ok {
		return []int{first}, nil
	}

	return []int{}, fmt.Errorf("no template step of class '%s' found", class)
}

func (p *templateStepsPatcher) patch() error {
	for _, values := range p.values {
		stepIdxs, err := p.defaultStepIdxs(values.TemplateSteps, stepClassValueable)
		if err != nil {
			return err
		}

		for _, stepIdx := range stepIdxs {
			if stepIdx < 0 || stepIdx >= len(p.templateSteps) {
				return fmt.Errorf("template step %d out of range", stepIdx)
			}
			if !p.stepHasClass(stepIdx, stepClassValueable) {
				return fmt.Errorf("template step %d does not support values", stepIdx)
			}

			templateStep := p.templateSteps[stepIdx]

			switch {
			case p.stepHasClass(stepIdx, stepClassYtt):
				p.yttPatcher.addValues(templateStep.Ytt, values)

			case p.stepHasClass(stepIdx, stepClassHelm):
				p.helmPatcher.addValues(templateStep.HelmTemplate, values)

			case p.stepHasClass(stepIdx, stepClassCue):
				p.cuePatcher.addValues(templateStep.Cue, values)
			}
		}
	}

	for _, stepIdx := range p.getClassifiedSteps(stepClassYtt) {
		p.yttPatcher.addPaths(p.templateSteps[stepIdx].Ytt)
		break // TODO: support multiple ytt steps
	}
	for _, stepIdx := range p.getClassifiedSteps(stepClassHelm) {
		ts := p.templateSteps[stepIdx].HelmTemplate
		p.helmPatcher.addPaths(ts)
		p.helmPatcher.setNameAndNamespace(ts)
		break // TODO: support multiple helm steps
	}

	return nil
}

func secretNamesFromAnn(installedPkg *pkgingv1alpha1.PackageInstall, annKey string) []string {
	var suffixes []string
	suffixToSecretName := map[string]string{}

	for ann, secretName := range installedPkg.Annotations {
		if ann == annKey {
			suffix := ""
			suffixToSecretName[suffix] = secretName
			suffixes = append(suffixes, suffix)
		} else if strings.HasPrefix(ann, annKey+".") {
			suffix := strings.TrimPrefix(ann, annKey+".")
			suffixToSecretName[suffix] = secretName
			suffixes = append(suffixes, suffix)
		}
	}

	sort.Strings(suffixes)

	var result []string
	for _, suffix := range suffixes {
		result = append(result, suffixToSecretName[suffix])
	}
	return result
}

func pkgiAnnotationValue(pkgi *pkgingv1alpha1.PackageInstall, key string) string {
	if anno, found := pkgi.Annotations[key]; found {
		return anno
	}
	return ""
}

func pkgiHasAnnotation(pkgi *pkgingv1alpha1.PackageInstall, key string) bool {
	_, found := pkgi.Annotations[key]
	return found
}

// addSecretAsInlinePath adds a secret as an inline path to the provided inline
// fetches. If the inline fetch is nil, it is initialized.
func addSecretAsInlinePath(inline **kcv1alpha1.AppFetchInline, secretName string) {
	if *inline == nil {
		*inline = &kcv1alpha1.AppFetchInline{}
	}
	(*inline).PathsFrom = append((*inline).PathsFrom, kcv1alpha1.AppFetchInlineSource{
		SecretRef: &kcv1alpha1.AppFetchInlineSourceRef{
			Name: secretName,
		},
	})
}

// addSecretAsValueSource adds a secret as a value source to the provided
// template values sources.
func addSecretAsValueSource(values *[]kcv1alpha1.AppTemplateValuesSource, secretName string) {
	*values = append(*values, kcv1alpha1.AppTemplateValuesSource{
		SecretRef: &kcv1alpha1.AppTemplateValuesSourceRef{
			Name: secretName,
		},
	})
}
