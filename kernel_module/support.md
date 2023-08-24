# I want to know who supports my kernel module

## Problem

You have created a kernel module, but now it has an issue and you want to get help fixing it.

## Solution

In general the creator or distributor of the kernel module is responsible for supporting it (although they may choose not to). The provider of the operating system, the compiled kernel, or the kernel source code do not support it.

Therefore if you have created a kernel module, your operating system vendor can not help you debug it or add features.


## Discussion

Kernel modules can be very complicated and add a very wide range of functionality to the kernel, generally the only person or organisation that understands them is the creator of the code, any third party is going to struggle to help even if they do have a deep knowledge of kernel internals or other kernel modules.

