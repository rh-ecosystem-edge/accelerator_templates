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

In which case you can run simply run it:

```
podman run quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```



```
podman push quay.io/example/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```


# Links

* podman https://podman.io/ 
* podman build https://docs.podman.io/en/latest/markdown/podman-build.1.html
* podman push https://docs.podman.io/en/latest/markdown/podman-push.1.html
* podman run https://docs.podman.io/en/latest/markdown/podman-run.1.html





