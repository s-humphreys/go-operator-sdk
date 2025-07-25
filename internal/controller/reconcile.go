package controller

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/s-humphreys/go-operator-sdk/api/v1alpha1"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s"
	"github.com/s-humphreys/go-operator-sdk/internal/k8s/resources"
)

func (r *SamtestReconciler) reconcileResource(
	log logr.Logger,
	ctx context.Context,
	crd *cachev1alpha1.Samtest,
	resource resources.Resource,
) (ctrl.Result, error) {
	desiredObj := resource.Generate()
	kind := resource.Kind()
	name := desiredObj.GetName()
	log.Info("reconciling resource", "kind", kind, "name", name)

	foundObj := desiredObj.DeepCopyObject().(client.Object)
	err := r.Get(ctx, client.ObjectKeyFromObject(desiredObj), foundObj)

	if errors.IsNotFound(err) {
		log.Info("resource not found, creating", "kind", kind)

		// Set owner references
		if err := ctrl.SetControllerReference(crd, foundObj, r.Scheme); err != nil {
			log.Error(err, "failed to set controller reference")
			return ctrl.Result{}, err
		}

		// Create the resource
		if err := r.Create(ctx, foundObj); err != nil {
			log.Error(err, "failed to create resource", "kind", kind)
			k8s.NewCreateErrorEvent(crd, r.Recorder, kind, name)
			return ctrl.Result{}, err
		}

		k8s.NewCreatedEvent(crd, r.Recorder, kind, name)
		log.Info("resource created", "kind", kind)
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Perform equality check & update if different
	// if !resource.IsEqual(foundObj) {
	// 	log.Info("resource is out of sync, updating", "kind", kind, "name", name)
	// 	k8s.NewOutOfSyncEvent(crd, r.Recorder, kind, name)

	// 	switch desired := desiredObj.(type) {
	// 	case *appsv1.Deployment:
	// 		found := foundObj.(*appsv1.Deployment)
	// 		found.Spec = desired.Spec
	// 	case *corev1.Service:
	// 		found := foundObj.(*corev1.Service)
	// 		// For services, we must preserve the ClusterIP
	// 		clusterIP := found.Spec.ClusterIP
	// 		found.Spec = desired.Spec
	// 		found.Spec.ClusterIP = clusterIP
	// 	}

	// 	if err := r.Update(ctx, foundObj); err != nil {
	// 		log.Error(err, "failed to update resource")
	// 		k8s.NewUpdateErrorEvent(crd, r.Recorder, kind, name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	k8s.NewUpdatedEvent(crd, r.Recorder, kind, name)
	// }

	// If resource is not managed by this controller, set the owner reference
	if !metav1.IsControlledBy(foundObj, crd) {
		log.Info("resource not managed, setting owner reference", "kind", kind)
		if err := ctrl.SetControllerReference(crd, foundObj, r.Scheme); err != nil {
			log.Error(err, "failed to set controller reference on existing Deployment")
			return ctrl.Result{}, err
		}
		if err := r.Update(ctx, foundObj); err != nil {
			log.Error(err, "failed to update resource with owner reference", "kind", kind)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
