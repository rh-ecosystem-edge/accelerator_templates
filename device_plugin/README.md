# Device Plugins

## Introduction

To provide more reliable scheduling and minimize node resource overcommitment kubernetes uses the concept of resources, and allows a pod to reserve a portion of the CPU and memory resources for its use. However the resources provided by the core of kuberenetes are limited to only CPU and Memory because they are the only resources that it can reasonably rely on existing.

To allow the addition of custom resource types Kuberenetes provides the device plug-in API, and the creation of extensions called Device Plugins

Device Plugins provide a way to extend the resource concept to allow pods to reserve custom resources including devices. Instead of customizing the code for Kubernetes itself, vendors can implement a device plugin that you deploy either manually or as a DaemonSet. The targeted devices include GPUs, high-performance NICs, FPGAs, InfiniBand adapters, and other similar computing resources that may require vendor specific initialization and setup.


## Cookbook
* [Question: who supports a device plugin?](support.md)
* [I want to write a device plugin for my driver](writing_a_device_plugin.md)
* [I want to deploy my device plugin alongside my driver](kmm_with_device_plugin.md)
* [I want to schedule a worker pod only where there resources](scheduling_with_device_plugins.md)


## Links

* [Openshift: Using device plugins to access external resources with pods](https://docs.openshift.com/container-platform/4.13/nodes/pods/nodes-pods-plugins.html) (Includes links to several example device plugins)

* [Simple Device Plugin](https://github.com/yevgeny-shnaidman/simple-device-plugin/) an example device plugin configurable via a file.

* [K8s Dummy Device Plugin](https://github.com/redhat-nfvpe/k8s-dummy-device-plugin) (for testing purpose only)

* [The NVIDIA device plugin for Kubernetes ](https://github.com/NVIDIA/k8s-device-plugin)
