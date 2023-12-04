# Using KMM with your operator

## Introduction

Often loading a kernel module by itself is only part of the problem, requiring other resources like user land processes, or carefully calculated kmod parameters, to provide their full benefit.

The recommended way to tie together the loading of kernel modules with other OpenShift resources is to create a device specific operator with its own Custom Resources that then creates and manages the other Kubernetes objects it needs. The users of your new operator can then create new instances of your Custom Resources, and your operator will then process this to automatically create Module resources for the KMM operator, DaemonSets to create pods, and any other resources needed.

The [Ptemplate Operator](../src/ptemplate-operator) is an example of a way to implement this. It is mostly designed to be a simple teaching tool rather than a useful operator but it does load a kmod and create a set of pods to utilise it.

## Overview

Without the Operator a Custom Resource does nothing, its just some data fields stored as a serialised object. The operator associated with the CR creates and updates other Kubernetes resources based on this data by sending API calls to Kubernetes. As such the Operators `Reconcile()` function needs to:

1. get the data fields in the CR its been called for

2. get the data fields for the child resources that need to exist via `Get` requests to the Kubernetes API server

3. determine if there are any differences between the Spec fields of the CR and the child resources

4. make appropriate `Create` and `Patch` calls to fix any discrepancies in the sub-resources

5. if necessary update the `Status` fields of the CR via one or more calls to `Update`

## controllers/ptemplate_controller.go Functions

### Reconcile()

```go
func (r *PtemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
```

This gets called by Kubernetes with a `req` object which includes the `NamespacedName` of the resource it was called for which we can then use to get all the details of that resource

```go
existingPtemplate := &cachev1alpha1.Ptemplate{}
err := r.Get(ctx, req.NamespacedName, existingPtemplate)
```

The `existingPtemplate` object at this point is a golang equivalent of the yaml that `oc describe` would show.

Once we have those details we can examine them, determine any differences between the `existingPtemplate.spec` and the child resources we want to exist (including if there are no child resources) and then make the appropriate changes to create or update them.

The rest of this function then calls separate functions to reconcile the two child resources (the kmm `Module` resource, and a DaemonSet for the consumer pods), and updates the `Ptemplate` resource's status

### ReconcileModule()

The managing of the KMM `Module` resource is broken out into its own set of functions just for code-clarity reasons.  The function `ReconcileModule()` is called by the main `Reconcile()` function (see above).

`ReconcileModule` starts of by looking at the `*cachev1alpha1.Ptemplate` object that `Reconcile()` has passed it and sets up some default values if they are not already defined. It then tries to get the details of the `Module` resource it is managing via a Get request to the Kubernetes API (handled in the `getExistingModule()` function).

If the `Module` resource is not found, then the Operator can assume its not been created yet, or has been accidentally deleted, and needs creating.

If the `Module` is there then it checks if its parameters match those defined in the `Ptemplate` Custom Resource and if necessary updates the `Module`.

Once everything matches what is expected reconciliation is complete and control is handed back to the main `Reconcile()` function.

If it needs to create the `Module` from scratch it calls `getNewModule()`. This simply builds a golang structure equivalent of the KMM `yaml` (or to be more precise, its an structure that meets the definitions in the imported `github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1/module_types.go` library, and the `Modules` yaml definition is a serialised version of the same definition).

The only "magic" in `getNewModule()` is the call

```go
controllerutil.SetControllerReference(ptResource, module, r.Scheme)
```

This marks the (not yet created) `Module` as owned by the `Ptemplate` resource and lets Kubernetes that this `Module` needs to be deleted when this `Ptemplate` is deleted. Without setting this you can end up with orphaned resources, and the kmod itself will be left loaded after it should be deleted.

Once the `Module` structure is created  then the `ReconcileModule()` function calls the Kubernetes API via the `r.Create()` call.

Updating an existing `Module` resource operates in a similar fashion except that instead of creating a new `Module` structure `ReconcileModule()` simply updates the one it received from `getExistingModule()` to match what it expects and updates Kubernetes via `r.Update()`.

Once the Kubernetes API server has been updated with the new, or changed `Module` resource details those details get stored in etcd, and Kubernetes will then take whatever steps in needs to make reality meet that etcd representation.

### ReconcileConsumer()

The `ReconcileConsumer()` function follows a very similar workflow to `ReconcileModule()`, except this time it checking for a DaemonSet that runs the "Consumer" pod on each node where the kmod is loaded.

So it sets up some defaults, gets the consumer DaemonSet if it exists (via `getExistingConsumer()`). If it doesn't exist it calls `r.getNewConsumer()` to create a new DaemonSet golang structure and fires off a `Create()` call to the API server to create it.

If `getExistingConsumer()` returns a structure then the DaemonSet already exists, so we check its parameters meet what is defined in the `Ptemplate` resource, and if not set them and update Kubernetes.

### SetupWithManager()

The final piece of "magic" is the `SetupWithManager()` function. This is very simple but provides a crucial piece of information to Kubernetes, specifically *when* to run the operators `Reconcile()` function.

```go
ctrl.NewControllerManagedBy(mgr).
    For(&cachev1alpha1.Ptemplate{}).
    Owns(&appsv1.DaemonSet{}).
    Complete(r)
```

This call defines that this Operator is for the `Ptemplate` resource and should be run whenever  a `Ptemplate` custom resource is added, updated, or changed. It also defines that this CR `Owns` a `DaemonSet` resource, and should therefore also be reconciled whenever a DaemonSet is updated, changed, or deleted. This allows the Ptemplate operator to notice when the DaemonSet it owns gets changed and fix it if necessary.
