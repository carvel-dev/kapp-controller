// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
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
	"github.com/go-logr/logr"
	"github.com/k14s/semver/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	genericopenapi "k8s.io/apiserver/pkg/endpoints/openapi"
	apirest "k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aggregatorclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

const (
	// selfSignedCertDir is the dir kapp-controller self signed certificates are created in.
	selfSignedCertDir = "/home/kapp-controller/kc-agg-api-selfsigned-certs"

	TokenPath = "/token-dir"

	kappctrlNSEnvKey  = "KAPPCTRL_SYSTEM_NAMESPACE"
	kappctrlSVCEnvKey = "KAPPCTRL_SYSTEM_SERVICE"

	apiServiceName = "v1alpha1.data.packaging.carvel.dev"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

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

	config, err := newServerConfig(aggClient, opts)
	if err != nil {
		return nil, err
	}

	server, err := config.Complete().New("kapp-controller-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
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

func newServerConfig(aggClient aggregatorclient.Interface, opts NewAPIServerOpts) (*genericapiserver.RecommendedConfig, error) {
	recommendedOptions := genericoptions.NewRecommendedOptions("", Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion))
	recommendedOptions.Etcd = nil

	// Set the PairName and CertDirectory to generate the certificate files.
	recommendedOptions.SecureServing.ServerCert.CertDirectory = selfSignedCertDir
	recommendedOptions.SecureServing.ServerCert.PairName = "kapp-controller"
	recommendedOptions.SecureServing.CipherSuites = opts.TLSCipherSuites

	// ports below 1024 are probably the wrong port, see https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Well-known_ports
	if opts.BindPort < 1024 {
		return nil, fmt.Errorf("error initializing API Port to %v - try passing a port above 1023", opts.BindPort)
	}
	recommendedOptions.SecureServing.BindPort = opts.BindPort

	if err := recommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("kapp-controller", []string{apiServiceEndoint()}, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	caContentProvider, err := dynamiccertificates.NewDynamicCAContentFromFile("self-signed cert", recommendedOptions.SecureServing.ServerCert.CertKey.CertFile)
	if err != nil {
		return nil, fmt.Errorf("error reading self-signed CA certificate: %v", err)
	}

	if err := updateAPIService(opts.Logger, aggClient, caContentProvider); err != nil {
		return nil, fmt.Errorf("error updating api service with generated certs: %v", err)
	}

	serverVersion, err := getServerVersion(aggClient.Discovery())
	if err != nil {
		return nil, err
	}
	// this feature gate is not enabled in k8s <1.26 as the
	// APIs it relies on were in v1beta2/v1beta1/alpha.
	// the apiserver library hardcodes the v1beta3 version of the resource
	// so the best we can do for older k8s clusters is to allow it to be disabled.
	minSupportedVersionForAPF, err := semver.New("1.26.0")
	if err != nil {
		return nil, err
	}
	isServerVerLTminSupportedVer := serverVersion.LT(*minSupportedVersionForAPF)
	if !opts.EnableAPIPriorityAndFairness || isServerVerLTminSupportedVer {
		if isServerVerLTminSupportedVer {
			opts.Logger.Info("The current version of kapp-controller does not support api-priority-and-fairness for versions of kubernets prior to 1.26, disabling this option")
		}
		err := feature.DefaultMutableFeatureGate.Set("APIPriorityAndFairness=false")
		if err != nil {
			return nil, fmt.Errorf("error updating disabling feature gate for APIPriorityAndFairness: %v", err)
		}
	}

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	if err := recommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(
		openapi.GetOpenAPIDefinitions,
		genericopenapi.NewDefinitionNamer(Scheme))
	serverConfig.OpenAPIV3Config.Info.Title = "Kapp-controller"
	return serverConfig, nil
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

func updateAPIService(logger logr.Logger, client aggregatorclient.Interface, caProvider dynamiccertificates.CAContentProvider) error {
	logger.Info("Syncing CA certificate with APIServices")
	apiService, err := client.ApiregistrationV1().APIServices().Get(context.TODO(), apiServiceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting APIService %s: %v", apiServiceName, err)
	}
	apiService.Spec.CABundle = caProvider.CurrentCABundleContent()
	if _, err := client.ApiregistrationV1().APIServices().Update(context.TODO(), apiService, metav1.UpdateOptions{}); err != nil {
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
