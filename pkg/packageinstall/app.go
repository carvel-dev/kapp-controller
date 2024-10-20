// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	kcv1alpha1 "carvel.dev/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	pkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apis/packaging/v1alpha1"
	datapkgingv1alpha1 "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"carvel.dev/kapp-controller/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	ManuallyControlledAnnKey = "ext.packaging.carvel.dev/manually-controlled"

	HelmTemplateOverlayNameKey         = "ext.packaging.carvel.dev/helm-template-name"
	HelmTemplateOverlayNameKeyFmt      = "ext.packaging.carvel.dev/helm-%d-template-name"
	HelmTemplateOverlayNameSpaceKey    = "ext.packaging.carvel.dev/helm-template-namespace"
	HelmTemplateOverlayNameSpaceKeyFmt = "ext.packaging.carvel.dev/helm-%d-template-namespace"

	// Resulting secret names are sorted deterministically by suffix
	ExtYttPathsFromSecretNameAnnKey     = "ext.packaging.carvel.dev/ytt-paths-from-secret-name"
	ExtYttPathsFromSecretNameAnnKeyFmt  = "ext.packaging.carvel.dev/ytt-%d-paths-from-secret-name"
	ExtHelmPathsFromSecretNameAnnKey    = "ext.packaging.carvel.dev/helm-template-values-from-secret-name"
	ExtHelmPathsFromSecretNameAnnKeyFmt = "ext.packaging.carvel.dev/helm-%d-template-values-from-secret-name"

	// ExtYttDataValuesOverlaysAnnKey if set, adds the pkgi's values secrets as overlays/paths, not as values, to the app
	ExtYttDataValuesOverlaysAnnKey    = "ext.packaging.carvel.dev/ytt-data-values-overlays"
	ExtYttDataValuesOverlaysAnnKeyFmt = "ext.packaging.carvel.dev/ytt-%d-data-values-overlays"

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
		templateSteps: desiredApp.Spec.Template,
		values:        pkgInstall.Spec.Values,
		annotations:   pkgInstall.Annotations,
	}
	if err := templatesPatcher.patch(); err != nil {
		return &v1alpha1.App{}, err
	}

	return desiredApp, nil
}

type (
	stepClass   string
	stepClasses = sets.Set[stepClass]
)

const (
	// anything that can take values
	stepClassTakesValues stepClass = "takesValues"
	// only helm template steps
	stepClassHelm stepClass = "helm"
	// only ytt template steps
	stepClassYtt stepClass = "ytt"
	// only cue template steps
	stepClassCue stepClass = "cue"
)

type templateStepsPatcher struct {
	templateSteps []kcv1alpha1.AppTemplate
	values        []pkgingv1alpha1.PackageInstallValues
	annotations   map[string]string

	classifiedSteps []stepClasses
	once            sync.Once
}

func (p *templateStepsPatcher) classifySteps() {
	p.classifiedSteps = make([]stepClasses, len(p.templateSteps))

	for i, step := range p.templateSteps {
		classes := stepClasses{}

		if step.HelmTemplate != nil {
			classes.Insert(stepClassHelm, stepClassTakesValues)
		}
		if step.Ytt != nil {
			classes.Insert(stepClassYtt, stepClassTakesValues)
		}
		if step.Cue != nil {
			classes.Insert(stepClassCue, stepClassTakesValues)
		}

		p.classifiedSteps[i] = classes
	}
}

func (p *templateStepsPatcher) stepHasClass(stepIdx int, class stepClass) bool {
	p.once.Do(p.classifySteps)
	return p.classifiedSteps[stepIdx].Has(class)
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
	if err := p.patchFromValues(); err != nil {
		return err
	}
	p.patchFromYttAnnotations()
	p.patchFromHelmAnnotations()

	return nil
}

