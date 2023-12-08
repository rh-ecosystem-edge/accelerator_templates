ARG DTK
FROM ${DTK} as builder
ARG KVER
RUN ls /lib/modules/
WORKDIR /build/
RUN git clone https://github.com/chr15p/partner_templates.git
WORKDIR /build/partner_templates/src/kernel_module
RUN make all install KVER=${KVER}
FROM registry.redhat.io/ubi9/ubi-minimal
ARG KVER
# Required for installing `modprobe`
RUN microdnf install -y kmod
COPY --from=builder /lib/modules/${KVER}/ptemplate_char_dev.ko /lib/modules/${KVER}/
COPY --from=builder /lib/modules/${KVER}/ptemplate_char_dev.ko /lib/modules/${KVER}/
RUN depmod ${KVER}
