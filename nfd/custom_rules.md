# I want to add my NFD config without interfering with existing setup

## Solution

The NodeFeatureRule custom resource allows the adding of multiple resources that instruct NFD what labels to add. These can then be used to add labels without disrupting any existing NFD configuration.

For example:

```
apiVersion: nfd.k8s-sigs.io/v1alpha1
kind: NodeFeatureRule
metadata:
  name: my-sample-rule-object
spec:
  rules:
    - name: "is fuse installed?"
      labels:
        "fuse-kmod": "true"
      matchFeatures:
        - feature: kernel.loadedmodule
          matchExpressions:
            fuse: {op: Exists}
```



this will add the label `feature.node.kubernetes.io/fuse-kmod=true`  if the fuse kernel module is installed regardless of any other configuration (including any other  configuration that looks for the fuse kernel module)

It will add this label as long as the NFD operator is installed, regardless of other configuration it may have.

## Discussion 

NodeFeatureRule objects provide an easy way to create vendor or application specific labels. It uses a flexible rule-based mechanism for creating labels based on node features. It also makes it easier to use custom label names within the `feature.node.kubernetes.io` namespace.

These rules can be used for anything that the NodeFeatureDiscovery resource can label but has much greater flexability, not only allowing customised labels, but also allowing complex rules to be used with labels that depend on multiple features being present.

For example

```
apiVersion: nfd.k8s-sigs.io/v1alpha1
kind: NodeFeatureRule
metadata:
  name: redhat-scsi-controller
spec:
  rules:
    - name: "my scsi rule"
      labels:
        "redhat-scsi-controller": "true"
      matchFeatures:
        - feature: pci.device
          matchExpressions:
            class: {op: In, value: ["0100"]}
            vendor: {op: In, value: ["1af4"]}

```

Which adds a label `feature.node.kubernetes.io/redhat-scsi-controller=true`
if a pci device with the class of `0100` (a scsi controller) and also a vendor code of `1af4` (Red Hat) exists on the system. We could allow for multiple classes from the same vendor by adding to the `class` expression, e.g

```
class: {op: In, value: ["0100", "0300"]}
```

This means that if a device from either of those classes exists, that also has the vendor id of `1af4` the label will be added.


It's also possible to create rules with many expressions utilising features from many features, with a whole range of different operators, for example:

```
apiVersion: nfd.k8s-sigs.io/v1alpha1
kind: NodeFeatureRule
metadata:
  name: my-sample-rule-object
spec:
  rules:
    - name: "my sample rule"
      labels:
        "my-sample-feature": "true"
      matchFeatures:
        - feature: kernel.loadedmodule
          matchExpressions:
            dummy: {op: Exists}
        - feature: kernel.config
          matchExpressions:
            X86: {op: In, value: ["y"]}
```





## Links
* [NFD NodeFeatureRule documentation](https://kubernetes-sigs.github.io/node-feature-discovery/v0.14/usage/customization-guide.html#nodefeaturerule-custom-resource)
