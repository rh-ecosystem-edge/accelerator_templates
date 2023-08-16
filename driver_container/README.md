# Driver Container

This will build a driver container from the source in the `https://github.com/chr15p/partner_templates.git` repo under the `device_driver` directory

To use it run:
```
podman build -f Dockerfile --build-arg KERNEL_VERSION=5.14.0-284.25.1.el9_2.x86_64 -t quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
podman push quay.io/chrisp262/pt-char-dev:5.14.0-284.25.1.el9_2.x86_64
```
