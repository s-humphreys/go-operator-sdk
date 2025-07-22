package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name      string
	Namespace string
	Image     string
}

// Creates default labels for the Deployment
func (d *Deployment) newLabels() Labels {
	return Labels{
		"app": d.Name,
	}
}

// Creates a new Deployment object
func (d *Deployment) New() *appsv1.Deployment {
	labels := d.newLabels()

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

// Creates a condition for the Deployment resource
func (d *Deployment) NewCondition(state ResourceState) *metav1.Condition {
	if state == ResourceCreated {
		return &metav1.Condition{
			Type:    "Available",
			Status:  metav1.ConditionTrue,
			Reason:  "DeploymentCreated",
			Message: fmt.Sprintf("Deployment has been created successfully - %s/%s", d.Namespace, d.Name),
		}
	}
	return nil
}
