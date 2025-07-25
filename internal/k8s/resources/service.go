package resources

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
)

type Service struct {
	Name      string
	Namespace string
	Labels    k8s.Labels
}

// New creates a new Service with default events.
func (s *Service) New(crd *cachev1alpha1.Samtest) Resource {
	return &Service{
		Name:      crd.Name,
		Namespace: crd.Namespace,
		Labels:    k8s.CreateLabels(crd.Name),
	}
}

// Returns the resource kind.
func (s *Service) Kind() string {
	return "Service"
}

// Creates a new Service Kubernetes object.
func (s *Service) Generate() client.Object {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       s.Kind(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
			Labels:    s.Labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: s.Labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromString("http"),
				},
			},
		},
	}
}

func (d *Service) IsEqual(found client.Object) bool {
	foundService, ok := found.(*corev1.Service)
	if !ok {
		return false
	}

	return equality.Semantic.DeepEqual(d.Generate().(*corev1.Service).Spec, foundService.Spec)
}
