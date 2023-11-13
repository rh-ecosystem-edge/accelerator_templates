# I want to use NFD to add feature label nodes

## Solution

The default NFD configuration adds labels for a wide range of useful features.

[See here](cpu_flags.yaml)

```
apiVersion: nfd.kubernetes.io/v1
kind: NodeFeatureDiscovery
metadata:
  name: my-nfd-deployment
  namespace: nfd
spec:
  operand:
    image: registry.k8s.io/nfd/node-feature-discovery:v0.13.4
    imagePullPolicy: IfNotPresent
  workerConfig:
    configData:
```

Will add a whole range of default labels:

```
feature.node.kubernetes.io/cpu-cpuid.ADX=true
feature.node.kubernetes.io/cpu-cpuid.AESNI=true
feature.node.kubernetes.io/cpu-cpuid.AVX2=true
feature.node.kubernetes.io/cpu-cpuid.AVX=true
feature.node.kubernetes.io/cpu-cpuid.CMPXCHG8=true
feature.node.kubernetes.io/cpu-cpuid.FMA3=true
feature.node.kubernetes.io/cpu-cpuid.FXSR=true
feature.node.kubernetes.io/cpu-cpuid.FXSROPT=true
feature.node.kubernetes.io/cpu-cpuid.HYPERVISOR=true
feature.node.kubernetes.io/cpu-cpuid.IA32_ARCH_CAP=true
feature.node.kubernetes.io/cpu-cpuid.IBPB=true
feature.node.kubernetes.io/cpu-cpuid.IBRS=true
feature.node.kubernetes.io/cpu-cpuid.LAHF=true
feature.node.kubernetes.io/cpu-cpuid.MD_CLEAR=true
feature.node.kubernetes.io/cpu-cpuid.MOVBE=true
feature.node.kubernetes.io/cpu-cpuid.MPX=true
feature.node.kubernetes.io/cpu-cpuid.OSXSAVE=true
feature.node.kubernetes.io/cpu-cpuid.SPEC_CTRL_SSBD=true
feature.node.kubernetes.io/cpu-cpuid.STIBP=true
feature.node.kubernetes.io/cpu-cpuid.SYSCALL=true
feature.node.kubernetes.io/cpu-cpuid.SYSEE=true
feature.node.kubernetes.io/cpu-cpuid.VMX=true
feature.node.kubernetes.io/cpu-cpuid.X87=true
feature.node.kubernetes.io/cpu-cpuid.XGETBV1=true
feature.node.kubernetes.io/cpu-cpuid.XSAVE=true
feature.node.kubernetes.io/cpu-cpuid.XSAVEC=true
feature.node.kubernetes.io/cpu-cpuid.XSAVEOPT=true
feature.node.kubernetes.io/cpu-cpuid.XSAVES=true
feature.node.kubernetes.io/cpu-hardware_multithreading=false
feature.node.kubernetes.io/cpu-model.family=6
feature.node.kubernetes.io/cpu-model.id=142
feature.node.kubernetes.io/cpu-model.vendor_id=Intel
feature.node.kubernetes.io/kernel-config.NO_HZ=true
feature.node.kubernetes.io/kernel-config.NO_HZ_FULL=true
feature.node.kubernetes.io/kernel-version.full=5.14.0-284.25.1.el9_2.x86_64
feature.node.kubernetes.io/kernel-version.major=5
feature.node.kubernetes.io/kernel-version.minor=14
feature.node.kubernetes.io/kernel-version.revision=0
feature.node.kubernetes.io/system-os_release.ID=rhel
feature.node.kubernetes.io/system-os_release.VERSION_ID.major=9
feature.node.kubernetes.io/system-os_release.VERSION_ID.minor=2
feature.node.kubernetes.io/system-os_release.VERSION_ID=9.2
```

Which are then available for use in selectors including by `kmm`


## Discussion

NFD can produce a wide range of labels based on the nodes hardware and OS, even the few that are available by default provide some highly useful information about the CPUs and OS that are running on the node.


## Links

* [NFD sources.CPU reference](https://kubernetes-sigs.github.io/node-feature-discovery/v0.13/reference/worker-configuration-reference.html)

