# Driver Container

## Introduction

To use a kernel module on Openshift it needs to be built into a container image which Openshift can then use to instantiate a pod. These images are commonly known as Driver Containers. 

Driver containers are container images used for building and deploying out-of-tree kernel modules and drivers on container OSs like Red Hat Enterprise Linux CoreOS (RHCOS).

Driver containers are the first layer of the software stack used to enable kernel modules on Kubernetes.

Building a kernel module inside a Dockerfile is made more complicated by the large number of additional packages that are required for the build process (e.g. a C compiler). Getting these right can be an annoying process of trial and error, and leaving them in a driver-container that is deployed in production both bloats the container image (often by hundreds of MB), and leaves an additional attack surface for no value. The Driver Toolkit provides a workaround for this



## Cookbook

1. [I need to build a driver container from scratch](driver_container_hard_way.md)

1. [I want an easier way to build a driver container, the Driver Toolkit](driver_container_easy_way.md)

1. [I need to run my driver container](driver_container_running.md)

1. I need my driver at boot time

## Links
* https://github.com/openshift/driver-toolkit
