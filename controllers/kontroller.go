/*
Copyright 2023.

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

package controllers

import (
	"context"
	"errors"

	"github.com/frantjc/go-fn"
	apiv1alpha1 "github.com/frantjc/kontrol/api/v1alpha1"
	"github.com/frantjc/kontrol/cr"
	"github.com/frantjc/kontrol/pkg"
	"github.com/frantjc/kontrol/pkg/lbl"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/cri-api/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KontrollerReconciler reconciles a Kontroller object.
type KontrollerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger logr.Logger
}

//+kubebuilder:rbac:groups="",resources=events;pods;namespaces;serviceaccounts;services;services/finalizers,verbs=*
//+kubebuilder:rbac:groups=apps,resources=deployments;replicasets,verbs=*
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=*
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=*

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *KontrollerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		_          = r.Logger
		packager   = new(lbl.Packager)
		kontroller = &apiv1alpha1.Kontroller{}
	)

	if err := r.Client.Get(ctx, req.NamespacedName, kontroller); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	image, err := cr.Pull(ctx, kontroller.Spec.Image)
	if err != nil {
		return ctrl.Result{}, err
	}

	packaging, err := packager.Unpackage(ctx, image)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, errors.Join(
		fn.Map(
			pkg.PackagingToObjects(kontroller.Spec.Image, req.Name, req.Namespace, packaging),
			func(obj client.Object, _ int) error {
				obj.SetOwnerReferences([]metav1.OwnerReference{
					{
						APIVersion:         apiv1alpha1.GroupVersion.String(),
						Kind:               "Kontroller",
						Name:               req.Name,
						BlockOwnerDeletion: fn.Ptr(true),
						Controller:         fn.Ptr(true),
						UID:                kontroller.GetUID(),
					},
				})

				if err := r.Client.Create(ctx, obj); err != nil {
					return err
				}

				return r.Client.Update(ctx, obj)
			},
		)...,
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *KontrollerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Kontroller{}).
		Complete(r)
}
