# I need to deploy my driver container on multiple kernels Openshift.

## Problem

Different versions of Openshift come with different kernels, and each new kernel version potentially requires changes in kernel modules, so to build a generic Module resource you need to be able to respond to these different versions.

## Solution

The `kernelMappings` section of the `Module` resource allows multiple mappings that allows mapping between a driver container image name and a regular expression that matches against the running kernel. For example:

([see full file](different_kernels.yaml))

```
      kernelMappings:
        - literal: 6.4.11-200.fc38.x86_64
          containerImage: "quay.io/chrisp262/pt-char-dev-f38:6.4.11-200.fc38.x86_64"

        - regexp: '^.+\fc37\.x86_64$'
          containerImage: "quay.io/chrisp262/pt-char-dev-f37:${KERNEL_FULL_VERSION}"

        - regexp: '^.+$'
          containerImage: "quay.io/chrisp262/pt-char-dev:${KERNEL_FULL_VERSION}"

```

This set of mappings will load a specific driver container for any node running the 6.4.11-200.fc38.x86_64 kernel, a Fedora 37 specific driver container image for nodes running a Fedora 37 kernel, and a default version for any other nodes in the cluster, which would include any of the released CoreOS kernels, and any Fedora 38 kernel that doesn't match the first (literal) mapping.


## Discussion

The regexp setting uses the golang regular expression syntax so can support complex regular expressions, and the mappings are tested in the order they are given from the top down. Together this makes it possible to support a wide range of different kernels from the same `module` resource, including not just different releases of Openshift, or different versions of the underlying OS, but different versions of Linux entirely. 

## Links

* [Golang regular expression syntax](https://pkg.go.dev/regexp/syntax)