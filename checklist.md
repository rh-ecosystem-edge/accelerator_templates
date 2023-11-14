# Operator Requirements

To capture requirements for operator development and enablement in OpenShift it is helpful to ask a series of questions. Not all are relevant to every solution but they are helpful for thinking about the components that make up the solution and how they will work together.

## Kernel Modules

* What is available/supported in-tree for which adaptors?
* What is needed for out-of-tree support for which adaptors? What kernel modules/drivers and how many?
* What are the kmod dependencies? Is there a symbol-dependency between OOT kernel module(s) and in-tree kernel modules that are not loaded on boot?
* If there is more than one OOT kernel module that needs to be loaded, is there a symbol-dependency between them?
* If there is more than one OOT kernel module that needs to be loaded, is there a logical dependency between them (do they need to be loaded in a specific order)?
* Is there a need for firmware loading for the device?
* Is there a need for a user-space application to be executed by the kernel module in user-space/container?
* Are there any proprietary libraries that need to be considered?
* Are those modules going to be in-cluster built, or off-cluster built only?

## Device plugin

* Does this already exist?
* What does kubelet need to know about the device?

## PCIe

* How do these devices appear in the OS from a hardware perspective?
* Can this be used for node labelling?

## Workload Orchestration

* Is there an existing application that does this today?
* How is this handled?
