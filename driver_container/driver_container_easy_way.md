# I need an easier way to build a driver container

## Problem

Building a driver container from scratch is hard to get right, I need an easier, more reliable way to do it.

## Solution


The Driver Toolkit (DTK from now on) is a container image in the OpenShift payload which is meant to be used as a base image on which to build driver containers. The Driver Toolkit image contains the all the commonly needed packages to build or install kernel modules as well as a few tools needed in driver containers. The version of these packages will match the kernel version running on the RHCOS nodes in the corresponding OpenShift release.

By using it as a base image for the builder stage our Dockerfile  for the ptemplate_char_dev driver container can be simplified like this:

```
ARG KERNEL_VERSION
ARG DTK=registry.redhat.io/openshift4/driver-toolkit-rhel9:v4.13.0-202308011445.p0.gd719bdc.assembly.stream as builder

FROM ${DTK} as builder
WORKDIR /usr/src
RUN ["git", "clone", "https://github.com/chr15p/partner_templates.git"]
WORKDIR /usr/src/partner_templates/kernel_module/
RUN KERNEL_SRC_DIR=/lib/modules/${KERNEL_VERSION}/build make all KVER=${KERNEL_VERSION}

FROM registry.redhat.io/ubi9/ubi-minimal
ARG KERNEL_VERSION
RUN ["microdnf", "install", "-y", "kmod"]
COPY --from=builder /usr/src/partner_templates/kernel_module/*.ko /opt/lib/modules/${KERNEL_VERSION}/
RUN depmod -b /opt ${KERNEL_VERSION}

CMD [ "modprobe", "-d", "/opt", "ptemplate_char_dev"]
```

Then as we did previously we build the driver container with:

```
podman build -f Dockerfile.hard --build-arg KERNEL_VERSION=$(uname -r) \
		-t quay.io/example/pt-char-dev:$(uname -r) .
```


## Discussion

This Docker file works in a very similar way to the previous version, downloading the source code from github, and building it with `make`, then in the second stage copying over only the .ko file created to a clean minimal image. The big difference is that the first stage does not need to install the build tools first, this removes any issues with entitlements and the availability of packages.




# Links

* Driver Toolkit https://github.com/openshift/driver-toolkit
* Dockerfile reference https://docs.docker.com/engine/reference/builder/
* podman https://podman.io/ 
* podman build https://docs.podman.io/en/latest/markdown/podman-build.1.html





