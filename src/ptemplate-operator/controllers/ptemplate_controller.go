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
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kmmv1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cachev1alpha1 "github.com/chr15p/partner_templates/api/v1alpha1"

	//import the monitoring pacakge
	"github.com/chr15p/partner_templates/monitoring"
)

// PtemplateReconciler reconciles a Ptemplate object
type PtemplateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pt.example.com,resources=ptemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pt.example.com,resources=ptemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pt.example.com,resources=ptemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PtemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	logger.Info("Reconcile Module")

	// make sure we actually have a ptemplate resource defined, if not we can just end here,
	// or if we get an error talking to k8s we can fail here and return an error.
	existingPtemplate := &cachev1alpha1.Ptemplate{}
	err := r.Get(ctx, req.NamespacedName, existingPtemplate)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// we have two major resources to control, the Module that loads the kmod
	// and the consumer pods that use it, so we need to reconcile each in turn

	// first the Module
	modResult, err := r.ReconcileModule(ctx, req, existingPtemplate)
	if err != nil {
		return ctrl.Result{}, err
	}
	if modResult != nil {
		return *modResult, nil
	}

	// Then the Consumer
	conResult, err := r.ReconcileConsumer(ctx, req, existingPtemplate)
	if err != nil {
		return ctrl.Result{}, err
	}
	if conResult != nil {
		return *conResult, nil
	}

	// Check if we need to update the Status
	podList := &v1.PodList{}

	listOpts := []client.ListOption{
		client.InNamespace(existingPtemplate.Namespace),
		client.MatchingLabels(labelsForApp(existingPtemplate.Name)),
	}

	if err = r.List(ctx, podList, listOpts...); err != nil {
		return ctrl.Result{}, err
	}

	podNames := getPodNames(podList.Items)

	if !reflect.DeepEqual(podNames, existingPtemplate.Status.Consumers) {
		patch := client.MergeFrom(existingPtemplate.DeepCopy())

		existingPtemplate.Status.Consumers = podNames
		logger.Info("Reconcile Status", "setting Status.Consumers to", podNames)
		err = r.Status().Patch(ctx, existingPtemplate, patch)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// report the number of consumer pods to prometheus
	monitoring.PtemplateGauge.Set(float64(len(podNames)))
	//report the number of times reconcile has run
	monitoring.PtemplateCounter.Inc()

	logger.Info("Reconcile completed")
	return ctrl.Result{}, nil
}

func labelsForApp(name string) map[string]string {
	return map[string]string{"app": name}
}

func getPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

/*
* ReconcileModule - Ensure the Module resource exists and matches our spec
 */
func (r *PtemplateReconciler) ReconcileModule(ctx context.Context, req ctrl.Request, ptResource *cachev1alpha1.Ptemplate) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if ptResource.Spec.MaxDev == 0 {
		ptResource.Spec.MaxDev = 5
	}
	if ptResource.Spec.DevicePlugin == "" {
		ptResource.Spec.DevicePlugin = "quay.io/chrisp262/pt-device-plugin:latest"
	}

	existingMod, err := r.getExistingModule(ctx, req.NamespacedName)
	// if there is no kmm.module resource already exisiting then we need to create it
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Define and create a new module.
			module := r.getNewModule(ctx, ptResource)
			if err = r.Create(ctx, module); err != nil {
				return &ctrl.Result{}, err
			}
			return &ctrl.Result{Requeue: true}, nil
		} else {
			return &ctrl.Result{}, err
		}
	}

	// if there is a kmm.module resource with our expected name we need to check if our CR has changed and update the
	// module with any fields that are differ

	// TODO: this is a clunky way to check if the modules params have changed, needs fixing.
	maxdev := fmt.Sprintf("max_dev=%d", ptResource.Spec.MaxDev)
	defaultmsg := fmt.Sprintf("default_msg=%s", ptResource.Spec.DefaultMsg)
	parameters := []string{maxdev, defaultmsg}
	existingParams := existingMod.Spec.ModuleLoader.Container.Modprobe.Parameters

	if !reflect.DeepEqual(parameters, existingParams) ||
		ptResource.Spec.DevicePlugin != existingMod.Spec.DevicePlugin.Container.Image {

		logger.Info("Module not in sync, updating")
		existingMod.Spec.DevicePlugin.Container.Image = ptResource.Spec.DevicePlugin
		existingMod.Spec.ModuleLoader.Container.Modprobe.Parameters = parameters
		if err = r.Update(ctx, existingMod); err != nil {
			logger.Info("Failed to Update Module", "error:", err)

			return &ctrl.Result{}, err
		}
		return &ctrl.Result{Requeue: true}, nil
	}

	// now our kmod should match our spec, so we check the consumer

	logger.Info("Module Reconciled", "mod=", existingMod)

	return nil, nil
}

