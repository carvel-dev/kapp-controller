package convert

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/davecgh/go-spew/spew"
	rukpakv1alpha1 "github.com/operator-framework/rukpak/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/rand"
)

func GenerateTemplateHash(template *rukpakv1alpha1.BundleTemplate) string {
	hasher := fnv.New32a()
	DeepHashObject(hasher, template)
	return rand.SafeEncodeString(fmt.Sprintf("%x", hasher.Sum32())[:6])
}

// DeepHashObject writes specified object to hash using the spew library
// which follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
// From https://github.com/operator-framework/operator-lifecycle-manager/blob/master/pkg/lib/kubernetes/pkg/util/hash/hash.go
func DeepHashObject(hasher hash.Hash, objectToWrite interface{}) {
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	printer.Fprintf(hasher, "%#v", objectToWrite)
}
