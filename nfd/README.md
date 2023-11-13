# Node Feature Discovery Operator

### The Problem

By default OpenShift can schedule a pod to run on any node in the cluster. This can be limited by the use node selectors that limit the pods to only running on nodes that have a matching label. Kernel modules run in the kernel so if they encounter an error they can lead to a kernel panic and the node failing. Loading a kernel module on a node without the hardware it controls is at best pointless and at worst can lead to significant failures.

Therefore discovering what hardware is available on a node and labelling it accordingly is important for the use of Driver Containers, but doing this by hand is time consuming and error prone.

### Node Feature Discovery Operator


[Node Feature Discovery (NFD)](https://github.com/kubernetes-sigs/node-feature-discovery) is a Kubernetes add-on for detecting hardware features and system configuration, and then labelling the nodes with that information.

It is deployed as an operator, this deploys a worker pod on each node in the cluster. The worker pods periodically check the hardware features and configuration of the node and report back to a master pod. The master pod then adds and removes labels for each node based on this information.


### Usage

[The OpenShift documentation](https://docs.openshift.com/container-platform/4.13/hardware_enablement/psap-node-feature-discovery-operator.html) describes the best way to deploy NFD and get started with labelling nodes

```
Labels are key/value pairs, and must obey some rules about their format.

The keys are made up of an optional prefix and a name, separated by a slash (/).

* The prefix is optional and can be up to 253 characters in length before the slash.
* The name segment is required and must be 63 characters or less,
* They can contain alphanumeric character ([a-z0-9A-Z]), dashes (-), underscores (_), dots (.)
* They must begin and end with an alphanumeric.

The kubernetes.io/ and k8s.io/ prefixes are reserved for Kubernetes core components.

Valid label value:
* must be 63 characters or less,but can be empty (often depicted as "")
* can contain alphanumeric character ([a-z0-9A-Z]), dashes (-), underscores (_), dots (.)
* if not empty must begin and end with an alphanumeric.
```


## Cookbook

* [I want NFD to label nodes if a CPU supports a given flag](cpu_flags.md)

* [I want to customise NFD to label nodes with a given PCI device](pci_devices.md)

* [I want to add my NFD configuration to an existing setup easily](custom_rules.md)

* [I want KMM to load kernel modules based on NFD labels](kmm.md)

## More Reading
* [OpenShift Documentation](https://docs.openshift.com/container-platform/4.13/hardware_enablement/psap-node-feature-discovery-operator.html)

* [NFD documentation](https://kubernetes-sigs.github.io/node-feature-discovery/v0.13/get-started/index.html)
