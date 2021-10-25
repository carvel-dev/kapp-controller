# Package

Manage package lifecycle operations.

## Usage

Package and Repository Operations ( Subject to Change ):

```sh
>>> tanzu package --help
Tanzu package management

Usage:
  tanzu package [command]

Available Commands:
    available   Manage available packages
    install     Install a package
    installed   Manage installed packages
    repository  Manage registered package repositories


Flags:
  -h, --help              help for package
      --log-file string   Log file path
  -v, --verbose int32     Number for the log level verbosity(0-9)

Use "tanzu package [command] --help" for more information about a command.
```

```sh
>>> tanzu package repository --help
Add, list, get or delete a repository for tanzu packages

Usage:
  tanzu package repository [command]

Available Commands:
  add         Add a repository
  delete      Delete a repository
  get         Get repository status
  list        List repository
  update      Update a repository

Flags:
  -h, --help   help for repository

Global Flags:
      --log-file string   Log file path
  -v, --verbose int32     Number for the log level verbosity(0-9)

Use "tanzu package repository [command] --help" for more information about a command.
```

## Test

1. Create a management cluster using latest tanzu cli

1. Use package commands to:
   * add a repository
   * list a repository
   * get a repository status
   * list packages
   * get a package information
   * get an installed package information
   * update a package
   * delete a repository

   Use the following image package bundles for testing:

   | S.no |                        Repository URL                                    |
   | :----| :------------------------------------------------------------------------|
   |  1.  |  projects-stg.registry.vmware.com/tkg/test-packages/test-repo:v1.0.0     |
   |  2.  |  projects-stg.registry.vmware.com/tkg/test-packages/standard-repo:v1.0.0 |

   Here is an example workflow

1. Add a repository

   ```sh
   >>> tanzu package repository add standard-repo --url projects-stg.registry.vmware.com/tkg/test-packages/standard-repo:v1.0.0 -n test-ns --create-namespace
   Added package repository 'standard-repo'
   ```

1. Get repository status

   ```sh
   >>> tanzu package repository get standard-repo -n test-ns
   NAME:        standard-repo
   VERSION:     88984
   REPOSITORY:  projects-stg.registry.vmware.com/tkg/test-packages/standard-repo:v1.0.0
   STATUS:      Reconcile succeeded
   REASON:
   ```

1. Update a repository

   ```sh
   >>> tanzu package repository update standard-repo --url projects-stg.registry.vmware.com/tkg/test-packages/standard-repo:v1.0.0 -n test-ns
   Updated package repository 'standard-repo' in namespace 'test-ns'
   ```

1. List the repository

   ```sh
   >>> tanzu package repository list -A
   NAME           REPOSITORY                                                               STATUS               DETAILS  NAMESPACE
   standard-repo  projects-stg.registry.vmware.com/tkg/test-packages/standard-repo:v1.0.0  Reconcile succeeded           test-ns
   repo           projects-stg.registry.vmware.com/tkg/test-packages/test-repo:v1.0.0      Reconcile succeeded           test-ns
   ```

1. Get information of a package

   Example 1: Get detailed information of a package

   ```sh
   >>> tanzu package available get contour.tanzu.vmware.com/1.15.1+vmware.1-tkg.1 --namespace test-ns
   / Retrieving package details for contour.tanzu.vmware.com/1.15.1+vmware.1-tkg.1...
     NAME:                           contour.tanzu.vmware.com
     VERSION:                        1.15.1+vmware.1-tkg.1
     RELEASED-AT:
     DISPLAY-NAME:                   contour
     SHORT-DESCRIPTION:              This package provides ingress functionality.
     PACKAGE-PROVIDER:
     MINIMUM-CAPACITY-REQUIREMENTS:
     LONG-DESCRIPTION:               This package provides ingress functionality.
     MAINTAINERS:                    []
     RELEASE-NOTES:
     LICENSE:                        []
   ```

   Example 2: Get openAPI schema of a package

   ```sh
   >>> tanzu package available get external-dns.tanzu.vmware.com/0.8.0+vmware.1-tkg.1 -n external-dns --values-schema
    KEY                         DEFAULT       TYPE     DESCRIPTION
    deployment.args             <nil>         array    List of arguments passed via command-line to external-dns.
                                                       For more guidance on configuration options for your
                                                       desired DNS provider, consult the ExternalDNS docs at
                                                       https://github.com/kubernetes-sigs/external-dns#running-externaldns

    deployment.env              <nil>         array    List of environment variables to set in the external-dns container.
    deployment.securityContext  <nil>         <nil>    SecurityContext defines the security options the external-dns container should be run with. More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
    deployment.volumeMounts     <nil>         array    Pod volumes to mount into the external-dns container's filesystem.
    deployment.volumes          <nil>         array    List of volumes that can be mounted by containers belonging to the external-dns pod. More info: https://kubernetes.io/docs/concepts/storage/volumes
    namespace                   external-dns  string   The namespace in which to deploy ExternalDNS.
   ```

