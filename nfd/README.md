# Node Feature Discovery Operator

Kubernetes uses labels to store and display metadata about its objects. As nodes are, within the kuberenets database, just another type of object they too can have labels attached to them, and these can then be used to distinguish features such as which nodes have particular type of PCI card installed, or CPU feature available.

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


Labels can be added to each node manually or the Node Feature Discovery Operator (NFD) manages the detection of hardware features and configuration in an OpenShift Container Platform cluster and labels the nodes with hardware-specific information.


## Cookbook

* [I want NFD to label nodes if a CPU supports a given flag](cpu_flags.md)

* [I want to customise NFD to label nodes with a given PCI device](pci_devices.md)

* [I want to add my NFD config to an existing setup easily](custom_rules.md)

* [I want KMM to load kernel modules based on NFD labels](kmm.md)

## More Reading
* [Openshift Documentation](https://docs.openshift.com/container-platform/4.13/hardware_enablement/psap-node-feature-discovery-operator.html)

* [NFD documentation](https://kubernetes-sigs.github.io/node-feature-discovery/v0.13/get-started/index.html)