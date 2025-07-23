/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s/resources"
)

// SamtestReconciler reconciles a Samtest object
type SamtestReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Samtest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *SamtestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	samtest := &cachev1alpha1.Samtest{}
	if err := r.Get(ctx, req.NamespacedName, samtest); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy := resources.NewDeployment(samtest.Name, samtest.Namespace, samtest.Spec.Image)
	deployObj := deploy.Generate()
	err := r.Get(ctx, client.ObjectKeyFromObject(deployObj), deployObj)

	// Deployment not found, create it
	if errors.IsNotFound(err) {
		// Set the owner references
		if err := ctrl.SetControllerReference(samtest, deployObj, r.Scheme); err != nil {
			log.Error(err, "failed to set controller reference")
			return ctrl.Result{}, err
		}

		// Create the deployment
		if err := r.Create(ctx, deployObj); err != nil {
			log.Error(err, "failed to create deployment")
			return ctrl.Result{}, err
		}

		k8s.NewEvent(samtest, r.Recorder, deploy.Events.Created)
		log.Info("created deployment", "name", deployObj.Name, "namespace", deployObj.Namespace)

		if err := r.updateStatus(ctx, samtest, *k8s.NewStatusCondition(k8s.ResourcesCreated)); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Handle changes to existing Deployment
	image := ""
	if len(deployObj.Spec.Template.Spec.Containers) > 0 {
		image = deployObj.Spec.Template.Spec.Containers[0].Image
	}

	// Check if the image is out of sync & update if necessary
	if image != samtest.Spec.Image {
		log.Info("deployment image is out of sync, updating", "current", image, "desired", samtest.Spec.Image)
		k8s.NewEvent(samtest, r.Recorder, deploy.Events.ImageOutOfSync)

		if err := r.updateStatus(ctx, samtest, *k8s.NewStatusCondition(k8s.ResourcesOutOfSync)); err != nil {
			return ctrl.Result{}, err
		}

		deployObj.Spec.Template.Spec.Containers[0].Image = samtest.Spec.Image
		if err := r.Update(ctx, deployObj); err != nil {
			log.Error(err, "failed to update deployment image")
			return ctrl.Result{}, err
		}

		k8s.NewEvent(samtest, r.Recorder, deploy.Events.Updated)
		if err := r.updateStatus(ctx, samtest, *k8s.NewStatusCondition(k8s.ResourcesReconciled)); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SamtestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Samtest{}).
		Owns(&appsv1.Deployment{}).
		Named("samtest").
		Complete(r)
}

// Updates the status of the Samtest resource with a provided condition
func (r *SamtestReconciler) updateStatus(ctx context.Context, samtest *cachev1alpha1.Samtest, condition metav1.Condition) error {
	log := logf.FromContext(ctx)
	meta.SetStatusCondition(&samtest.Status.Conditions, condition)
	if err := r.Status().Update(ctx, samtest); err != nil {
		log.Error(err, "failed to update samtest status", "conditionType", condition.Type, "conditionStatus", condition.Status)
		return err
	}
	return nil
}
