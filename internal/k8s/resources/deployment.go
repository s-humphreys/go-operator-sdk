package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
)

type DeploymentEvents struct {
	Created        k8s.Event
	Updated        k8s.Event
	ImageOutOfSync k8s.Event
}

// NewDeploymentEvents returns a new DeploymentEvents struct with default values.
func NewDeploymentEvents(name string) DeploymentEvents {
	return DeploymentEvents{
		Created: k8s.Event{
			EventType: k8s.EventTypeNormal,
			Reason:    "DeploymentCreated",
			Message:   fmt.Sprintf("The Deployment '%s' has been created successfully", name),
		},
		Updated: k8s.Event{
			EventType: k8s.EventTypeNormal,
			Reason:    "DeploymentUpdated",
			Message:   fmt.Sprintf("The Deployment '%s' has been updated successfully", name),
		},
		ImageOutOfSync: k8s.Event{
			EventType: k8s.EventTypeWarning,
			Reason:    "DeploymentImageOutOfSync",
			Message:   fmt.Sprintf("The Deployment '%s' image is being updated to match the spec", name),
		},
	}
}

type Deployment struct {
	Name      string
	Namespace string
	Image     string
	Events    DeploymentEvents
}

// New creates a new Deployment with default events.
func NewDeployment(name, namespace, image string) *Deployment {
	return &Deployment{
		Name:      name,
		Namespace: namespace,
		Image:     image,
		Events:    NewDeploymentEvents(name),
	}
}

// Creates default labels for the Deployment.
func (d *Deployment) generateLabels() k8s.Labels {
	return k8s.Labels{
		"app": d.Name,
	}
}

// Creates a new Deployment Kubernetes object.
func (d *Deployment) Generate() *appsv1.Deployment {
	labels := d.generateLabels()

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "main",
							Image: d.Image,
						},
					},
				},
			},
		},
	}
}
