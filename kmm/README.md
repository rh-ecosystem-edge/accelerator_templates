# Kernel Module Management (KMM) Operator
The Kernel Module Management Operator is designed too manage out-of-tree kernel modules in Kubernetes and Openshift. It does this by managing the loading and unloading of driver-containers across the cluster using node selectors to determine wich nodes require the driver deployed. Optionally it can also manage the building of driver containers and signing the kernel modules themselves for secureboot installations.





## Links

* [https://kmm.sigs.k8s.io](Upstream documentation)
* [https://github.com/kubernetes-sigs/kernel-module-management/](Upstream) (where development is done)
* [https://github.com/rh-ecosystem-edge/kernel-module-management/](Midstream) (forms the basis of the Openshift release)
