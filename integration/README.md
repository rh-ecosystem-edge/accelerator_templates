# Using KMM with other operators

1. [Introduction](#introduction)

2. [Loading a driver](#loading-a-driver)

3. [Unloading a driver](#unloading-a-driver)

4. [Example Implementation](ptemplate-implementation.md)

5. [Links](#links)

## Introduction

Often loading a kernel module by itself is only part of the problem, requiring other resources like user land processes, or carefully calculated kmod parameters, to provide their full benefit.

The recommended way to tie together the loading of kernel modules with other OpenShift resources is to create a device specific operator with its own Custom Resources that then creates and manages the other Kubernetes objects it needs. The users of your new operator can then create new instances of your Custom Resources, and your operator will then process this to automatically create Module resources for the KMM operator, DaemonSets to create pods, and any other resources needed.

The [Ptemplate Operator](../src/ptemplate-operator) is an example of a way to implement this. It is mostly designed to be a simple teaching tool rather than a useful operator but it does load a kmod and create a set of pods to utilise it. See [here](ptemplate-implementation.md) for more details on how it is implemented.

## Loading a driver

Loading and unloading a driver manually requires the creation and removal of a `Module` resource, by hand this would require an `oc create` or `oc remove` command, to do the same from an operator requires using the Kubernetes API to create and destroy the resource.

Once the `Module` object is created the KMM operator will take over responsibility for the actual loading of the kmod.

To create a `Module` resource from an operator all that is needed is to create a golang structure that meets the `type Module struct` defined in the [KMM API source code](https://github.com/rh-ecosystem-edge/kernel-module-management/blob/main/api/v1beta1/module_types.go), and then call the `Create` method of the Kubernetes API.

For example:

```golang
import kmmv1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"

module := &kmmv1beta1.Module{
    ObjectMeta: metav1.ObjectMeta{
        Name:      myModuleResourceName
        Namespace: myModuleResourceNamespace,
    },
    Spec: kmmv1beta1.ModuleSpec{
        ModuleLoader: kmmv1beta1.ModuleLoaderSpec{
            Container: kmmv1beta1.ModuleLoaderContainerSpec{
                Modprobe: kmmv1beta1.ModprobeSpec{
                    ModuleName: myKmodBaseName,
                },
                ImagePullPolicy: "Always",
                KernelMappings: []kmmv1beta1.KernelMapping{
                    {
                        Regexp:         "^.+$",
                        ContainerImage: myDriverContainerImage
                    },
                },
            },
        },
        ImageRepoSecret: myImagePullSecret,
        Selector:        mySelector,
    },
}

r.Create(ctx, module)
```

## Unloading a driver

There are two possibilities for Unloading a kernel module, it can either be done directly from the code, or indirectly by declaring the `Module` as a child of another resource.

Setting the `Module` resource to be a child is as simple as calling the `SetControllerReference()` function on it:

```golang
controllerutil.SetControllerReference(myParentResource, myModuleResource, r.Scheme)
```

Because in Kubernetes deleting a parent resource will cause its children to be deleted first, setting this reference is all that is required to unload the kmod when the parent is deleted.

The other method for `Module` deletion is to manually delete it via the `Delete()` method of the Kubernetes Client:

```golang
r.Delete(ctx, myModuleResource)
```

Care must obviously be taken to ensure that the Operator doesn't not then attempt to recreate it the next time `Reconcile()` is run.

## Links

Walk through of the `Reconcile()` implementation for the [Ptemplate Operator](../src/ptemplate-operator)