// patchFromValues patches all template steps that take values with values from
// the packageInstall
func (p *templateStepsPatcher) patchFromValues() error {
	for _, values := range p.values {
		stepIdxs, err := p.defaultStepIdxs(values.TemplateSteps, stepClassTakesValues)
		if err != nil {
			return err
		}

		for _, stepIdx := range stepIdxs {
			if stepIdx < 0 || stepIdx >= len(p.templateSteps) {
				return fmt.Errorf("template step %d out of range", stepIdx)
			}
			if !p.stepHasClass(stepIdx, stepClassTakesValues) {
				return fmt.Errorf("template step %d does not support values", stepIdx)
			}

			templateStep := p.templateSteps[stepIdx]

			switch {
			case p.stepHasClass(stepIdx, stepClassYtt):
				// ytt is a bit special: when we find the indexed annotation (or the
				// naked one for the first ytt step) to apply the values as inline
				// paths on the app, we do so; else, by default, we apply the pkgi's
				// values as values on the app.
				valuesAsPath := false
				if firstYttStepIdx, ok := p.firstOf(stepClassYtt); ok && stepIdx == firstYttStepIdx {
					if _, ok := p.annotations[ExtYttDataValuesOverlaysAnnKey]; ok {
						valuesAsPath = true
					}
				}
				if _, ok := p.annotations[fmt.Sprintf(ExtYttDataValuesOverlaysAnnKeyFmt, stepIdx)]; ok {
					valuesAsPath = true
				}
				if valuesAsPath {
					addSecretAsInlinePath(&templateStep.Ytt.Inline, values.SecretRef.Name)
				} else {
					addSecretAsValueSource(&templateStep.Ytt.ValuesFrom, values.SecretRef.Name)
				}

			case p.stepHasClass(stepIdx, stepClassHelm):
				addSecretAsValueSource(&templateStep.HelmTemplate.ValuesFrom, values.SecretRef.Name)

			case p.stepHasClass(stepIdx, stepClassCue):
				addSecretAsValueSource(&templateStep.Cue.ValuesFrom, values.SecretRef.Name)
			}
		}
	}

	return nil
}

// patchFromYttAnnotations patches ytt template steps with values from
// annotations from the packageInstall
func (p *templateStepsPatcher) patchFromYttAnnotations() {
	firstYttIdx, hasYtt := p.firstOf(stepClassYtt)

	if !hasYtt {
		return
	}

	patcher := func(ts *kcv1alpha1.AppTemplateYtt, pathsAnno string) {
		for _, secretName := range secretNamesFromAnn(p.annotations, pathsAnno) {
			addSecretAsInlinePath(&ts.Inline, secretName)
		}
	}

	for _, stepIdx := range p.getClassifiedSteps(stepClassYtt) {
		ts := p.templateSteps[stepIdx].Ytt

		if stepIdx == firstYttIdx {
			// annotations that are not indexed are only applied to the first ytt
			// step, so that we are backwards compatible
			patcher(ts, ExtYttPathsFromSecretNameAnnKey)
		}

		patcher(ts, fmt.Sprintf(ExtYttPathsFromSecretNameAnnKeyFmt, stepIdx))
	}
}

// patchFromHelmAnnotations patches helm template steps with values from
// annotations from the packageInstall
func (p *templateStepsPatcher) patchFromHelmAnnotations() {
	firstHelmIdx, hasHelm := p.firstOf(stepClassHelm)

	if !hasHelm {
		return
	}

	patcher := func(ts *kcv1alpha1.AppTemplateHelmTemplate, nameAnno, namespaceAnno, pathsAnno string) {
		if name, ok := p.annotations[nameAnno]; ok && name != "" {
			ts.Name = name
		}
		if namespace, ok := p.annotations[namespaceAnno]; ok && namespace != "" {
			ts.Namespace = namespace
		}
		for _, secretName := range secretNamesFromAnn(p.annotations, pathsAnno) {
			addSecretAsValueSource(&ts.ValuesFrom, secretName)
		}
	}

	for _, stepIdx := range p.getClassifiedSteps(stepClassHelm) {
		ts := p.templateSteps[stepIdx].HelmTemplate

		if stepIdx == firstHelmIdx {
			// annotations that are not indexed are only applied to the first helm
			// step, so that we are backwards compatible
			patcher(ts, HelmTemplateOverlayNameKey, HelmTemplateOverlayNameSpaceKey, ExtHelmPathsFromSecretNameAnnKey)
		}

		patcher(ts,
			fmt.Sprintf(HelmTemplateOverlayNameKeyFmt, stepIdx),
			fmt.Sprintf(HelmTemplateOverlayNameSpaceKeyFmt, stepIdx),
			fmt.Sprintf(ExtHelmPathsFromSecretNameAnnKeyFmt, stepIdx),
		)
	}
}

func secretNamesFromAnn(annotations map[string]string, annKey string) []string {
	var suffixes []string
	suffixToSecretName := map[string]string{}

	for ann, secretName := range annotations {
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
