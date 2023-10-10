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
      image: quay.io/chrisp262/simple-device-plugin:latest
```



## Discussion

Alongside deploying a kernel module the KMM operator can create a daemonset that manages a device plugin pod on each of the nodes in the cluster. 

