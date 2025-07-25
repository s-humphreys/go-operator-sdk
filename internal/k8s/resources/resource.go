package resources

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
)

type Resource interface {
	New(crd *cachev1alpha1.Samtest) Resource
	Kind() string
	Generate() client.Object
	IsEqual(client.Object) bool
}