/*
* getExistingModule
* Call into Kubernetes and get any Module resources that match our name and namespace (i.e our namespacedName)
 */
func (r *PtemplateReconciler) getExistingModule(ctx context.Context, namespacedName types.NamespacedName) (*kmmv1beta1.Module, error) {
	logger := log.FromContext(ctx)

	mod := kmmv1beta1.Module{}

	if err := r.Client.Get(ctx, namespacedName, &mod); err != nil {
		logger.Info("Failed to getExistingModule", "error", err)
		return nil, err
	}
	return &mod, nil
}

/*
* getNewModule
* Returns a new Module object that the r.Create(ctx, module) call in ReconcileModule() will send to k8s for it to create
* once k8s has created the Module reources with our settings then the KMM operator will see it and reconcile it
* so we are loading our kernel module at one step removed becuase leveraging it saves us a lot of code and pain
 */

func (r *PtemplateReconciler) getNewModule(ctx context.Context, ptResource *cachev1alpha1.Ptemplate) *kmmv1beta1.Module {
	logger := log.FromContext(ctx)

	maxdev := fmt.Sprintf("max_dev=%d", ptResource.Spec.MaxDev)
	defaultmsg := fmt.Sprintf("default_msg=%s", ptResource.Spec.DefaultMsg)

	logger.Info("getNewModule", "max_dev=", maxdev, "default_msg=", defaultmsg)

	module := &kmmv1beta1.Module{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ptResource.Name,
			Namespace: ptResource.Namespace,
		},
		Spec: kmmv1beta1.ModuleSpec{
			DevicePlugin: &kmmv1beta1.DevicePluginSpec{
				Container: kmmv1beta1.DevicePluginContainerSpec{
					Image: "quay.io/chrisp262/pt-device-plugin:plugin-latest",
				},
			},
			ModuleLoader: kmmv1beta1.ModuleLoaderSpec{
				Container: kmmv1beta1.ModuleLoaderContainerSpec{
					Modprobe: kmmv1beta1.ModprobeSpec{
						ModuleName: "ptemplate_char_dev",
						Parameters: []string{maxdev, defaultmsg},
					},
					ImagePullPolicy: "Always",
					KernelMappings: []kmmv1beta1.KernelMapping{
						{
							Regexp:         "^.+$",
							ContainerImage: "quay.io/chrisp262/pt-char-dev:${KERNEL_FULL_VERSION}-sb",
						},
					},
				},
			},
			ImageRepoSecret: ptResource.Spec.ImagePullSecret,
			Selector:        ptResource.Spec.Selector,
		},
	}

	logger.Info("getNewModule", "module=", module)

	// Set Ptemplate instance as the owner in the metadata.ownerReferences stanza
	// NOTE: calling SetControllerReference, and setting owner references in
	// general, is important as it allows deleted objects to be garbage collected.
	// TODO: this function can return an error type but theres not much of valuei
	// we can do with it so we're ignoing it for the moment
	controllerutil.SetControllerReference(ptResource, module, r.Scheme)
	return module
}

// ********** Consumer **********
/*
* ReconcileConsumer
* the consumer pods will be managed by a daemonset so we get one pod per node that matches our selector
* this method checks if the daemonset exists and if not creates it
* and if it does checks it against our ptemplate.spec and updates the daemonset if anything has changed
 */
