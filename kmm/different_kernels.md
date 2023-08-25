# I need to deploy my driver container on multiple kernels Openshift.

## Problem

Different versions of Openshift come with different kernels, and each new kernel version potentially requires changes in kernel modules, so to build a generic Module resource you need to be able to respond to these different versions.

## Solution

The `kernelMappings` section of the `Module` resource allows multiple mappings that allows mapping between a driver container image name and a regular expression that matches against the running kernel. For example:


```
      kernelMappings:
        - literal: 6.4.11-200.fc38.x86_64
          containerImage: "quay.io/chrisp262/pt-char-dev-f37:6.4.11-200.fc38.x86_64"

        - regexp: '^.+\fc37\.x86_64$'
          containerImage: "quay.io/chrisp262/pt-char-dev-f37:${KERNEL_FULL_VERSION}"

        - regexp: '^.+$'
          containerImage: "quay.io/chrisp262/pt-char-dev:${KERNEL_FULL_VERSION}"

```


