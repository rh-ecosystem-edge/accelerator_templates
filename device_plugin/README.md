# Device Plugins

To provide more reliable scheduling and minimise node resource over commitment Kubernetes uses the concept of resources, and allows a pod to reserve a portion of the CPU and memory resources for its use. However the resources provided by the core of Kubernetes are limited to only CPU and Memory because they are the only resources that it can reasonably rely on existing.

To allow the addition of custom resource types Kubernetes provides the device plug-in API, and the creation of extensions called Device Plugins

Device Plugins provide a way to extend the resource concept to allow pods to reserve custom resources including devices. Instead of customising the code for Kubernetes itself, vendors can implement a device plugin that you deploy either manually or as a DaemonSet. The targeted devices include GPUs, high-performance network cards, FPGAs, InfiniBand adaptors, and other similar computing resources that may require vendor specific initialisation and setup.

## Implementing A Device Plugins

The kubelet exports a Registration gRPC service which the device plugin uses to register itself.

Once registered it then needs to implement the DevicePlugin service

```golang
service DevicePlugin {
      rpc GetDevicePluginOptions(Empty) returns (DevicePluginOptions) {}
      rpc ListAndWatch(Empty) returns (stream ListAndWatchResponse) {}
      rpc Allocate(AllocateRequest) returns (AllocateResponse) {}
      rpc GetPreferredAllocation(PreferredAllocationRequest) returns (PreferredAllocationResponse) {}
      rpc PreStartContainer(PreStartContainerRequest) returns (PreStartContainerResponse) {}
}
```

Of which the most important functions are `ListAndWatch()`, and `Allocate()`

### ListAndWatch()

The `ListAndWatch()` function is called by kubelet to get the status of the devices on the node, their number, their IDs and their health.

This data will be propagated by kubelet to the Kubernetes API server, and will be used by the Kubernetes scheduler to decide where to schedule any pod that requests the resource.

`ListAndWatch()` responds to the kubelet request with an array of devices, each represented as a `struct` containing an ID for the device and its Health. The ID can be set to any value by the device-plugin, as long as it knows how to correlate them to the actual devices. Those IDs will be used by the kubelet to decide which device(s) should be allocated to each pod.

### Allocate(AllocateRequest)

The `Allocate()` function is called by the kubelet prior to scheduling pod on the node. The kubelet sends the IDs of all devices that kubelet wants to allocated for the scheduled pod, and `Aloocate()` function can then

* verify the state/health of each requested device and return an error if it's not available, causing the requesting pod to fail to schedule.

* It can add specific annotations to the response. Those annotation will be set by kubelet on the container using container runtime

* It can add environment variable to the response. Those variable will be passed  on to the pod by the kubelet so can include such things as device node names, PCI buss addresses or any thing else that allows the pod to know details about precisely what has been allocated.

* It can also define mounting of the devices on the container file-system. The  volume/mount definitions will be added by the kubelet to the container. This means that the Pod does not necessarily need to mount /dev host directory it can be mounted for it by the kubelet

For more details see the [Kubernetes Device Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/) page, and the examples listed in [Links](#links) below.

## Cookbook

* [I want to deploy my device plugin alongside my driver](kmm_with_device_plugin.md)

## Links

* [Kubernetes Device Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/)

* [OpenShift: Using device plugins to access external resources with pods](https://docs.openshift.com/container-platform/4.13/nodes/pods/nodes-pods-plugins.html) (Includes links to several example device plugins)

* [Partner Templates Device Plugin example](../src/ptemplate-device-plugin/main.go)

* [Simple Device Plugin](https://github.com/yevgeny-shnaidman/simple-device-plugin/) an example device plugin configurable via a file.

* [K8s Dummy Device Plugin](https://github.com/redhat-nfvpe/k8s-dummy-device-plugin) (for testing purpose only)

* [The NVIDIA device plugin for Kubernetes](https://github.com/NVIDIA/k8s-device-plugin)
