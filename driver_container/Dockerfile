ARG DTK=registry.redhat.io/openshift4/driver-toolkit-rhel9:v4.13.0-202308011445.p0.gd719bdc.assembly.stream as builder
FROM ${DTK} as builder
ARG KERNEL_VERSION
WORKDIR /usr/src
RUN ["git", "clone", "https://github.com/chr15p/partner_templates.git"]
WORKDIR /usr/src/partner_templates/kernel_module/
RUN KERNEL_SRC_DIR=/lib/modules/${KERNEL_VERSION}/build make all KVER=${KERNEL_VERSION}

FROM registry.redhat.io/ubi9/ubi-minimal
ARG KERNEL_VERSION
RUN ["microdnf", "install", "-y", "kmod"]
COPY --from=builder /usr/src/partner_templates/kernel_module/ptemplate_char_dev.ko /opt/lib/modules/${KERNEL_VERSION}/
RUN depmod -b /opt ${KERNEL_VERSION}
