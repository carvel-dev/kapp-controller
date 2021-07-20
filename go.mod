module github.com/vmware-tanzu/carvel-kapp-controller

go 1.16

require (
	cloud.google.com/go v0.60.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/fatih/color v1.12.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.4.0
	github.com/go-openapi/spec v0.19.3
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/go-containerregistry v0.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/onsi/ginkgo v1.16.1 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/vmware-tanzu/carvel-vendir v0.21.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/genproto v0.0.0-20200707001353-8e8330bf89df // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	honnef.co/go/tools v0.2.0 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/apiserver v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/code-generator v0.19.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.19.2
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/controller-tools v0.4.1
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
