# I want to manage which nodes my driver is loaded on

## Problem

You only want to load a driver on nodes that have the particular piece of hardware that it enables.

## Solution 

Kubernetes offers a robust labelling and selector system that KMM can use to select the nodes to load your kmod on.

Firstly label the required nodes with the `kubectl label nodes` command
For example:

```
kubectl label node kube93.cp.chrisprocter.co.uk  ptemplate=required
```

Once the nodes have the required labels update the selector field of the Module resource yaml to include it:

([see full file](node_selectors.yaml)):

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
    ptemplate: "required"
```

This should now only install our ptemplate_char_dev kernel module on nodes that have the       node-role.kubernetes.io/worker=worker label AND the ptemplate=required label.


## Discussion

You can check what labesl a node has with :

```
kubectl get nodes --show-labels
```

Although because they print as a comma separated list once you get to more than a handful of labels I generally find it easier to read them by replacing the commas with newlines. e.g

```
kubectl get nodes kube93.cp.chrisprocter.co.uk --show-labels | tr "," "\n" 
```

Or if you know what labels should exist its easy to search for just the nodes with that label e.g.

```
[root@kube92 ~]# kubectl get nodes -l ptemplate=required
NAME                           STATUS   ROLES    AGE    VERSION
kube93.cp.chrisprocter.co.uk   Ready    worker   186d   v1.27.4
```

While it might seem that labels are an ideal place to store data they do come with some restrictions, valid label values:

*   must be 63 characters or less (can be empty),
*    unless empty, must begin and end with an alphanumeric character ([a-z0-9A-Z]),
*   could contain dashes (-), underscores (_), dots (.), and alphanumerics between.



Several useful labels are added to nodes by default including `kubernetes.io/hostname` and `kubernetes.io/arch`, and then kmm itself adds several more, such as the kernel version `
kmm.node.kubernetes.io/kernel-version.full` and a label for each of the kernel modules it has loaded `kmm.node.kubernetes.io/pt-char-dev.ready`.  These are very useful for things like only loading kmods on the architecture they are compiled for or for only loading them when other pre-requisite modules already exist.

Other operators will add their own labels, or especial note is the the Node Feature Discovery operator whose entire purpose is to automatically label nodes based on hardware features such as pci or usb hardware available, and cpu features



## Links

https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/#add-a-label-to-a-node

https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