1. Install a package

   Example 1: Install the specified version for package name "fluent-bit.tkg-standard.tanzu.vmware", while providing the values.yaml file and without waiting for package reconciliation to complete

   ```sh

   >>> tanzu package install fluentbit --package-name fluent-bit.tanzu.vmware.com --namespace test-ns --create-namespace --version 1.7.5+vmware.1-tkg.1 --values-file values.yaml --wait=false
   \ Installing package 'fluent-bit.tanzu.vmware.com'
   | Getting package metadata for fluent-bit.tanzu.vmware.com
   / Creating service account 'fluentbit-test-ns-sa'

   - Creating cluster role binding 'fluentbit-test-ns-cluster-rolebinding'

   Added installed package 'fluentbit' in namespace 'test-ns'
   ```

   An example values.yaml is as follows:

   ```yaml
   fluent_bit:
      config:
        outputs: |
          [OUTPUT]
            Name     stdout
            Match    *
    ```

    Example 2: Install the latest version for package name "contour.tanzu.vmware.com". If the namespace does not exist beforehand, it gets created.

    ```sh
    >>> tanzu package install contour-pkg --package-name contour.tanzu.vmware.com --namespace test-ns --version 1.15.1+vmware.1-tkg.1
    \ Installing package 'contour.tanzu.vmware.com'
    / Getting package metadata for contour.tanzu.vmware.com

    - Creating cluster admin role 'contour-pkg-test-ns-cluster-role'

    / Creating package resource
    / Package install status: Reconciling

     Added installed package 'contour-pkg' in namespace 'test-ns'
    ```

1. Get information of an installed package

   Example 1: Get information of an installed package

   ```sh
   >>> tanzu package installed get contour-pkg --namespace test-ns
   NAME:                 contour.tanzu.vmware.com
   PACKAGE-NAME:         contour-pkg
   PACKAGE-VERSION:      1.15.1+vmware.1-tkg.1
   STATUS:               Reconcile succeeded
   CONDITIONS:
   USEFUL-ERROR-MESSAGE:
   ```

   Example 2: Get data value secret of an installed package and save it to file (example: config.yaml)

   ```sh
   >>> tanzu package installed get fluent-bit --namespace test-ns --values-file config.yaml
   / Retrieving installation details for myfb...

   cat config.yaml
   fluent_bit:
     config:
       outputs: |
         [OUTPUT]
           Name     stdout
           Match    *
   ```

1. Update a package

   Example 1: Update a package with different version

   ```sh
   >>> tanzu package installed update mypkg --version 3.0.0-rc.1 --namespace test-ns
   | Updating package 'mypkg'
   / Getting package install for 'mypkg'
   - Getting package metadata for 'pkg.test.carvel.dev'

   Updated package install 'mypkg' in namespace 'test-ns'
   ```

   Example 2: Update a package which is not installed

   ```sh
   >>> tanzu package installed update fluent-bit --package-name fluent-bit.tanzu.vmware.com --version 1.7.5+vmware.1-tkg.1 --namespace test-ns --install
   / Getting package install for 'fluent-bit'

   - Getting package metadata for fluent-bit.tanzu.vmware.com
   \ Creating service account 'fluent-bit-test-ns-sa'

   | Creating cluster role binding 'fluent-bit-test-ns-cluster-rolebinding'
   - Creating package resource
   \ Package install status: Reconciling

   Updated package install 'fluent-bit' in namespace 'test-ns'
   ```

   Example 3: Update an installed package with providing values.yaml file

   ```sh
   >>> tanzu package installed update fluent-bit --version 1.7.5+vmware.1-tkg.1 --namespace test-ns --values-file values.yaml
   | Updating package 'fluent-bit'
   | Getting package install for 'fluent-bit'
   / Updating secret 'fluent-bit-test-ns-values'

   Updated package install 'fluent-bit' in namespace 'test-ns'
   ```

   An example values.yaml is as follows:

   ```yaml
   fluent_bit:
      config:
        outputs: |
          [OUTPUT]
            Name     stdout
            Match    /
   ```

