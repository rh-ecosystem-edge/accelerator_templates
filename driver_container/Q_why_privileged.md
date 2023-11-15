# Question: Why does my driver container need to be run as privileged

## Answer

A driver container needs to be run as `--privileged`, by default a container is only allowed limited access to devices, which prevents loading kernel modules. A "privileged" container is given the same access to devices as the user launching the container. Running a container with `--privileged` turns off the security features that isolate the container from the host. Dropped Capabilities, limited devices, read-only mount points, AppArmor/SELinux separation, and seccomp filters are all disabled. Clearly this is not a thing you want to do for all containers, but in the case of driver containers is necessary.

Failing to add the `--privileged` flag results in an error similar to:

```bash
modprobe: ERROR: could not insert 'ptemplate_char_dev': Operation not permitted
```

## Links

* [What the --privileged flag does with container engines such as Podman, Docker, and Buildah](https://www.redhat.com/sysadmin/privileged-flag-container-engines)
