# Openshift Partner Templates

Adding hardware support to Openshift clusters is more complicated than on non-containerized Linux. While CoreOS, the operating system Openshift sits upon, is essentially Red Hat Enterprise Linux (RHEL) its appliance nature gives it two important differences.

Firstly access to the operating system is restricted. While it is possible to connect to an Openshiftnode via ssh and change to the root user this behaviour is discouraged, and installing software on the underlying CoreOS operating system will lead to invalidating support contracts.

Secondly CoreOS is designed to be immutable. This means rebooting the node deletes any changes made to the operating system.

Together these make the traditional approach of installing packages of drivers and agent software that system administrators can install problematic. Instead a new approach of containerised drivers managed by Openshift operators is required.

Partner Templates will explore the various aspects of this new approach, and attempt to provide examples of how top approach various problems that might be encountered.

To achieve this we break the problem down into a number of chapters. Not all solutions will need to use all of these and we have attempted to make each section as stand alone as possible. 

## Table Of Contents

1. [Kernel Modules](kernel_module/README.md)

1. [Driver Containers](driver_container/README.md)

1. [The Kernel Module Management (kmm) operator](kmm/README.md)

1. [The Node Feature Discovery (nfd) operator](nfd/README.md)

1. [Device Plugin](device_plugin/README.md)

1. operators

1. Certification

1. [Glossary Of Terms](GLOSSARY.md)



## Related Operators

[Intel Technology Enabling for OpenShift](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main) with it's related [device plugins](https://github.com/intel/intel-technology-enabling-for-openshift/tree/main)





