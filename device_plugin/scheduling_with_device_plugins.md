# I want to schedule a worker pod only where there resources

## Problem

Hardware resources on a node are limited, you need to schedule pods only on nodes that have the required resources available.

## Solution

Once a device plugin has been set up to advertise the resources then pods can consume them using the `resources` stanza in their `container` definition

e.g. ([see full file](consumer.yaml)):

```yaml
spec:
  containers:
      ...
      resources:
        limits:
          example.com/ptemplate: 1
```

When this pod is created it will request one example.com/ptemplate resource from the Kubelet, the Kubelet will then contact the device plugin for the node and if one is available it will be allocated to the pod. If not the Kubelet will not run the pod.

Multiple instances of the same resource can be allocated, although the container has no control over which resources it receives from the pool
