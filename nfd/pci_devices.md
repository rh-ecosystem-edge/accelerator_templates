# I want to customise NFD to label nodes with a given PCI device

## Solution

NFD can be configured to add labels if a device with a given PCI class is present.

For example:

[See here](pci_devices.yaml)

```
apiVersion: nfd.kubernetes.io/v1
kind: NodeFeatureDiscovery
metadata:
  name: pt-nfd-config
  namespace: node-feature-discovery-operator
spec:
  instance: "" ## instance is empty by default
  operand:
    image: registry.k8s.io/nfd/node-feature-discovery:v0.13.4
    imagePullPolicy: Always
    servicePort: 12000
  workerConfig:
    configData: |
      core:
        labelSources: ["all"]
      sources:
        pci:
          deviceClassWhitelist:
            - "0100"
          deviceLabelFields:
            - "class"
            - "vendor"
            - "subsystem_vendor"
            - "subsystem_device"
```


This will label the node with
`feature.node.kubernetes.io/pci-0100_1af4_1042_1af4_1100.present=true` if the appropriate device exists (in this case a Red Hat, Inc. Virtio SCSI controller)

The `deviceLabelFields` entries determine what information the label encodes, in this case if a device with the whitelisted class exists it results in `pci-<class>_<vendor>_<subsystem_vendor>_<subsystem_device>.present=true`


## Discussion

Certain workloads required very specialist hardware to be available and for scheduling pods that perform those tasks being able to reference a label in their selector is vital.

The PCI class of a device provides the type of that device, so the `0100` device in our example above identifies it as a SCSI controllers, if its class was reported as `0201` it would be a Token Ring controller(!).  These break down into two codes, the first two numbers (`01`) refers to Storage Controllers, and the second two gives the subclass where `00` identifies SCSI devices. `02` identifies a network controller, and a subclass of `01` is a Token Ring controller (more usefully `0200` is the class of Ethernet controllers and `0207` InfiniBand controllers).

On a running Linux box the PCI classes of installed devices can be found either by the `lspci -nn` command or via the `/sys/bus/pci/devices/` directory structure which lists all the addresses on the PCI bus and the devices they address.

Sysfs files such as `/sys/bus/pci/devices/0000:04:00.0/class` contain a 6 digit hex string that gives the Class, Sub-class, and Programming Interface which can provide more information about the device type. NFD only uses the Class and Sub-class, it strips off the Programming Interface.


## Links
* [NFD PCI documentation](https://kubernetes-sigs.github.io/node-feature-discovery/v0.13/reference/worker-configuration-reference.html#sourcespci)
* [PCI Code And ID Assignment Specification](https://pcisig.com/sites/default/files/files/PCI_Code-ID_r_1_12__v9_Jan_2020.pdf)
