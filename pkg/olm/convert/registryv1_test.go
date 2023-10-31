package convert

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRegistryV1Converter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RegstryV1 suite")
}

var _ = Describe("RegistryV1 Suite", func() {
	var _ = Describe("Convert", func() {
		var (
			registryv1Bundle RegistryV1
			installNamespace string
			targetNamespaces []string
		)
		Context("Should set the namespaces of object correctly", func() {
			var (
				svc corev1.Service
				csv v1alpha1.ClusterServiceVersion
			)
			BeforeEach(func() {
				csv = v1alpha1.ClusterServiceVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name: "testCSV",
					},
					Spec: v1alpha1.ClusterServiceVersionSpec{
						InstallModes: []v1alpha1.InstallMode{{Type: v1alpha1.InstallModeTypeAllNamespaces, Supported: true}},
					},
				}
				svc = corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name: "testService",
					},
				}
				svc.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
				installNamespace = "testInstallNamespace"
			})

			It("should set the namespace to installnamespace if not available", func() {
				By("creating a registry v1 bundle")
				unstructuredSvc := convertToUnstructured(svc)
				registryv1Bundle = RegistryV1{
					PackageName: "testPkg",
					CSV:         csv,
					Others:      []unstructured.Unstructured{unstructuredSvc},
				}

				By("converting to plain")
				plainBundle, err := Convert(registryv1Bundle, installNamespace, targetNamespaces)
				Expect(err).NotTo(HaveOccurred())

				By("verifying if plain bundle has required objects")
				Expect(plainBundle).NotTo(BeNil())
				Expect(len(plainBundle.Objects)).To(BeEquivalentTo(2))

				By("verifying if ns has been set correctly")
				resObj := containsObject(unstructuredSvc, plainBundle.Objects)
				Expect(resObj).NotTo(BeNil())
				Expect(resObj.GetNamespace()).To(BeEquivalentTo(installNamespace))
			})

			It("should override namespace if already available", func() {
				By("creating a registry v1 bundle")
				svc.SetNamespace("otherNs")
				unstructuredSvc := convertToUnstructured(svc)
				unstructuredSvc.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})

				registryv1Bundle = RegistryV1{
					PackageName: "testPkg",
					CSV:         csv,
					Others:      []unstructured.Unstructured{unstructuredSvc},
				}

				By("converting to plain")
				plainBundle, err := Convert(registryv1Bundle, installNamespace, targetNamespaces)
				Expect(err).NotTo(HaveOccurred())

				By("verifying if plain bundle has required objects")
				Expect(plainBundle).NotTo(BeNil())
				Expect(len(plainBundle.Objects)).To(BeEquivalentTo(2))

				By("verifying if ns has been set correctly")
				resObj := containsObject(unstructuredSvc, plainBundle.Objects)
				Expect(resObj).NotTo(BeNil())
				Expect(resObj.GetNamespace()).To(BeEquivalentTo(installNamespace))
			})

			Context("Should error when object is not supported", func() {
				It("should error when unsupported GVK is passed", func() {
					By("creating an unsupported kind")
					event := corev1.Event{
						ObjectMeta: metav1.ObjectMeta{
							Name: "testEvent",
						},
					}

					unstructuredEvt := convertToUnstructured(event)
					unstructuredEvt.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Event"})

					registryv1Bundle = RegistryV1{
						PackageName: "testPkg",
						CSV:         csv,
						Others:      []unstructured.Unstructured{unstructuredEvt},
					}

					By("converting to plain")
					plainBundle, err := Convert(registryv1Bundle, installNamespace, targetNamespaces)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("bundle contains unsupported resource"))
					Expect(plainBundle).To(BeNil())
				})
			})

			Context("Should not set ns cluster scoped object is passed", func() {
				It("should not error when cluster scoped obj is passed and not set its namespace", func() {
					By("creating an unsupported kind")
					pc := schedulingv1.PriorityClass{
						ObjectMeta: metav1.ObjectMeta{
							Name: "testPriorityClass",
						},
					}

					unstructuredpriorityclass := convertToUnstructured(pc)
					unstructuredpriorityclass.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PriorityClass"})

					registryv1Bundle = RegistryV1{
						PackageName: "testPkg",
						CSV:         csv,
						Others:      []unstructured.Unstructured{unstructuredpriorityclass},
					}

					By("converting to plain")
					plainBundle, err := Convert(registryv1Bundle, installNamespace, targetNamespaces)
					Expect(err).NotTo(HaveOccurred())

					By("verifying if plain bundle has required objects")
					Expect(plainBundle).NotTo(BeNil())
					Expect(len(plainBundle.Objects)).To(BeEquivalentTo(2))

					By("verifying if ns has been set correctly")
					resObj := containsObject(unstructuredpriorityclass, plainBundle.Objects)
					Expect(resObj).NotTo(BeNil())
					Expect(resObj.GetNamespace()).To(BeEmpty())
				})
			})
		})

	})
})

func convertToUnstructured(obj interface{}) unstructured.Unstructured {
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&obj)
	Expect(err).NotTo(HaveOccurred())
	Expect(unstructuredObj).NotTo(BeNil())
	return unstructured.Unstructured{Object: unstructuredObj}
}

func containsObject(obj unstructured.Unstructured, result []client.Object) client.Object {
	for _, o := range result {
		// Since this is a controlled env, comparing only the names is sufficient for now.
		// In future, compare GVKs too by ensuring its set on the unstructuredObj.
		if o.GetName() == obj.GetName() {
			return o
		}
	}
	return nil
}
