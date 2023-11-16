# I want to load custom firmware with my driver

## Solution

Use the `.spec.moduleLoader.container.modprobe.firmwarePath` parameter. This points to a directory in the driver container whose contents are copied into /var/lib/firmware on the node before the driver is loaded.

This path needs to be added to the kernels firmware search path by adding the  `firmware_class.path=/var/lib/firmware` boot parameter ([see here](https://openshift-kmm.netlify.app/documentation/firmwares/#configuring-the-lookup-path-on-nodes), and then the kernel module can load the firmware via the kernels [firmware API](https://docs.kernel.org/driver-api/firmware/introduction.html).

These files are then tidied up again when the pod is terminated and the kernel module unloaded.

## Discussion

Because the firmware version is normally closely coupled to the driver version distributing it with the driver container is generally a good practice. They should be copied into the image at build time as part of the Dockerfile.

## Links

* [KMM firmware loading docs](https://openshift-kmm.netlify.app/documentation/firmwares/)

* [Linux kernel firmware API](https://docs.kernel.org/driver-api/firmware/introduction.html)
