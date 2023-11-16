# I need to support SecureBoot and signed drivers

## Solution

The `module.spec.moduleLoader.container.sign` section of kmm can be used to sign the kernel modules inside the device container using the `signfile` binary from the Linux kernel and a public/private key pair passed in as `configmap` objects.

## Discussion

On a secure-boot enabled system all kernel modules must be signed with a public/private key-pair enrolled into the Machine Owner's Key or MOK database. Any kernel module that is not signed, or is signed by a private key for which the OS does not have the corresponding public key, will be rejected with an error message.

Drivers distributed as part of a distribution should already be signed by the distributions private key, which should shown in the `signer` and `signature` fields of the `modinfo` output. The public key is distributed as part of the operating system install and so is available for checking on every machine with that OS.

Generally Linux distribution developers such as Red Hat are only willing to sign kernel modules that they distribute and are willing to on the support for. The reason for this is obviously that SecureBoot is intended to provide security against an attacker creating a malicious kernel module with full, unfettered, access to all the data on a system. Signing anything that you cannot vouch for the security of provides a vector for those attackers and essentially makes Secure Boot pointless.

An alternative is to add a suitable public key to the OS and use its corresponding private key to sign the kmod.

Adding a public key requires using the `mokutil` utility to stage the key and then rebooting the server to allow the UEFI firmware to activate the key.  Once this is done then KMM can sign the kernel modules.

## Links

* [KMM Secure Boot documentation](https://openshift-kmm.netlify.app/documentation/secure_boot/)

* [How to use Secure Boot to validate startup software](https://www.redhat.com/sysadmin/secure-boot-systemtap)