1. Uninstall a package

   ```sh
   >>> tanzu package installed delete contour-pkg --namespace test-ns
   | Uninstalling package 'contour-pkg' from namespace 'test-ns'
   - Getting package install for 'contour-pkg'
   - Deleting package install 'contour-pkg' from namespace 'test-ns'
   | Package uninstall status: Deleting
   / Deleting service account 'contour-pkg-test-ns-sa'


   Uninstalled package 'contour-pkg' from namespace 'test-ns'
   ```

1. List the packages

   ```sh
   #List installed packages in the default namespace
   >>> tanzu package installed list
   NAME  DISPLAY-NAME  SHORT-DESCRIPTION

   #List installed packages across all namespaces
   >>> tanzu package installed list -A
   - Retrieving installed packages...
     NAME         PACKAGE-NAME              PACKAGE-VERSION        STATUS               NAMESPACE
     contour-pkg  contour.tanzu.vmware.com  1.15.1+vmware.1-tkg.1  Reconcile succeeded  test-ns
     mypkg        pkg.test.carvel.dev       2.0.0                  Reconcile succeeded  test-ns


   #List installed packages in user provided namespace
   >>> tanzu package installed list --namespace test-ns
   / Retrieving installed packages...
     NAME         PACKAGE-NAME              PACKAGE-VERSION        STATUS
     contour-pkg  contour.tanzu.vmware.com  1.15.1+vmware.1-tkg.1  Reconcile succeeded
     mypkg        pkg.test.carvel.dev       2.0.0                  Reconcile succeeded

   #List all available package CRs in default namespace
   >>> tanzu package available list
   / Retrieving available packages...
     NAME  DISPLAY-NAME  SHORT-DESCRIPTION

   #List all available package CRs across all namespace
   >>> tanzu package available list -A
   | Retrieving available packages...
     NAME                           DISPLAY-NAME          SHORT-DESCRIPTION                                                                    NAMESPACE
     harbor.tanzu.vmware.com        harbor                This package provides cloud native container registry service.                       test-ns
     pkg.test.carvel.dev            Test Package in repo  Package used for testing                                                             test-ns
     prometheus.tanzu.vmware.com    prometheus            This package provides an open-source systems monitoring and alerting toolkit         test-ns
     external-dns.tanzu.vmware.com  external-dns          This package provides DNS synchronization functionality.                             test-ns
     fluent-bit.tanzu.vmware.com    fluent-bit            This package provides log shipping functionality.                                    test-ns
     grafana.tanzu.vmware.com       grafana               This package allows you to visualize and analyze metrics data                        test-ns
     multus-cni.tanzu.vmware.com    multus-cni            This package provides ability for attaching multiple network interfaces to the pod.  test-ns
     cert-manager.tanzu.vmware.com  cert-manager          This package provides certificate management functionality.                          test-ns
     contour.tanzu.vmware.com       contour               This package provides ingress functionality.                                         test-ns

   #List all available packages for package name
   >>> tanzu package available list contour.tanzu.vmware.com -A
   / Retrieving package versions for contour.tanzu.vmware.com...
   NAME                      VERSION                RELEASED-AT  NAMESPACE
   contour.tanzu.vmware.com  1.15.1+vmware.1-tkg.1               test-ns
   ```

1. Delete the repository

   ```sh
   >>> tanzu package repository delete standard-repo --namespace test-ns
   Deleted package repository 'standard-repo' in namespace 'test-ns''
   ```

All the above commands are equipped with --kubeconfig flag to perform the package and repository operations on the desired cluster.

Example:

```sh
>>> tanzu package installed list -A --kubeconfig wc-kc-alpha8
 - Retrieving installed packages...
 NAME         PACKAGE-NAME              PACKAGE-VERSION        STATUS               NAMESPACE
 contour-pkg  contour.tanzu.vmware.com  1.15.1+vmware.1-tkg.1  Reconcile succeeded  test-ns
 mypkg        pkg.test.carvel.dev       2.0.0                  Reconcile succeeded  test-ns
```
