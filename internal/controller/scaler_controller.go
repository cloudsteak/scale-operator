/*
Copyright 2024 CloudSteak.

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
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/cloudsteak/scale-operator.git/api/v1alpha1"
)

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.scaler.cloudsteak.com,resources=scalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.scaler.cloudsteak.com,resources=scalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.scaler.cloudsteak.com,resources=scalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("Reconciling Scaler", "name", req.NamespacedName, "namespace", req.Namespace)

	scaler := &apiv1alpha1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		l.Error(err, "Failed to get Scaler")
		return ctrl.Result{}, nil
	}

	// Define the start and end time for scaling
	startTime := scaler.Spec.Start
	endTime := scaler.Spec.End

	// Get the current hour
	currentHour := time.Now().UTC().Hour()

	l.Info(fmt.Sprintf("############## Current hour: %d", currentHour))

	if currentHour >= startTime && currentHour <= endTime {
		l.Info("--- Scaling up deployments")
		if err = scaleDeployment(scaler, r, ctx, int32(scaler.Spec.Replicas)); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		l.Info("--- Scaling down deployments")
		if err = scaleDeployment(scaler, r, ctx, 1); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func scaleDeployment(scaler *apiv1alpha1.Scaler, r *ScalerReconciler, ctx context.Context, replicas int32) error {
	for _, deploy := range scaler.Spec.Deployments {
		dep := &v1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{
			Name:      deploy.Name,
			Namespace: deploy.Namespace,
		}, dep)
		if err != nil {
			log.Log.Error(err, "Failed to get Deployment")
			return nil
		}

		if dep.Spec.Replicas != &replicas {
			log.Log.Info("Scaling Deployment", "name", dep.Name, "namespace", dep.Namespace, "replicas_to", replicas, "replicas_from", dep.Spec.Replicas)
			dep.Spec.Replicas = &replicas
			err := r.Update(ctx, dep)
			if err != nil {
				log.Log.Error(err, "Failed to update Deployment")
				return nil
			}

			err = r.Status().Update(ctx, scaler)
			if err != nil {
				log.Log.Error(err, "Failed to update Scaler status")
				return nil
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Scaler{}).
		Complete(r)
}
