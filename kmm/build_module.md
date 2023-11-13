# I don't want to worry about maintaining many driver containers

## Problem

Prebuilt driver containers require maintaining many different images covering every possible kernel version it may be deployed on, is there a better way?

## Solution


The `kernelMappings.build` section of the Module resource takes the name of a `ConfigMap` resource that contains the Dockerfile required to create the driver container. So our example ptemplate_char_dev.ko loader would look like:

([see full file](build_module.yaml))

```
apiVersion: kmm.sigs.x-k8s.io/v1beta1
kind: Module
metadata:
  name: pt-char-dev
spec:
  moduleLoader:
    container:
      modprobe:
        moduleName: ptemplate_char_dev
        dirName: /opt
        parameters:
          - max_dev=5
          - default_msg=ptemplate
      imagePullPolicy: Always
      kernelMappings:
        - regexp: '^.+$'
          containerImage: "quay.io/chrisp262/pt-char-dev:${KERNEL_FULL_VERSION}"
          build:
            dockerfileConfigMap:
              name: pt-char-dev-dockerfile
  imageRepoSecret:
    name: pt-char-pull-secret
  selector:
    node-role.kubernetes.io/worker: "worker"

apiVersion: kmm.sigs.x-k8s.io/v1beta1
kind: ConfigMap
metadata:
  name: pt-char-dev-dockerfile
spec:
  data:
    dockerfile: |
      FROM registry.redhat.io/openshift4/driver-toolkit-rhel9:v4.13.0-202308011445.p0.gd719bdc.assembly.stream as builder
 as builder
      ARG KERNEL_VERSION
      WORKDIR /usr/src
      RUN ["git", "clone", "https://github.com/chr15p/partner_templates.git"]
      WORKDIR /usr/src/partner_templates/kernel_module/
      RUN KERNEL_SRC_DIR=/lib/modules/${KERNEL_VERSION}/build make all KVER=${KERNEL_VERSION}

      FROM registry.redhat.io/ubi9/ubi-minimal
      ARG KERNEL_VERSION
      RUN ["microdnf", "install", "-y", "kmod"]
      COPY --from=builder /usr/src/partner_templates/kernel_module/*.ko /opt/lib/modules/$ {KERNEL_VERSION}/
      RUN depmod -b /opt ${KERNEL_VERSION}

      CMD [ "modprobe", "-d", "/opt", "ptemplate_char_dev"]
```


## Discussion

If the source code for your driver is available you can build the driver container "in-cluster" using KMM's `build` parameter.  This allows the OpenShift cluster to build its own driver container customised for its own kernel version (or versions if they're not homogeneous.

This approach moves the requirements for security patching and updating the driver container to the cluster, and allows the kmod provider to concentrate on their kmod without the overhead of maintaining an array of container images.

On the other hand, it requires the source code to be available to end users such as on github. It also obviously takes the exact build tools, kernel, and even version of the source code compiled out of the hands of the driver writers which can make supporting the driver complicated, if not impossible.



## Links

* [KMM upstream documentation](https://kmm.sigs.k8s.io/documentation/module_loader_image/)

