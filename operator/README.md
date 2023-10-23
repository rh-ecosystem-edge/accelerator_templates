# Operator

## Introduction

An Operator is an extension to Kubernetes and Openshift that adds a Custom Resource type that provides higher level functionality.


### The Ptemplate-operator example

Our example [ptemplate-operator](ptemplate-operator) implements the `Ptemplate` resource type (or Kind). When a resource of this Kind is created Kubernetes notifies the ptemplate-operator instance which runs its Reconcile() function which then creates a number of other resources needed to run our example. It does this by creating a Module resource that `KMM` picks up and processes, and a `Daemonset` resource for or consumer pods that is picked up by core Kubernetes 

```
apiVersion: ptemplates.pt.example.com/v1alpha
kind: Ptemplate
metadata:
  name: consumer-pod
spec:
  maxdev: 5
  defaultmsg: ptemplate
  consumerimage: quay.io/chrisp262/pt-device-plugin:consumer-latest
  requiredDevices: 2
  imageRepoSecret:
    name: pt-char-pull-secret
  selector:
    node-role.kubernetes.io/worker: worker
```

Implementing an operator requires a large amount of boilerplate code.
The easiest way to create this code is to use the operator-sdk. It's [Quickstart](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/) and [tutorial](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/) walks you through generating the boilerplate code for the operator and adding an api to it (which defines the yaml structure used to create the resource instance).

Once its generated all the code, a minimum the following files need customising to implement the functionality that this Operator requires.

#### api/v1alpha1/ptemplate_types.go
This contains the data types that define the Custom Resources specification as a series of nested golang structs. Each field has a `json:`  tag that is used to define the associated field in the yaml (or json) used to create the resource
e.g. 
```
ConsumerImage   string                   `json:"consumer"`
```
means that if a `consumer` field is defined in the yaml that value is used to populate the ConsumerImage field in the relevant golang struct


#### controllers/ptemplate_controller.go
The heart of the operator is the Reconcile() function in this file which implements the actual functionality. 

This function is run whenever an object of a type owned by the operator (listed in the `SetupWithManager()` function), so it needs to be indempotent. 
In the ptemplate-operator we own a daemonset so reconcile is run for each ptemplate resource, whenever any daemonset in the cluster is changed. It needs to get the config for the relevant ptemplate resource, check the sub-resources that should exist (in this case a KMM Module, and a Daemonset) and reconcile them, creating them if they dont exist, correcting them if they differ from the `ptemplate.Spec` , update the `ptemplate.Status` if needed, and most importantly of all, do nothing if nothing has changed.

For production purposes it would make sense to break [this file](../ptemplate-operator/controllers/ptemplate_controller.go) into separate files for the Module and consumer Daemonset, remove some of the logging, handle errors better, and generally do all the things that turn example code into production quality. But that would also reduce its clarity!


#### config/rbac/role.yaml
By default the operator only has permissions on its own CRD, to be able to access any other resource in the cluster (e.g. to list, and create Daemonsets or Pods) permissions need to be added to this file. These then get deployed as part of the operator-manager-role e.g. 
```
 kubectl get clusterrole ptemplate-operator-manager-role 
```

If the operator logs show permissions errors this is probably the file that needs changing.


#### Makefile
The variable `IMAGE_TAG_BASE` was set to have the required image name and registry, for convenience. 


### Building

Building the operator is as easy as running ` make build` this will also regenerate any of the boilerplate code to reflect changes in the api (`api/v1alpha1/ptemplate_types.go`). 


Once the code is building successfully you can build the operator image running the `docker-build` target

```
make docker-build
```
Will build the code and create a docker image with it named for the  `IMAGE_TAG_BASE` variable, which can then be pushed to a registry.


```
make deploy
```
Will deploy the image as a pod. It will also create all the other required resource such as the `CRD` itself, `serviceaccounts`, `clusterroles` and `rolebindings` etc. These are initially created by operator-sdk, but can be changed if needed.


(During testing I had a strange error where I was writing to the `ptemplate.status.consumers` field to report the status but it wasn't getting changed, and no errors were being thrown, eventually I realised I had updated the CRD, and the `config/crd/bases/pt.example.com_ptemplates.yaml` file had been updated by Make, but I hadn't re-run `make deploy` so the kubernetes wasn't aware of my changes. Once I'd done that everything worked correctly)



## Cookbook

## Links
* [Operator SDK](https://sdk.operatorframework.io/docs/building-operators/golang/quickstart/) 

* [Kuberenets API](https://kubernetes.io/docs/reference/kubernetes-api/) 

* [Kubebuilder](https://book.kubebuilder.io/introduction) is not technically the same as the Operator SDK but is close enough to provide a good alternate view when the Operator SDK docs fail you