package resources

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
)

type Deployment struct {
	Name      string
	Namespace string
	Image     string
	Labels    k8s.Labels
}

// New creates a new Deployment with default values
func (d *Deployment) New(crd *cachev1alpha1.Samtest) Resource {
	return &Deployment{
		Name:      crd.Name,
		Namespace: crd.Namespace,
		Image:     crd.Spec.Image,
		Labels:    k8s.CreateLabels(crd.Name),
	}
}

// Returns the resource kind.
func (d *Deployment) Kind() string {
	return "Deployment"
}

// Creates a new Deployment Kubernetes object.
func (d *Deployment) Generate() client.Object {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       d.Kind(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
			Labels:    d.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: d.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: d.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: d.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *Deployment) IsEqual(found client.Object) bool {
	foundDeployment, ok := found.(*appsv1.Deployment)
	if !ok {
		return false
	}

	return equality.Semantic.DeepEqual(d.Generate().(*appsv1.Deployment).Spec, foundDeployment.Spec)
}
