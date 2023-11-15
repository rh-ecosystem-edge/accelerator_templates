# I need to build a driver container with DTK

## Problem

You have your driver source code ready to compile, you want to use the Driver-Toolkit image, but can't use OpenShift BuildConfig resources for whatever reason,

## Solution

While using a BuildConfig is the recommended solution for creating DTK based driver containers, it possible to use the DTK image within a standard Dockerfile, and build an image with `podman`

For example for the ptemplate_char_dev driver container a standalone [Dockerfile](./Dockerfile.easy) might look like:

```Dockerfile
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

Then we build the driver container with podman via:

```bash
podman build -f Dockerfile.hard --build-arg KERNEL_VERSION=$(uname -r) \
             -t quay.io/example/pt-char-dev:$(uname -r) .
```

## Discussion

This Docker file works in a very similar way to the [non-DTK version](driver_container_hard_way.md), downloading the source code from github, and building it with `make`, then in the second stage copying over only the .ko file created to a clean minimal image. The big difference is that the first stage does not need to install the build tools first, this removes any issues with entitlements and the availability of packages.

## Links

* [Driver Toolkit](https://github.com/openshift/driver-toolkit)

* [Dockerfile reference](https://docs.docker.com/engine/reference/builder/)

* [podman](https://podman.io/)

* [podman build](https://docs.podman.io/en/latest/markdown/podman-build.1.html)
