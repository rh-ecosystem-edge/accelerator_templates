# OpenShift Partner Templates

Adding hardware support to OpenShift clusters is more complicated than on non-containerised Linux. While Red Hat CoreOS (RHCOS), the operating system OpenShift sits upon, is essentially Red Hat Enterprise Linux (RHEL) its appliance nature gives it two important differences.

Firstly RHCOS is designed to be immutable. This means rebooting the node deletes any changes made to the operating system. Secondly access to the operating system is restricted. While it is possible to connect to an OpenShift node via ssh and change to the root user this behaviour is discouraged, and installing software on the underlying RHCOS operating system will lead to invalidating support contracts.

Together these make the traditional approach of installing packages of drivers and agent software that system administrators can install problematic. Instead a new approach of containerised drivers managed by OpenShift operators is required.

To support this Red Hat has developed a number of technologies:

* [Driver Toolkit (DTK)](driver_container/README.md) - A toolkit to help create OCI images containing kernel modules and drivers (known as driver containers)

* [Node Feature Discovery (NFD) Operator](nfd/README.md) - An operator and application that labels nodes according to hardware and operating system features

* [Kernel Module Management (KMM) Operator](kmm/README.md) - An operator to load driver containers on nodes that meet set criteria (normally nodes with the given hardware)

Together these provide a rich set of tools for third party developers to build on to support their own drivers with custom operators and other tooling.

## Suggested workflow

Every operator is different and will need different components so the steps required to build the solution will be different, but the following checklist should provide a good starting point for most projects.

* [] Work through the [Operator Checklist](checklist.md) to assess what work has already been done and what is required to be done before shipping.

* [] Create the Device Driver and any user land tools. This is the same as for any other version of Linux including RHEL. ([Example](src/kernel_module/README.md)  [Source Code](src/kernel_module/))

* [] Package the device driver into a [Driver Container](driver_container/README.md) using the  [Driver Tool Kit (DTK)](https://github.com/openshift/driver-toolkit) a toolkit to help create OCI images for kernel modules and drivers.  ([Example Source](src/driver_container))

* [] Package any user land component required into container images via a [Dockerfile](https://docs.docker.com/engine/reference/builder/) and [podman build](https://docs.podman.io/en/latest/markdown/podman-build.1.html)

* [] A config for the [Node Feature Discovery (NFD)](nfd/README.md) operator, labelling nodes based on its hardware and operating system features.  ([Example Source](src/nfd/pci_devices.yaml))

* [] [Create a Device Plugin](device_plugin/README.md) to allow the user land components to request the hardware and make sure they are not being oversubscribed. ([Example Source](src/ptemplate-device-plugin/))

* [] Create a [DaemonSet configuration](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) (as a yaml file), to deploy the user land components. The yaml created here is useful for testing, but will also translate directly into the golang structures that the operator will use to instantiate Kubernetes objects.

* [] Create a configuration for the [Kernel Module Management](kmm/README.md) operator to load driver containers on nodes that meet set criteria (normally nodes with the given hardware).  Again building this as a manually applied yaml file is both a good sanity check that all the parts work together manually before they get automated with an Operator, and translate directly into the golang structures the operator needs.

* [] Automate the deployment of the components by creating a custom Operator that deploys the Driver Container and the user land components it needs. ([Operators](operator/README.md) [Integration with KMM](integration/README.md) [Example source](src/ptemplate-operator/)

* [] Add [metrics](observability/README.md) the operator to report the state of the hardware to the cluster manager.

* [] Create any PrometheusRule configuration yaml that might be needed.

* [] Add [Grafana dashboards](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/create-dashboard/) and any other supporting components needed to make the cluster operator's life easier.

* [] [Certify your driver](certification/README.md). As only a specific version of the images are certified this needs to be the final step before release.

&nbsp;

## Direct Links

1. [Kernel Modules](src/kernel_module/README.md)

1. [Driver Containers](driver_container/README.md)

1. [The Kernel Module Management (kmm) operator](kmm/README.md)

1. [The Node Feature Discovery (nfd) operator](nfd/README.md)

1. [Device Plugin](device_plugin/README.md)

1. [Operator](operator/README.md)

1. [Integrating with KMM](integration/README.md)

1. [Observability and Metrics](observability/README.md)

1. [Certification For Containers and Operators](certification/README.md)

1. [Support](support.md)

### Appendices

1. [Checklist](checklist.md)

1. [Glossary Of Terms](GLOSSARY.md)

## Links and Related Operators

[Red Hat OpenShift Container Platform Life Cycle Policy](https://access.redhat.com/support/policy/updates/openshift)

[Intel Technology Enabling for OpenShift](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main) with it's related [device plugins](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main)

&nbsp;

## Corrections and Omissions

If there is something you think needs adding, expanding, or correcting, please file an Issue, or even better raise a PR
