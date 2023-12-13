# Accelerator Templates

The Accelerator Templates Operator is a set of components to demonstrate the deploying of drivers and their associated workload on OpenShift it includes

- The ptemplate_char_dev kernel driver that creates a number of character devices (/dev/ptemplate-X) that you you can write up 100 bytes (configurable via a parameter) to, and read them back.

- A Driver Toolkit based Dockerfile to build a container that will load the ptemplate_char_dev kmod

- A Node Feature Discovery configuration that will label the nodes on which the ptemplate kmod should be loaded

- A device plugin that will report to OpenShift how man /dev/ptemplate-X devices are available and allow Pods to request exclusive access to them via Kubernetes Resource Requests.

- a Consumer program and Dockerfile that will make use of the /dev/ptemplate devices the device plugin allocates to it

- An Operator that will create a Module resource that the Kernel-Module-Management operator will use to manage loading and unloading the kernel module, and a DaemonSet resource that will deploy a 'consumer' pod on every node to read and write to the /dev/ptemplate devices.

**NOTE: This is a learning tool, designed to explain the principles and the various components that can be used together to build a solution, the code here is not production quality and is NOT supported by Red Hat or anyone else.**

**SECOND NOTE: If you find something you think could be better feel free to file an Issue in github, or (even better) a pull request. The aim is to make this generally useful as a learning tool so any improvements are welcomed!**

# Installation

## Install Prerequisites

### Node Feature Discovery Operator
[Node Feature Discovery (NFD) Operator](https://docs.openshift.com/container-platform/4.14/hardware_enablement/psap-node-feature-discovery-operator.html) manages the deployment and life cycle of the NFD add-on to detect hardware features and system configuration, such as PCI cards, kernel, operating system version, etc.

#### Install NFD Operator
Follow the guide below to install the NFD operator using the command line or web console. 

- [Install from the CLI](https://docs.openshift.com/container-platform/4.14/hardware_enablement/psap-node-feature-discovery-operator.html#install-operator-cli_node-feature-discovery-operator)
- [Install from the web console](https://docs.openshift.com/container-platform/4.14/hardware_enablement/psap-node-feature-discovery-operator.html#install-operator-web-console_node-feature-discovery-operator)

### Install Kernel-Module-Management Operator

The Kernel Module Management Operator manages the loading and unloading of kernel modules on the nodes, as well as building out-of-tree drivers from source code if required.

- [Install KMM On Kubernetes](https://kmm.sigs.k8s.io/documentation/install/)

- [Install On OpenShift](https://openshift-kmm.netlify.app/documentation/install/)

