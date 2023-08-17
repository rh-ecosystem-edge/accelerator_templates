# Driver Container

To use a kernel module on Openshift it needs to be built into a container image which Openshift can then use to instantiate a pod. These images are commonly known as Driver Containers. 

Building a kernel module inside a Dockerfile is made more complicated by the large number of additional packages that are required for the build process (e.g. a C compiler). Getting these right can be an annoying process of trial and error, and leaving them in a driver-container that is deployed in production both bloats the container image (often by hundreds of MB), and leaves an additional attack surface for no value.

To remedy this Red Hat provides the [driver toolkit](https://github.com/openshift/driver-toolkit), an image with the required tools for building kernel modules. This can then be used in a multi-stage dockerfile to build a container image with just the minimal files to hold and load the kernel module.

# Example

The following Dockerfile uses the Driver-Toolkit image. It downloads the source code from github and builds it (via `make all`) for the kernel version given in the KERNEL_VERSION argument.

It then downloads the ubi-minimal image, installs the kmod package in it to provide the `modprobe` and `rmmod` tools, and copies the build kernel module (.ko file) from the builder image (our driver toolkit image) into the /opt/lib/modules/${KERNEL_VERSION}/

Finally it sets the default command to execute when the container loads to load the build kernel module.
 
```
ARG KERNEL_VERSION
ARG DTK=registry.redhat.io/openshift4/driver-toolkit-rhel9:v4.13.0-202308011445.p0.gd719bdc.assembly.stream as builder
ARG KMOD_SRC=https://github.com/chr15p/partner_templates.git
ARG REPO_DIR=partner_templates/kernel_module/

FROM ${DTK} as builder
WORKDIR /usr/src
RUN ["git", "clone", "${KMOD_SRC}"]
WORKDIR /usr/src/${REPO_DIR}
RUN KERNEL_SRC_DIR=/lib/modules/${KERNEL_VERSION}/build make all KVER=${KERNEL_VERSION}

FROM registry.redhat.io/ubi9/ubi-minimal
ARG KERNEL_VERSION
RUN ["microdnf", "install", "-y", "kmod"]
COPY --from=builder /usr/src/${REPO_DIR}/*.ko /opt/lib/modules/${KERNEL_VERSION}/
RUN depmod -b /opt ${KERNEL_VERSION}

CMD [ "modprobe", "-d", "/opt", "ptemplate_char_dev"]
```


To use it run (modify for your running kernel):

```
podman build -f Dockerfile --build-arg KERNEL_VERSION=5.14.0-284.25.1.el9_2.x86_64 -t quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
podman push quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64

podman run quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```

# Links

* Driver Toolkit https://github.com/openshift/driver-toolkit
* Dockerfile reference https://docs.docker.com/engine/reference/builder/
* podman https://podman.io/ 
* podman build https://docs.podman.io/en/latest/markdown/podman-build.1.html
* podman push https://docs.podman.io/en/latest/markdown/podman-push.1.html
* podman run https://docs.podman.io/en/latest/markdown/podman-run.1.html


