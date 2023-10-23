#FROM registry.access.redhat.com/ubi9/go-toolset:1.19 as builder

FROM registry.access.redhat.com/ubi9/ubi-minimal
#COPY  --from=builder simple-device-plugin /simple-device-plugin
RUN mkdir /deviceplugin
COPY ptemplate-device-plugin /deviceplugin/
ENTRYPOINT ["/deviceplugin/ptemplate-device-plugin"]
