module github.com/vmware-tanzu/carvel-kapp-controller/cli

go 1.16

require (
	github.com/cppforlife/cobrautil v0.0.0-20200514214827-bb86e6965d72
	github.com/cppforlife/go-cli-ui v0.0.0-20200716203538-1e47f820817f
	github.com/getkin/kin-openapi v0.81.0
	github.com/google/go-containerregistry v0.1.2
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/vmware-tanzu/carvel-kapp-controller v0.27.0
	github.com/vmware-tanzu/carvel-vendir v0.23.0
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	k8s.io/api v0.22.1 // kubernetes-1.22.1
	k8s.io/apimachinery v0.22.1 // kubernetes-1.22.1
	k8s.io/client-go v0.22.1 // kubernetes-1.22.1
)

replace github.com/spf13/cobra => github.com/spf13/cobra v1.1.1
