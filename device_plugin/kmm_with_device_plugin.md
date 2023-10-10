# I want to deploy my device plugin alongside my driver

## Problem

You want to deploy your device plugin alongside the driver it reports on, and have the two deploy and un-deploy together.

## Solution

KMM can deploy a device plugin as well as a driver container as part of its `Module` resource.

Add a `spec.devicePlugin` stanza to an existing `Module` resource that deploys the kernel module.

e.g. ([see full file](kmm_with_device_plugin.yaml)):

```
spec:
  devicePlugin:
    container:
      image: quay.io/chrisp262/pt-device-plugin:latest
```


## Discussion

Alongside deploying a kernel module the KMM operator can create a daemonset that manages a device plugin pod on each of the nodes in the cluster. A pod will be created on every node in the cluster that matches the `spec.selector` field in the `module` resource. This pod will run the device plugin and register its resources with kubernetes.


As the device plugin needs to be able to get information about the underlying host and its hardware it's quite common to need to mount a local volume within the device plugin container. This can be done by creating `volume` object with in the `spec.devicePlugin.volume` section of the module resource and then setting up the required `volumeMounts` within the device plugin container section at `spec.devicePlugin.container.volumeMounts`. Both these section are standard kuberentes objects and can take all the fields that you would expect from any other `volume` and `volumemount`.

It is also possibel to runt he device plugin under a specific service account by setting the `spec.devicePlugin.serviceAccountName` 