func (r *PtemplateReconciler) ReconcileConsumer(ctx context.Context, req ctrl.Request, ptResource *cachev1alpha1.Ptemplate) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// build the consumer name from the CR
	consumerDSName := types.NamespacedName{
		Namespace: req.NamespacedName.Namespace,
		Name:      req.NamespacedName.Name + "-consumer",
	}

	// set some default values that are option in the Spec but we actually need
	if ptResource.Spec.RequiredDev == 0 {
		ptResource.Spec.RequiredDev = 1
	}
	if ptResource.Spec.ConsumerImage == "" {
		ptResource.Spec.ConsumerImage = "quay.io/chrisp262/pt-device-plugin:consumer-latest"
	}

	existingDS, err := r.getExistingConsumer(ctx, consumerDSName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Define and create a new module.
			module := r.getNewConsumer(ctx, ptResource, consumerDSName)
			if err = r.Create(ctx, module); err != nil {
				return &ctrl.Result{}, err
			}
			return &ctrl.Result{Requeue: true}, nil
		} else {
			return &ctrl.Result{}, err
		}
	}

	consumerSpec := existingDS.Spec.Template.Spec
	consumerContainer := existingDS.Spec.Template.Spec.Containers[0]

	consumerLimit := consumerContainer.Resources.Limits["example.com/ptemplate"]

	if consumerContainer.Image != ptResource.Spec.ConsumerImage ||
		!reflect.DeepEqual(consumerSpec.NodeSelector, ptResource.Spec.Selector) ||
		consumerLimit.Value() != ptResource.Spec.RequiredDev {

		logger.Info("Consumer not in sync, updating")
		consumerContainer.Image = ptResource.Spec.ConsumerImage
		consumerSpec.NodeSelector = ptResource.Spec.Selector
		consumerContainer.Resources.Limits["example.com/ptemplate"] = *resource.NewQuantity(ptResource.Spec.RequiredDev, resource.DecimalSI)

		if err = r.Update(ctx, existingDS); err != nil {
			logger.Info("Failed to Update Consumer", "error:", err)

			return &ctrl.Result{}, err
		}
		return &ctrl.Result{Requeue: true}, nil
	}

	logger.Info("Consumer Reconciled", "consumer=", existingDS)

	return nil, nil
}

/*
* getExistingConsumer
* Call into Kubernetes and get any Daemonset resources that match our name and namespace (i.e our namespacedName)
 */
func (r *PtemplateReconciler) getExistingConsumer(ctx context.Context, namespacedName types.NamespacedName) (*appsv1.DaemonSet, error) {
	logger := log.FromContext(ctx)

	ds := appsv1.DaemonSet{}

	if err := r.Client.Get(ctx, namespacedName, &ds); err != nil {
		logger.Info("Failed to get existing Consumer Daemonset", "error", err)
		return nil, err
	}
	return &ds, nil
}

/*
* getNewConsumer
* Returns a new Daemonset object that the r.Create(ctx, module) call in ReconcileConsumer() will send to k8s for it to create
*
 */
func (r *PtemplateReconciler) getNewConsumer(ctx context.Context, ptResource *cachev1alpha1.Ptemplate, consumerDSName types.NamespacedName) *appsv1.DaemonSet {
	logger := log.FromContext(ctx)

	maxdev := fmt.Sprintf("max_dev=%d", ptResource.Spec.MaxDev)
	defaultmsg := fmt.Sprintf("default_msg=%s", ptResource.Spec.DefaultMsg)

	logger.Info("getNewConsumer", "max_dev=", maxdev, "default_msg=", defaultmsg)

	// we want /dev to be mounted into our container so we need Volume and VolumeMount objects to enable that
	containerVolumeMounts := []v1.VolumeMount{
		{
			Name:      "host-dev",
			MountPath: "/hostdev",
		},
	}

	containerVolume := []v1.Volume{
		{
			Name: "host-dev",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/dev",
				},
			},
		},
	}

	// generate the spec fo the daemonset
	consumer := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      consumerDSName.Name,
			Namespace: consumerDSName.Namespace,
			Labels:    labelsForApp(ptResource.Name),
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsForApp(ptResource.Name),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelsForApp(ptResource.Name),
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "consumer",
							Image:           ptResource.Spec.ConsumerImage,
							SecurityContext: &v1.SecurityContext{Privileged: pointer.Bool(true)},
							ImagePullPolicy: "Always",
							Resources: v1.ResourceRequirements{
								Limits: map[v1.ResourceName]resource.Quantity{
									"example.com/ptemplate": *resource.NewQuantity(ptResource.Spec.RequiredDev, resource.DecimalSI),
								},
							},
							VolumeMounts: containerVolumeMounts,
						},
					},
					ImagePullSecrets: []v1.LocalObjectReference{
						*ptResource.Spec.ImagePullSecret,
					},
					NodeSelector: ptResource.Spec.Selector,
					Volumes:      containerVolume,
				},
			},
		},
	}

	// Set Ptemplate instance (i.e. ptResource) as the metadata.ownerReferences dor the daemonset (i.e. consumer)
	// NOTE: calling SetControllerReference, and setting owner references in
	// general, is important as it allows deleted objects to be garbage collected.
	// TODO: this function can return an error type but theres not much of valuei
	// we can do with it so we're ignoing it for the moment
	controllerutil.SetControllerReference(ptResource, consumer, r.Scheme)

	return consumer

}

// SetupWithManager sets up the controller with the Manager.
// The 'Owns()' call tells k8s to re-run the reconciler whenever any of that type of resources change
// in our case  whenever a Daemonset is updated we want to be reconciled in case its our daemonset
func (r *PtemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Ptemplate{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}
