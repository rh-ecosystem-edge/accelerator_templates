# I want to schedule a worker pod only where there resources

## Problem

Hardware resources on a node are limited, you need to schedule pods only on nodes that have the required resources available.



## Solution

Once a device plugin has been set up to advertise the resources then pods can consume them using the `resources` stanza in their `container` definition


e.g. ([see full file](consumer.yaml)):

```
spec:
  containers:
      ...
      resources:
        limits:
          example.com/ptemplate: 1
```

When this pod is created it will request one example.com/ptemplate resource from the kublet, the kubelet will then contact the device plugin for the node and 






```
---
apiVersion: v1
kind: Pod
metadata:
  name: consumer-pod
spec:
  containers:
    - name: consumer
      image: quay.io/chrisp262/pt-device-plugin:consumer-latest
      securityContext:
        privileged: true
      resources:
        limits:
          example.com/ptemplate: 2
      volumeMounts:
        - mountPath: /hostdev
          name: host-dev
          readOnly: false
  volumes:
    - hostPath:
        path: /dev
      name: host-dev
```


## Discussion