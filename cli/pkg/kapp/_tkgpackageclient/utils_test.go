// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tkgpackageclient

import (
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kapppkg "github.com/vmware-tanzu/carvel-kapp-controller/pkg/apiserver/apis/datapackaging/v1alpha1"
)

var rawInvalidDummyPackageSchema = `
foo:
  bar: xxx
`
var rawValidDummyPackageSchema = `
title: external-dns.community.tanzu.vmware.com.0.8.0 values schema
examples:
  - namespace: tanzu-system-service-discovery
    deployment:
      args:
        - --source=service
        - --txt-owner-id=k8s
        - --domain-filter=k8s.example.org
        - --namespace=tanzu-system-service-discovery
        - --provider=rfc2136
        - --rfc2136-host=100.69.97.77
        - --rfc2136-port=53
        - --rfc2136-zone=k8s.example.org
        - --rfc2136-tsig-secret=MTlQs3NNU=
        - --rfc2136-tsig-secret-alg=hmac-sha256
        - --rfc2136-tsig-keyname=externaldns-key
        - --rfc2136-tsig-axfr
      env: []
      securityContext: {}
      volumeMounts: []
      volumes: []
properties:
  namespace:
    type: string
    description: The namespace in which to deploy ExternalDNS.
    default: external-dns
    examples:
      - external-dns
  deployment:
    type: object
    description: Deployment related configuration
    properties:
      args:
        type: array
        description: |
          List of arguments passed via command-line to external-dns.  For
          more guidance on configuration options for your desired DNS
          provider, consult the ExternalDNS docs at
          https://github.com/kubernetes-sigs/external-dns#running-externaldns
        items:
          type: string
      env:
        type: array
        description: "List of environment variables to set in the external-dns container."
        items:
          $ref: "#/definitions/io.k8s.api.core.v1.EnvVar"
      securityContext:
        description: "SecurityContext defines the security options the external-dns container should be run with. More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/"
        $ref: "#/definitions/io.k8s.api.core.v1.SecurityContext"
      volumeMounts:
        type: array
        description: "Pod volumes to mount into the external-dns container's filesystem."
        items:
          $ref: "#/definitions/io.k8s.api.core.v1.VolumeMount"
      volumes:
        type: array
        description: "List of volumes that can be mounted by containers belonging to the external-dns pod. More info: https://kubernetes.io/docs/concepts/storage/volumes"
        items:
          $ref: "#/definitions/io.k8s.api.core.v1.Volume"
`

var _ = Describe("Test utils", func() {
	Context("parsing package values schema", func() {
		var (
			valuesSchema kapppkg.ValuesSchema
			parser       *PackageValuesSchemaParser
			err          error
		)
		It("should not have error if values schema is valid", func() {
			valuesSchema = kapppkg.ValuesSchema{
				OpenAPIv3: runtime.RawExtension{
					Raw: []byte(rawValidDummyPackageSchema),
				},
			}
			parser, err = NewValuesSchemaParser(valuesSchema)
			Expect(err).NotTo(HaveOccurred())
			parsedProperties, err := parser.ParseProperties()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(parsedProperties)).ShouldNot(Equal(0))
			Expect(len(parser.DataValueProperties)).ShouldNot(Equal(0))
		})
		It("should have error if values schema is invalid", func() {
			valuesSchema = kapppkg.ValuesSchema{
				OpenAPIv3: runtime.RawExtension{
					Raw: []byte(rawInvalidDummyPackageSchema),
				},
			}
			parser, err = NewValuesSchemaParser(valuesSchema)
			Expect(err).NotTo(HaveOccurred())
			parsedProperties, err := parser.ParseProperties()
			Expect(err).To(HaveOccurred())
			Expect(len(parsedProperties)).Should(Equal(0))
			Expect(len(parser.DataValueProperties)).Should(Equal(0))
		})
	})
})
