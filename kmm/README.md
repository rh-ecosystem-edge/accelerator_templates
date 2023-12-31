# Kernel Module Management (KMM) Operator

## Introduction

The Kernel Module Management Operator is designed too manage out-of-tree kernel modules in Kubernetes and OpenShift. It does this by managing the loading and unloading of driver-containers across the cluster using node selectors to determine which nodes require the driver deployed. Optionally it can also manage the building of driver containers and signing the kernel modules themselves for SecureBoot installations.

The KMM operator implements its own Custom Resource definition for resources of `kind: Module`. When a `Module` resource is defined KMM creates a DaemonSet that runs the referenced driver container on each of the nodes with the command

```bash
/bin/sh -c modprobe -v -d "/opt" "[driver_name]"
```

This loads the kmod on the referenced nodes. The kmod will now appear in the output of `lsmod` for the node and can be used just like any other loaded kmod.

When the Module resource is deleted then the driver container pod is deleted and the kmod is unloaded via the `-r` argument to the `modprobe` command

```bash
/bin/sh -c modprobe -rv -d "/opt" "[driver_name]"
```

To check what Modules have been created in a cluster you can use `kubectl get module` and describe them as you would any other Kubernetes resource with `kubectl describe module pt-char-dev`

## Cookbook

* [I want to deploy a prebuilt driver container with KMM](load_module.md)

* [I want to manage which nodes my driver is loaded on](node_selectors.md)

* [I want to load different drivers on different kernels](different_kernels.md)

* [I need KMM to build my driver containers for me](build_module.md)

* [I use SecureBoot and need my drivers signed before loading](secureboot.md)

* [I want to load multiple kernel modules]

* [I want to load custom firmware with my driver](loading_firmware.md)

* [I need to unload an in-tree driver and replace it with my new out-of-tree version]

## Links

* [KMM's Upstream documentation](https://kmm.sigs.k8s.io)

* [OpenShift NFD documentation](https://docs.openshift.com/container-platform/4.13/hardware_enablement/psap-node-feature-discovery-operator.html)

* [Upstream Repository](https://github.com/kubernetes-sigs/kernel-module-management/) where most development work is done.

* [KMM's OpenShift Repository](https://github.com/rh-ecosystem-edge/kernel-module-management/) forms the basis for the OpenShift release of KMM. Upstream changes are merged into here regularly, and OpenShift specific development and testing is performed.
