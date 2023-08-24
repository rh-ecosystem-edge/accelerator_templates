# I need an easier way to build a driver container

## Problem

I have a driver container with my kernel module, now I need to run it on my machine.

## Solution

A driver container can simply be run using the `podman` command:

```
podman run --privileged quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```

If all goes well this will load the kernel module and running `lsmod` on the host operating system will list our new kernel module as loaded.

The container doing the loading should finish immediately after completing the load, so to unload the kmod again you will need to use `rmmod` either from the host or in its own containerised command such as:

```
podman run --privileged quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64 rmmod ptemplate_char_dev
```



## Discussion

A driver container image should be built to run `modprobe -d /opt <driver_name>` by default, in which case if we simply run it then it loads the driver and exists leaving the kmod in the kernel.

You can check the default startup command simply, with 

```
# podman inspect -f "{{.Config.Cmd}}" <image_name>
```

Which should give an output something like:

```
[modprobe -d /opt <driver_name>]
```

E.g.

```
# podman inspect -f "{{.Config.Cmd}}" quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64

[modprobe -d /opt ptemplate_char_dev]

```

If not this can simply be overridden on the `podman` command line such as:

```
podman run --privileged quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64 modinfo -d /opt ptemplate_char_dev
``


A driver container needs to be run as `--privileged`, by default a container is only allowed limited access to devices, which prevents loading kernel modules. A "privileged" container is given the same access to devices as the user launching the container. Running a container with `--privileged` turns off the security features that isolate the container from the host. Dropped Capabilities, limited devices, read-only mount points, Apparmor/SELinux separation, and Seccomp filters are all disabled. Clearly this is not a thing you want to do for all containers, but in the case of driver containers is necessary.

Failing to add the `--privileged` flag results in an error similar to:

```
modprobe: ERROR: could not insert 'ptemplate_char_dev': Operation not permitted
```


# Links

* podman https://podman.io/ 
* podman build https://docs.podman.io/en/latest/markdown/podman-build.1.html
* podman push https://docs.podman.io/en/latest/markdown/podman-push.1.html
* podman run https://docs.podman.io/en/latest/markdown/podman-run.1.html





