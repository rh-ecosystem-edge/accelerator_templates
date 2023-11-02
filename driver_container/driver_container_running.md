# I need an easier way to build a driver container

## Problem

I have a driver container with my kernel module, now I need to run it on my machine.

## Solution

A driver container can simply be run using the `podman` command:

```
podman run --privileged quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```

This will load the kernel module using `modprobe` and immediately exit. Running `lsmod` on the host should show the newly loaded kernel module present in the kernel.

Because the driver container exists after running `modprobe` unloading the kmod again requires running the `rmmod` command from the host, either directly or in its own containerised command e.g.:

```
podman run --privileged quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64 rmmod ptemplate_char_dev
```

To gain more control of this process, and to use with on Openshift the [Kernel Module Management](../kmm/README.md) operator automates the loading and unloading process.


## Discussion

A driver container image is built with is `CMD` setting set to `modprobe -d /opt <driver_name>` so when the container is started it runs that command and loads the driver. Once this command is completed the container exists but as the kmod is already in the kernel it remains there

You can check the default startup command for a container image simply, with 

```
# podman inspect -f "{{.Config.Cmd}}" <image_name>
```

Which should show the `modprobe` command

E.g.

```
# podman inspect -f "{{.Config.Cmd}}" quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64

[modprobe -d /opt ptemplate_char_dev]

```

It is also possible to override an incorrectly set `CMD` by passing the command when the container is run:

```
podman run --privileged quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64 modinfo -d /opt ptemplate_char_dev
```

Although this is clearly a far more complex and unwieldy method so setting the `CMD` setting correctly is far superior.


# Links

* [Podman](https://podman.io/)
* [Podman build](https://docs.podman.io/en/latest/markdown/podman-build.1.html)
* [Podman push](https://docs.podman.io/en/latest/markdown/podman-push.1.html)
* [Podman run](https://docs.podman.io/en/latest/markdown/podman-run.1.html)





