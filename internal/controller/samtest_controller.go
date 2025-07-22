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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
)

// SamtestReconciler reconciles a Samtest object
type SamtestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.k8s.capitalontap.com,resources=samtests/finalizers,verbs=update
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
	l := logf.FromContext(ctx)

	samtest := &cachev1alpha1.Samtest{}
	if err := r.Get(ctx, req.NamespacedName, samtest); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy := &k8s.Deployment{
		Name:      samtest.Name,
		Namespace: samtest.Namespace,
		Image:     samtest.Spec.Image,
	}

	deployObj := deploy.New()
	err := r.Get(ctx, client.ObjectKeyFromObject(deployObj), deployObj)

	// Deployment not found, create it
	if errors.IsNotFound(err) {
		// Set the owner references
		if err := ctrl.SetControllerReference(samtest, deployObj, r.Scheme); err != nil {
			l.Error(err, "failed to set controller reference")
			return ctrl.Result{}, err
		}

		// Create the deployment
		if err := r.Create(ctx, deployObj); err != nil {
			l.Error(err, "failed to create deployment")
			return ctrl.Result{}, err
		}
		l.Info("created deployment", "name", deployObj.Name, "namespace", deployObj.Namespace)

		// Set the condition on Samtest status
		c := deploy.NewCondition(k8s.ResourceCreated)
		meta.SetStatusCondition(&samtest.Status.Conditions, *c)

		// Update the status subresource
		if err := r.Status().Update(ctx, samtest); err != nil {
			l.Error(err, "failed to update samtest status")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
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
