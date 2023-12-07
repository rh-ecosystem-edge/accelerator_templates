# Metrics

1. [Introduction](#introduction)

2. [Grafana Dashboards](#grafana-dashboards)

3. [PrometheusRule Objects](#prometheusrule-objects)

4. [Links](#links)

## Introduction

Logging metrics is important for cluster admins to understand if their cluster and the devices attached to it are in good working order and to do capacity planning.  They can also be used in other places such as defining custom auto-scaling rules.

An operator can create its own [custom metrics](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#monitoring-and-observability) using the   [Prometheus golang client library](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus). For an example of how this can be done see the [Ptemplate example operator implementation](ptemplate-observability.md).

## Grafana Dashboards

For an Operator to integrate seamlessly with the OpenShift admin dashboard it is highly recommended to provide a Grafana dashboard spec in the operator. This allows cluster admins to have an out-of-the-box view of the devices and have a single pane of glass to examine in the even of problems.

For details on how to create a dashboard see the [Grafana Documentation](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/create-dashboard/)

## PrometheusRule Objects

Once metrics are being passed to to Prometheus they can be used to trigger alerts that can then be caught by Alert Manager and wake up someone at 2 AM with the integration to PagerDuty or in the best case trigger automated remediation.

This requires the creation of `PrometheusRule` objects with expressions (`spec.groups.expr`) that reference the metric name.

For example, this rule will trigger an `InstanceDown` alert when the `Ptemplate_consumers` metric is zero for more then 1 minute:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: example-alert
  namespace: ns1
spec:
  groups:
  - alert: InstanceDown
    expr: Ptemplate_consumers == 0
    for: 1m
```

## Links

[Configuring The Monitoring Stack In OpenShift](https://access.redhat.com/documentation/en-us/openshift_container_platform/4.14/html/monitoring/index)

[OpenShift PrometheusRule reference](https://docs.openshift.com/container-platform/4.14/rest_api/monitoring_apis/prometheusrule-monitoring-coreos-com-v1.html)

[Prometheus Alerting Rules Syntax](https://www.prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)

[OpenShift Monitoring overview](https://access.redhat.com/documentation/en-us/openshift_container_platform/4.14/html/monitoring/monitoring-overview)

[Creating Alerting Rules For User Defined Projects](https://access.redhat.com/documentation/en-us/openshift_container_platform/4.14/html/monitoring/managing-alerts#creating-alerting-rules-for-user-defined-projects_managing-alerts)
