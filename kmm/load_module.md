# I want to deploy a pre-built driver container with KMM

## Problem
You have a pre-built Driver Container and want to load it on all the nodes in a cluster that match a particular criteria


## Solution

Create a yaml file defining a `Module` resource with the `containerImage` parameter set to the name of the Driver Container that contains the kmod you want to deploy.

([see full file](load_module.yaml)):

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
          containerImage: "quay.io/chrisp262/pt-char-dev"
  imageRepoSecret:
    name: pt-char-pull-secret
  selector:
    node-role.kubernetes.io/worker: "worker"

```

Then apply this file with `oc apply -f module.yaml` or `kubectl apply -f module.yaml`

## Discussion

The above yaml defines a `Module` resource that the KMM operator will read and use to create a number of child resources including a daemonset and pods.

There are a number of fields that are worth looking at in more detail.

### The spec.moduleLoader.container section

Firstly the `spec.moduleLoader.container` section contains a number of sub-fields that define the driver container, how to create it, and how to use it to load the kernel module. It contains three sub-sections, `modprobe`, `imagePullPolicy`, and `kernelMappings`.

The `modprobe` section is used to build the modprobe command that will be run. The most important field is `moduleName` which simply sets the name of the kernel module to be loaded. Its possible for a driver container to hold multiple valid .ko files so we need to explicitly tell KMM which one this module resource needs to load.

The last field in the `modprobe` section is `parameters`. This is simply an array of configuration parameters that are passed to the kernel module at load time, in our case we are setting the `max_dev` parameter to 5 and the `default_msg` field to the string `ptemplate`.

These fields are combined by KMM to build the command line. In our example this will result in the driver container being started with the command:

```
/bin/sh -c modprobe -v -d "/opt" "ptemplate_char_dev" max_dev=5 default_msg=ptemplate
```

### The spec.moduleLoader.imagePullPolicy field

The `imagePullPolicy` setting is the simplest to understand, it simply determines when Openshift should pull the Driver Container image from the registry and when it should use a locally cached version. Generally it is advisable to set this to `Always` especially when debugging, there is nothing as frustrating as spending hours trying to hunt down a bug only to find you've been updating the registry and Openshift has been using the same cached image all the time.


### The spec.moduleLoader.kernelMappings section

The last section gives the name of the Driver Container image with the kernel module in it that Openshift should pull down and run.

As the kernel ABI can change between different kernel versions, different distributions can have different kernel version schemes, and sometimes you just want different functionality on different OS versions its possible use regular expressions matching against the running kernel version, to define which Driver Container gets used.

Our ptemplate_char_dev.ko example is quite simple and  so for the moment we will make the assumption that it will work everywhere and use a catch-all regexp to say that for all kernels use the same image.


### Other Sections
There are two other sections defined in the yaml. 

`spec.imageRepoSecret` this simply gives the \name of a previously defined Secret resource containing the credentials needed to pull the Driver Container image. In this simple case these credentials just need read permissions on the image.

`spec.selector` defines a series of label/value pairs that are matched against the node labels to determine where the kernel module is needed, in our example we only want the kmod on nodes explicitly labelled as worker nodes, but a selector defining a driver should only be loaded on nodes with a particular piece of hardware, or where some prerequisite kernel module is already loaded are also possible.



