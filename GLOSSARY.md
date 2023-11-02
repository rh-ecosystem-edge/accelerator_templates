# Glossary

### Driver Container
Container images used for building and deploying out-of-tree kernel modules and drivers on container operating systems like Red Hat Enterprise Linux CoreOS (RHCOS). They have an `entrypoint` that loads a device driver (kmod) on startup.

### kmod
A kernel module. A binary file, normally with a .ko file extension that provides a device driver or other piece of optional functionality that can be loaded into the Linus Kernel if required.

### .ko file
A kernel object file. A kernel module compiled into a single file ready for loading into the kernel

### kmm
The Kernel Module Management operator used for managing the deployment of kernel modules across OpenShift clusters using [driver containers](#driver container).

### nfd
The Node Feature Discovery operator and components used for adding a wide range of hardware labels to OpenShift nodes

### Out-of-tree (OOT) drivers
Drivers and kernel modules that are maintained outside of the Linux source tree. They are distributed outside of the kernel source and maintained by third parties (who may be a Linux distribution provider or more commonly a hardware provider). As the kernel is updated the third-party is responsible for updating the driver to ensure compatibility.

### in-tree drivers
Drivers and kernel modules maintained as part of the Linux source tree with their source code available from kernel.org under the terms of the GPL. As part of the Linux source they are updated to ensure they work with the version of the kernel.
