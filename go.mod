module github.com/vmware-tanzu/carvel-kapp-controller

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/google/go-containerregistry v0.1.2
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/vmware-tanzu/carvel-vendir v0.19.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	k8s.io/code-generator v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/yaml v1.2.0
)
