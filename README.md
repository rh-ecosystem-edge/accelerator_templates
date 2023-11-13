# OpenShift Partner Templates

Adding hardware support to OpenShift clusters is more complicated than on non-containerised Linux. While Red Hat CoreOS (RHCOS), the operating system OpenShift sits upon, is essentially Red Hat Enterprise Linux (RHEL) its appliance nature gives it two important differences.

Firstly access to the operating system is restricted. While it is possible to connect to an OpenShift node via ssh and change to the root user this behaviour is discouraged, and installing software on the underlying RHCOS operating system will lead to invalidating support contracts.

Secondly RHCOS is designed to be immutable. This means rebooting the node deletes any changes made to the operating system.

Together these make the traditional approach of installing packages of drivers and agent software that system administrators can install problematic. Instead a new approach of containerised drivers managed by OpenShift operators is required.


Partner Templates is made up of three sections

* [The PTemplate Operator](src/README.md): An example driver, Operator, and other code to provide a working example of one approach to implementing a driver solution on OpenShift or Kubernetes

* [The "Partner Templates Cookbook"](#Cookbook): A collection of solutions, links, and FAQs to provide help for developers getting started with adding OpenShift support to their hardware.

* [Links and Related Operators](#links-and-related-operators): A collections of Links to useful resources and real-world implementations to provide examples of how other people have approached the subject.

&nbsp;

&nbsp;

**NOTE: If there is something you think needs adding, expanding, or correcting, please file an Issue, or even better raise a PR**

&nbsp;

## Cookbook

Hardware Operators are built up from a number of components not all solutions will need to use all of these and we have attempted to make each section as stand alone as possible.

### Table Of Contents

1. [Kernel Modules](kernel_module/README.md)

1. [Driver Containers](driver_container/README.md)

1. [The Kernel Module Management (kmm) operator](kmm/README.md)

1. [The Node Feature Discovery (nfd) operator](nfd/README.md)

1. [Device Plugin](device_plugin/README.md)

1. [Operator](operator/README.md)

1. [Certification For Containers and Operators](certfication/README.md)

1. [Support](support.md)

**Appendices**

1. [Checklist](Checklist.md)

1. [Glossary Of Terms](GLOSSARY.md)



&nbsp;

&nbsp;


## Links and Related Operators

[Glossary](glossary.md) a collection of useful terms.

[Red Hat OpenShift Container Platform Life Cycle Policy](https://access.redhat.com/support/policy/updates/openshift)

[Intel Technology Enabling for OpenShift](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main) with it's related [device plugins](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main)





