# I want KMM to load kernel modules based on NFD labels

## Problem

You want to use the labels NFD creates to determine where your driver is loaded.

## Solution

[KMM](../kmm/README.md) provides a `.selector` field to select which nodes its kernel module payload is deployed on.


For example ([see full file](kmm.yaml)):

```
  selector:
    feature.node.kubernetes.io/cpu-model.vendor_id: Intel
    feature.node.kubernetes.io/cpu-model.family: 6
    feature.node.kubernetes.io/cpu-cpuid.VMX: true
    feature.node.kubernetes.io/pci-0100_1af4_1042_1af4_1100.present: true
```

will load the kernel module on any node NFD has labelled as having a sixth generation Intel processor (actually a rather elderly Kaby Lake processor!), that supports the VMX extensions, and has a pci device that matches the class and ID given (this is actually a "Red Hat, Inc. Virtio network device", these labels come from a virtual machine).

## Discussion

The KMM selectors form a logical AND, the node must have all the labels for the driver to be deployed. This is great for restricting the kmod deployment quite precisely, but less good if you want to deploy it on a range of machines. In the above example the kmod will only be deployed on a machine with `feature.node.kubernetes.io/cpu-model.family: 6` but more normally you would want to deploy it on a range of CPU families, like "family 6 or above", or "CPUs between  6 and 8".

For these more complex setups you can use NFD rules to add a specific label based around the rules, for example (see [file](kmm_nfd_rule.yaml)):

```
apiVersion: nfd.k8s-sigs.io/v1alpha1
kind: NodeFeatureRule
metadata:
  name: ptemplate-required-rule
spec:
  rules:
    - name: "ptemplate_required"
      labels:
        "ptemplate": "required"
      matchFeatures:
        - feature: cpu.model
          matchExpressions:
            family: {op: In, value: ["6", "7","8"]}
```



## Links

* [KMM working with labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
* [Kubernetes Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)