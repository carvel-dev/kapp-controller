// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	kcinstall "carvel.dev/kapp-controller/pkg/apis/kappctrl/install"
	"carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkginginstall "carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/install"
	"carvel.dev/kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	"carvel.dev/kapp-controller/pkg/apiserver/openapi"
	packagerest "carvel.dev/kapp-controller/pkg/apiserver/registry/datapackaging"
	kcclient "carvel.dev/kapp-controller/pkg/client/clientset/versioned"
	"github.com/carvel-dev/semver/v4"
	"github.com/go-logr/logr"
	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/admission/plugin/namespace/lifecycle"
	mutatingadmissionpolicy "k8s.io/apiserver/pkg/admission/plugin/policy/mutating"
	validatingadmissionpolicy "k8s.io/apiserver/pkg/admission/plugin/policy/validating"
	mutatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	validatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
	genericopenapi "k8s.io/apiserver/pkg/endpoints/openapi"
	apirest "k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	apiservercompatibility "k8s.io/apiserver/pkg/util/compatibility"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aggregatorclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	genericopenapicommon "k8s.io/kube-openapi/pkg/common"
	openapispec "k8s.io/kube-openapi/pkg/validation/spec"
)

const (
	// selfSignedCertDir is the dir kapp-controller self signed certificates are created in.
	selfSignedCertDir = "/home/kapp-controller/kc-agg-api-selfsigned-certs"

	TokenPath = "/token-dir"

	kappctrlNSEnvKey  = "KAPPCTRL_SYSTEM_NAMESPACE"
	kappctrlSVCEnvKey = "KAPPCTRL_SYSTEM_SERVICE"

	apiServiceName = "v1alpha1.data.packaging.carvel.dev"

	apiServiceReconcileInterval = 30 * time.Second
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

// getOpenAPIDefinitions wraps the generated definitions and injects
// k8s.io/apimachinery/pkg/version.Info, which is required by the /version
// endpoint that the generic API server always installs. Without it,
// routes.InstallV2 calls klog.Fatalf and crashes the process.
func getOpenAPIDefinitions(ref genericopenapicommon.ReferenceCallback) map[string]genericopenapicommon.OpenAPIDefinition {
	defs := openapi.GetOpenAPIDefinitions(ref)
	// version.Info implements OpenAPIModelNamer, returning the dot-notation name
	// "io.k8s.apimachinery.pkg.version.Info", so the map key must match that.
	defs["io.k8s.apimachinery.pkg.version.Info"] = genericopenapicommon.OpenAPIDefinition{
		Schema: openapispec.Schema{
			SchemaProps: openapispec.SchemaProps{
				Description: "Info contains versioning information.",
				Type:        []string{"object"},
				Properties: map[string]openapispec.Schema{
					"major":                 {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"minor":                 {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"emulationMajor":        {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"emulationMinor":        {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"minCompatibilityMajor": {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"minCompatibilityMinor": {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"gitVersion":            {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"gitCommit":             {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"gitTreeState":          {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"buildDate":             {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"goVersion":             {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"compiler":              {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
					"platform":              {SchemaProps: openapispec.SchemaProps{Type: []string{"string"}, Format: ""}},
				},
				Required: []string{"major", "minor", "gitVersion", "gitCommit", "gitTreeState", "buildDate", "goVersion", "compiler", "platform"},
			},
		},
	}
	return defs
}

func init() {
	// Setup the scheme the server will use
	datapkginginstall.Install(Scheme)
	kcinstall.Install(Scheme)
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type APIServer struct {
	server    *genericapiserver.GenericAPIServer
	stopCh    chan struct{}
	aggClient aggregatorclient.Interface
	logger    logr.Logger
}

// NewAPIServerOpts is a collection of scalar arguments for the NewAPIServer function
type NewAPIServerOpts struct {
	// GlobalNamespace sets the special namespace that kc will always check,
	// so things can be installed to either the ns you specify or this special global ns
	GlobalNamespace string
	// BindPort is the port on which to serve HTTPS with authentication and authorization
	BindPort int
	// EnableAPIPriorityAndFairness sets a featuregate to allow us backwards compatibility with
	// v1.19 and earlier clusters - our libraries use the beta version of those APIs but they used to be alpha.
	EnableAPIPriorityAndFairness bool

	// TLSCipherSuites is the list of cipher suites the api server will be willing to use. Empty list defaults to the underlying
	// libraries' defaults, which is usually fine especially if you don't expose the APIServer outside the cluster.
	// see also: https://golang.org/pkg/crypto/tls/#pkg-constants
	// According to Antrea, who we mostly copied:
	// Note that TLS1.3 Cipher Suites cannot be added to the list. But the apiserver will always
	// prefer TLS1.3 Cipher Suites whenever possible.
	TLSCipherSuites []string

	// Logger is a logger
	Logger logr.Logger
}

func NewAPIServer(clientConfig *rest.Config, coreClient kubernetes.Interface, kcClient kcclient.Interface, opts NewAPIServerOpts) (*APIServer, error) { //nolint
	aggClient, err := aggregatorclient.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("building aggregation client: %v", err)
	}

	config, caContentProvider, err := newServerConfig(aggClient, coreClient, opts)
	if err != nil {
		return nil, err
	}

	server, err := config.Complete().New("kapp-controller-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	// Register the PostStartHook to reconcile the CA Bundle
	if err := server.AddPostStartHook("apiservice-ca-reconciler", func(hookContext genericapiserver.PostStartHookContext) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := updateAPIService(ctx, opts.Logger, aggClient, caContentProvider); err != nil {
			opts.Logger.Error(err, "Initial APIService CA sync failed")
			return err
		}

		// Background Reconciliation
		go wait.Until(func() {
			ctx, syncCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer syncCancel()
			if err := updateAPIService(ctx, opts.Logger, aggClient, caContentProvider); err != nil {
				opts.Logger.Error(err, "Background APIService CA reconciliation failed")
			}
		}, apiServiceReconcileInterval, hookContext.Done())

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error registering APIService CA reconciler hook: %v", err)
	}

	packageMetadatasStorage := packagerest.NewPackageMetadataCRDREST(kcClient, coreClient, opts.GlobalNamespace)
	packageStorage := packagerest.NewPackageCRDREST(kcClient, coreClient, opts.GlobalNamespace, opts.Logger)

	pkgGroup := genericapiserver.NewDefaultAPIGroupInfo(datapackaging.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	pkgv1alpha1Storage := map[string]apirest.Storage{}
	pkgv1alpha1Storage["packagemetadatas"] = packageMetadatasStorage
	pkgv1alpha1Storage["packages"] = packageStorage
	pkgGroup.VersionedResourcesStorageMap["v1alpha1"] = pkgv1alpha1Storage

	err = server.InstallAPIGroup(&pkgGroup)
	if err != nil {
		return nil, err
	}

	return &APIServer{server, make(chan struct{}), aggClient, opts.Logger}, nil
}

// Run spawns a go routine that exits when apiserver is stopped.
func (as *APIServer) Run() error {
	go func() {
		err := as.server.PrepareRun().Run(as.stopCh)
		if err != nil {
			as.logger.Error(err, "API service stopped")
		}
	}()

	return wait.PollInfinite(time.Second, func() (bool, error) {
		as.logger.Info("waiting for API service to become ready. Check the status by running `kubectl get apiservices v1alpha1.data.packaging.carvel.dev -o yaml`")
		return as.isReady()
	})
}

func (as *APIServer) Stop() {
	close(as.stopCh)
}

func (as *APIServer) isReady() (bool, error) {
	apiService, err := as.aggClient.ApiregistrationV1().APIServices().Get(context.TODO(), apiServiceName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("error getting APIService %s: %v", apiServiceName, err)
	}

	for _, condition := range apiService.Status.Conditions {
		if condition.Type == apiregv1.Available {
			return condition.Status == apiregv1.ConditionTrue, nil
		}
	}

	return false, nil
}

func newServerConfig(aggClient aggregatorclient.Interface, coreClient kubernetes.Interface, opts NewAPIServerOpts) (*genericapiserver.RecommendedConfig, *dynamiccertificates.DynamicFileCAContent, error) {

	recommendedOptions := genericoptions.NewRecommendedOptions("", Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion))
	recommendedOptions.Etcd = nil

	// Set the PairName and CertDirectory to generate the certificate files.
	recommendedOptions.SecureServing.ServerCert.CertDirectory = selfSignedCertDir
	recommendedOptions.SecureServing.ServerCert.PairName = "kapp-controller"
	recommendedOptions.SecureServing.CipherSuites = opts.TLSCipherSuites

	// ports below 1024 are probably the wrong port, see https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Well-known_ports
	if opts.BindPort < 1024 {
		return nil, nil, fmt.Errorf("error initializing API Port to %v - try passing a port above 1023", opts.BindPort)
	}
	recommendedOptions.SecureServing.BindPort = opts.BindPort

	if err := recommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("kapp-controller", []string{apiServiceEndoint()}, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	caContentProvider, err := dynamiccertificates.NewDynamicCAContentFromFile("self-signed cert", recommendedOptions.SecureServing.ServerCert.CertKey.CertFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading self-signed CA certificate: %v", err)
	}

	serverVersion, err := getServerVersion(aggClient.Discovery())
	if err != nil {
		return nil, nil, err
	}

	// this feature gate is not enabled in k8s <1.29 as the
	// APIs it relies on were in v1beta3/v1beta2/v1beta1/alpha.
	// the apiserver library hardcodes the v1 version of the resource
	// so the best we can do for older k8s clusters is to allow it to be disabled.
	minSupportedVersionForAPF, err := semver.New("1.29.0")
	if err != nil {
		return nil, nil, err
	}
	isServerVerLTminSupportedVer := serverVersion.LT(*minSupportedVersionForAPF)
	if !opts.EnableAPIPriorityAndFairness || isServerVerLTminSupportedVer {
		if isServerVerLTminSupportedVer {
			opts.Logger.Info("The current version of kapp-controller does not support api-priority-and-fairness for versions of kubernets prior to 1.29, disabling this option")
		}
		recommendedOptions.Features.EnablePriorityAndFairness = false
	}

	// Always explicitly set the admission plugin order. Only include an admission
	// policy plugin if BOTH of these hold:
	//   - Its backing API resource is served by the cluster (some plugins are
	//     gated by feature flags — e.g. MutatingAdmissionPolicy defaults on only
	//     from k8s 1.36 — and their resources are absent on older clusters).
	//   - The ServiceAccount kapp-controller runs under can list+watch that
	//     resource. If discovery says the resource exists but our RBAC does not
	//     grant list+watch, loading the plugin causes its SharedInformer to
	//     never sync. WaitForReady() then never returns true and every mutating
	//     admission decision fails with "not yet ready to handle request
	//     (reason: Forbidden)" — a very hard failure mode to diagnose because
	//     the APIService still reports Available=True. Probing RBAC up-front
	//     lets us skip the plugin cleanly and log an actionable warning
	//     instead.
	pluginOrder := []string{lifecycle.PluginName, mutatingwebhook.PluginName, validatingwebhook.PluginName}
	if canUseAdmissionPolicyPlugin(aggClient.Discovery(), coreClient, "validatingadmissionpolicies", "validatingadmissionpolicybindings", opts.Logger) {
		pluginOrder = append(pluginOrder, validatingadmissionpolicy.PluginName)
	}
	if canUseAdmissionPolicyPlugin(aggClient.Discovery(), coreClient, "mutatingadmissionpolicies", "mutatingadmissionpolicybindings", opts.Logger) {
		pluginOrder = append(pluginOrder, mutatingadmissionpolicy.PluginName)
	}
	recommendedOptions.Admission.RecommendedPluginOrder = pluginOrder

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)

	// EffectiveVersion must be initialized for k8s.io/apiserver v0.36.0+;
	// config.Complete() dereferences it and panics if nil.
	serverConfig.EffectiveVersion = apiservercompatibility.DefaultBuildEffectiveVersion()

	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(
		getOpenAPIDefinitions,
		genericopenapi.NewDefinitionNamer(Scheme))
	serverConfig.OpenAPIV3Config.Info.Title = "Kapp-controller"
	serverConfig.OpenAPIV3Config.Info.Version = "v1alpha1"

	// Register OpenAPI v2 configuration to satisfy the Kubernetes Aggregation Controller.
	// In K8s 1.30+, the aggregator syncs both v2 and v3 specs; providing v2 prevents
	// "resource not found" errors in the kube-apiserver logs (Fixes #1703).
	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		getOpenAPIDefinitions,
		genericopenapi.NewDefinitionNamer(Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "Kapp-controller"
	serverConfig.OpenAPIConfig.Info.Version = "v1alpha1"

	if err := recommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, nil, err
	}

	return serverConfig, caContentProvider, nil
}

// admissionPolicyAPIGroupVersion is the API group/version the vendored admission
// policy plugins operate against. Alpha/beta groupversions are intentionally not
// probed: the vendored plugin implementations only support the v1 shape, so on
// clusters where only alpha/beta is served we treat the plugin as unavailable.
const admissionPolicyAPIGroupVersion = "admissionregistration.k8s.io/v1"

// canUseAdmissionPolicyPlugin reports whether it is safe to load the admission
// policy plugin identified by (policyResource, bindingResource) into the
// aggregated apiserver's admission chain.
//
// The plugin is safe to load only if:
//  1. Both resources are present in cluster discovery under
//     admissionPolicyAPIGroupVersion (skips gracefully on older k8s or when
//     the plugin's feature gate is disabled).
//  2. The ServiceAccount this aggregated apiserver runs under is allowed to
//     list AND watch both resources at cluster scope. This is required
//     because the plugin's own SharedInformer performs a List+Watch at
//     startup; if List returns 403 the informer's cache never syncs and
//     every mutating admission decision that plugin evaluates will
//     permanently fail with "not yet ready to handle request".
//
// Discovery or SubjectAccessReview errors are treated conservatively as
// "cannot use" and cause the plugin to be skipped with a clear warning log.
// The aggregated apiserver's own startup is never blocked by this check.
func canUseAdmissionPolicyPlugin(discoveryClient discovery.DiscoveryInterface, coreClient kubernetes.Interface, policyResource, bindingResource string, logger logr.Logger) bool {
	for _, resource := range []string{policyResource, bindingResource} {
		if !serverResourceExists(discoveryClient, admissionPolicyAPIGroupVersion, resource) {
			logger.Info("Skipping admission-policy plugin: resource not present in cluster discovery",
				"apiGroupVersion", admissionPolicyAPIGroupVersion, "resource", resource)
			return false
		}
		for _, verb := range []string{"list", "watch"} {
			allowed, err := selfSubjectAccessAllowed(coreClient, "admissionregistration.k8s.io", resource, verb)
			if err != nil {
				logger.Error(err, "Skipping admission-policy plugin: SelfSubjectAccessReview failed",
					"resource", resource, "verb", verb)
				return false
			}
			if !allowed {
				logger.Info("Skipping admission-policy plugin: ServiceAccount lacks required permission. "+
					"Loading this plugin without list+watch permission would cause its informer cache to never sync, "+
					"resulting in every mutating admission request being rejected with \"not yet ready to handle request\". "+
					"Grant list+watch on this resource (in group admissionregistration.k8s.io) to the kapp-controller ClusterRole to enable this plugin.",
					"resource", resource, "verb", verb)
				return false
			}
		}
	}
	return true
}

// selfSubjectAccessAllowed checks whether the caller is permitted to perform
// verb on the cluster-scoped resource in group. It uses SelfSubjectAccessReview
// so no additional RBAC is required beyond the subjectaccessreviews/create
// permission already granted to kapp-controller.
func selfSubjectAccessAllowed(coreClient kubernetes.Interface, group, resource, verb string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	review, err := coreClient.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Group:    group,
				Resource: resource,
				Verb:     verb,
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return review.Status.Allowed, nil
}

// serverResourceExists returns true if the given groupVersion + resource kind is
// available in the cluster. It is used to guard admission plugins whose backing
// resources may not be present (e.g. MutatingAdmissionPolicy behind a feature gate).
func serverResourceExists(discoveryClient discovery.DiscoveryInterface, groupVersion, resource string) bool {
	list, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return false
	}
	for _, r := range list.APIResources {
		if r.Name == resource {
			return true
		}
	}
	return false
}

func getServerVersion(discoveryClient discovery.DiscoveryInterface) (semver.Version, error) {
	version, err := discoveryClient.ServerVersion()
	if err != nil {
		return semver.Version{}, err
	}
	retv, err := semver.ParseTolerant(version.String())
	if err != nil {
		return retv, err
	}
	retv.Pre = semver.PRVersion{}
	retv.Build = semver.BuildMeta{}
	return retv, nil
}

func updateAPIService(ctx context.Context, logger logr.Logger, client aggregatorclient.Interface, caProvider dynamiccertificates.CAContentProvider) error {
	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		apiService, err := client.ApiregistrationV1().APIServices().Get(ctx, apiServiceName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting APIService %s: %v", apiServiceName, err)
		}

		caBundle := caProvider.CurrentCABundleContent()
		if bytes.Equal(apiService.Spec.CABundle, caBundle) {
			return nil
		}

		logger.Info("Syncing CA certificate with APIServices")
		apiService.Spec.CABundle = caBundle
		_, err = client.ApiregistrationV1().APIServices().Update(ctx, apiService, metav1.UpdateOptions{})
		return err
	}); err != nil {
		return fmt.Errorf("error updating kapp-controller CA cert of APIService %s: %v", apiServiceName, err)
	}
	return nil
}

func apiServiceEndoint() string {
	var apiServiceName = getEnvVal(kappctrlSVCEnvKey, "packaging-api")
	ns := os.Getenv(kappctrlNSEnvKey)
	if ns == "" {
		panic("Cannot get api service endpoint, Kapp-controller namespace is empty")
	}

	return fmt.Sprintf("%s.%s.svc", apiServiceName, ns)
}

func getEnvVal(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
