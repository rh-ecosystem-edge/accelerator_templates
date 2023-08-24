# Example Kernel Module

This is an toy kernel module created for OpenShift Partner Templates. It has no guarantees that it will work and may well cause random kernel panics, security flaws, and other chaos. Do not use in production under any circumstances!

It does, however, illustrate some topics in creating a robust solution for developing, distributing, and deploying third party, out-of-tree kernel modules on Openshift.


When loaded it creates a number of character devices named /dev/ptemplate-0 to /dev/ptemplate-X where X is the value of of the `max_dev` parameter (or 5 by default). 

When read from it gives a short message, and when written to it replaces that message. Max message length is 100 characters. The initial message stored is governed by the `default_msg`

```
[root@kube93 nfd]# cat /dev/ptemplate-0
ptemplate
[root@kube93 nfd]# echo "hello world" >  /dev/ptemplate-0
[root@kube93 nfd]# cat /dev/ptemplate-0
hello world
```

## Building
To build for the current kernel run `make`:

```
> # cd kernel_module/
> [kernel_module]# make
> make -C /usr/src/kernels/6.4.7-200.fc38.x86_64/ M=/home/cprocter/engineering/partner_templates/kernel_module EXTRA_CFLAGS=-DKMODVER=\\\"2c7beaf\\\" modules
> make[1]: Entering directory '/usr/src/kernels/6.4.7-200.fc38.x86_64'
> warning: the compiler differs from the one used to build the kernel
>   The kernel was built by: gcc (GCC) 13.1.1 20230614 (Red Hat 13.1.1-4)
>   You are using:           gcc (GCC) 13.2.1 20230728 (Red Hat 13.2.1-1)
>   CC [M]  /home/example/partner_templates/kernel_module/ptemplate_char_dev.o
>   MODPOST /home/example/partner_templates/kernel_module/Module.symvers
>   CC [M]  /home/exampleg/partner_templates/kernel_module/ptemplate_char_dev.mod.o
>   LD [M]  /home/example/partner_templates/kernel_module/ptemplate_char_dev.ko
>   BTF [M] /home/example/partner_templates/kernel_module/ptemplate_char_dev.ko
> Skipping BTF generation for /home/example/partner_templates/kernel_module/ptemplate_char_dev.ko due to unavailability of vmlinux
> make[1]: Leaving directory '/usr/src/kernels/6.4.7-200.fc38.x86_64'
```

## Loading and unloading

To load from the current directory:

```
insmod ./ptemplate_char_dev.ko default_msg="hello world"
```

or you can install it to the modules directory for your kernel then use modprobe (which is generally the prefered tool for loading kmods) e.g.

```
cp ./ptemplate_char_dev.ko  /lib/modules/`uname -r`/
modprobe ptemplate_char_dev default_msg="hello world"
```

And to unload:

```
rmmod ptemplate_char_dev
```

## Further Reading:
There are many guides only to writing and building kernel modules. 

* A good place to start is the RHEL documentation on kernel modules https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/managing_monitoring_and_updating_the_kernel/managing-kernel-modules_managing-monitoring-and-updating-the-kernel
* Linux kernel documentation
https://docs.kernel.org/kbuild/modules.html 




