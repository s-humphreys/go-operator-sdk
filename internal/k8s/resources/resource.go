package resources

import "sigs.k8s.io/controller-runtime/pkg/client"

type Resource interface {
	Generate() *client.Object
}
