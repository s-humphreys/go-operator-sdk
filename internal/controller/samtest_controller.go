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
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

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

	if samtest.Spec.Suspend {
		log.Info("resource is suspended, skipping reconciliation")
		return ctrl.Result{}, nil
	}

	managedResources := []resources.Resource{
		&resources.Deployment{},
		&resources.Service{},
	}

	// TODO(Sam) - determine if there is a diff (do in prior loop, early return & update status)

	if err := r.updateStatus(ctx, samtest, k8s.NewStatusCondition(k8s.ProgressingResources)); err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile each resource
	var wg sync.WaitGroup
	errs := make(chan error, len(managedResources))

	for _, resource := range managedResources {
		res := resource.New(samtest)
		wg.Add(1)

		go func(reconcileResource resources.Resource) {
			defer wg.Done()
			if _, err := r.reconcileResource(log, ctx, samtest, reconcileResource); err != nil {
				errs <- err
			}
		}(res)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			_ = r.updateStatus(ctx, samtest, k8s.NewStatusCondition(k8s.ResourcesFailed))
			return ctrl.Result{}, err
		}
	}

	if err := r.updateStatus(ctx, samtest, k8s.NewStatusCondition(k8s.ResourcesReady)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SamtestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Samtest{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("samtest").
		Complete(r)
}

// Updates the status of the Samtest resource with a provided condition.
func (r *SamtestReconciler) updateStatus(ctx context.Context, samtest *cachev1alpha1.Samtest, condition metav1.Condition) error {
	log := logf.FromContext(ctx)
	meta.SetStatusCondition(&samtest.Status.Conditions, condition)
	if err := r.Status().Update(ctx, samtest); err != nil {
		log.Error(err, "failed to update status", "conditionType", condition.Type, "conditionStatus", condition.Status)
		return err
	}
	return nil
}
