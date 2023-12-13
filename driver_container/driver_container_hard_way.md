# I need a way to build a driver container without DTK

## Problem

You have your driver source code ready to compile, you can't use the Driver-Toolkit image for whatever reason, but need to build it into a driver container for use on RHCOS without bloating the image up with build tools.

## Solution

You have to essentially recreate the DTK base image from scratch, with all it's libraries and build packages.

Use `podman build` and a multistage Dockerfile to build the image

For example for the ptemplate_char_dev kernel module we create a Dockerfile with the following content:

```Dockerfile
FROM registry.redhat.io/ubi9/ubi as builder
ARG KERNEL_VERSION=''
RUN dnf -y install \
    kernel-devel${KERNEL_VERSION:+-}${KERNEL_VERSION} \
    kernel-headers${KERNEL_VERSION:+-}${KERNEL_VERSION} \
    kernel-modules${KERNEL_VERSION:+-}${KERNEL_VERSION} \
    elfutils-libelf-devel \
    kmod \
    binutils \
    kabi-dw \
    glibc \
    gcc \
    make \
    git


WORKDIR /usr/src
RUN ["git", "clone", "https://github.com/rh-ecosystem-edge/accelerator_templates.git"]
WORKDIR /usr/src/accelerator_templates/kernel_module
RUN KERNEL_SRC_DIR=/lib/modules/${KERNEL_VERSION}/build make all KVER=${KERNEL_VERSION}

FROM registry.redhat.io/ubi9/ubi-minimal
ARG KERNEL_VERSION
RUN ["microdnf", "install", "-y", "kmod"]
COPY --from=builder /usr/src/accelerator_templates/kernel_module/*.ko /opt/lib/modules/${KERNEL_VERSION}/
RUN depmod -b /opt ${KERNEL_VERSION}

CMD [ "modprobe", "-d", "/opt", "ptemplate_char_dev"]
```

Then to build the driver container run:

```bash
podman build -f Dockerfile.hard --build-arg KERNEL_VERSION=$(uname -r) \
             -t quay.io/example/pt-char-dev:$(uname -r) .
```

## Discussion

This Dockerfile is relatively straightforward for a multi-stage Dockerfile. The first stage pulls in all the packages needed for the build, then we use `git clone` to get the source code for the driver, and build it using its `Makefile`. The second stage takes the much smaller ubi-minimal image, installs only the `kmod` package (which includes `modprobe` and `depmod`) and copies over any .ko files from the first stage image, finally it uses `depmod` to build the metadata files that `modprobe` needs. The final line then sets the default startup command for the image to run `modprobe, so loading our kmod.

This multi-stage approach results in a much smaller image. All the build requirements result in an image of over 1GB, including compiles and kernel sources that have no value outside the build.  The second stage image contains everything we need at run time and on my system is a much skinnier 184MB.  (While I was testing this Dockerfile I kept getting strange random seeming errors that caused the first stage to fail, I eventually tracked it down to podman not being able to store all the intermediate layers it was creating because I only had a couple of GB free disk space!)

However there are significant issues with this approach. We need access to the correct repositories to install all the build components, which requires building on a host with Red Hat entitlements. We also require the kernel packages for the exact kernel version we are deploying to, so the above solution builds correctly on RHEL9 but you cannot use it for RHEL8 builds or for other Linux versions.

This means you need multiple versions of this Dockerfile, one per OS version that the driver will be released for.

## Links

* [How to use entitled image builds to build Driver Containers with UBI on OpenShift](https://cloud.redhat.com/blog/how-to-use-entitled-image-builds-to-build-drivercontainers-with-ubi-on-openshift)

* [Dockerfile reference](https://docs.docker.com/engine/reference/builder/)

* [podman](https://podman.io/)

* [podman build](https://docs.podman.io/en/latest/markdown/podman-build.1.html)
