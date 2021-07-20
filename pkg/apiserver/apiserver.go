// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/openapi"

	kcinstall "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apis/kappctrl/install"
	kcclient "github.com/vmware-tanzu/carvel-kapp-controller/pkg/client/clientset/versioned"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging"
	datapkginginstall "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/install"
	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
	packagerest "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/registry/datapackaging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apirest "k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	genericopenapi "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/klog"
	apiregv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aggregatorclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

const (
	// selfSignedCertDir is the dir kapp-controller self signed certificates are created in.
	selfSignedCertDir = "/home/kapp-controller/kc-agg-api-selfsigned-certs"

	TokenPath = "/token-dir"

	kappctrlNSEnvKey = "KAPPCTRL_SYSTEM_NAMESPACE"
	apiServiceName   = "v1alpha1.data.packaging.carvel.dev"
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
}

func NewAPIServer(clientConfig *rest.Config, coreClient kubernetes.Interface, kcClient kcclient.Interface, globalNamespace string, bindPort int) (*APIServer, error) {
	aggClient, err := aggregatorclient.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("building aggregation client: %v", err)
	}

	config, err := newServerConfig(aggClient, bindPort)
	if err != nil {
		return nil, err
	}

	server, err := config.Complete().New("kapp-controller-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	packageMetadatasStorage := packagerest.NewPackageMetadataCRDREST(kcClient, coreClient, globalNamespace)
	packageStorage := packagerest.NewPackageCRDREST(kcClient, coreClient, globalNamespace)

	pkgGroup := genericapiserver.NewDefaultAPIGroupInfo(datapackaging.GroupName, Scheme, metav1.ParameterCodec, Codecs)
	pkgv1alpha1Storage := map[string]apirest.Storage{}
	pkgv1alpha1Storage["packagemetadatas"] = packageMetadatasStorage
	pkgv1alpha1Storage["Packages"] = packageStorage
	pkgGroup.VersionedResourcesStorageMap["v1alpha1"] = pkgv1alpha1Storage

	err = server.InstallAPIGroup(&pkgGroup)
	if err != nil {
		return nil, err
	}

	return &APIServer{server, make(chan struct{}), aggClient}, nil
}

// Run spawns a go routine that exits when apiserver is stopped.
func (as *APIServer) Run() error {
	const (
		retries = 30
	)
	go as.server.PrepareRun().Run(as.stopCh)

	for i := 0; i < retries; i++ {
		ready, err := as.isReady()
		if err != nil {
			return fmt.Errorf("checking readiness: %v", err)
		}

		if ready {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timed out after %s waiting for api server to become healthy", retries*time.Second)
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

func newServerConfig(aggClient aggregatorclient.Interface, bindPort int) (*genericapiserver.RecommendedConfig, error) {
	recommendedOptions := genericoptions.NewRecommendedOptions("", Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion))
	recommendedOptions.Etcd = nil

	// Set the PairName and CertDirectory to generate the certificate files.
	recommendedOptions.SecureServing.ServerCert.CertDirectory = selfSignedCertDir
	recommendedOptions.SecureServing.ServerCert.PairName = "kapp-controller"

	// ports below 1024 are probably the wrong port, see https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Well-known_ports
	if bindPort < 1024 {
		return nil, fmt.Errorf("error initializing API Port to %v - try passing a port above 1023", bindPort)
	}
	recommendedOptions.SecureServing.BindPort = bindPort

	if err := recommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("kapp-controller", []string{apiServiceEndoint()}, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	caContentProvider, err := dynamiccertificates.NewDynamicCAContentFromFile("self-signed cert", recommendedOptions.SecureServing.ServerCert.CertKey.CertFile)
	if err != nil {
		return nil, fmt.Errorf("error reading self-signed CA certificate: %v", err)
	}

	if err := updateAPIService(aggClient, caContentProvider); err != nil {
		return nil, fmt.Errorf("error updating api service with generated certs: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	if err := recommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		openapi.GetOpenAPIDefinitions,
		genericopenapi.NewDefinitionNamer(Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "Kapp-controller"
	return serverConfig, nil
}

func updateAPIService(client aggregatorclient.Interface, caProvider dynamiccertificates.CAContentProvider) error {
	klog.Info("Syncing CA certificate with APIServices")
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
	const apiServiceName = "packaging-api"
	ns := os.Getenv(kappctrlNSEnvKey)
	if ns == "" {
		panic("Cannot get api service endpoint, Kapp-controller namespace is empty")
	}

	return fmt.Sprintf("%s.%s.svc", apiServiceName, ns)
}
