# Monitoring and Metrics

## Introduction

Kubernetes is designed to talk to Prometheus for the logging and graphing of metrics. This means there is good support for gathering and communicating values in the golang libraries.

There are two things that need to be done in the operator, firstly registering a data structure for each metric to be gathered, and secondly actually gathering the metrics. The actual sending of the metrics to Prometheus is then handled automatically.

## Metrics in the Ptemplate Operator

### Registering Metrics

In the Ptemplate operator we use a separate golang package to register the metrics we are going to be gathering. This simply creates a new object of the appropriate type for each metric, passing the name of the metric and some help text. Then it calls the `MustRegister()` function for each of them.

```golang
PtemplateGauge = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Name: "Ptemplate_consumers",
        Help: "Number of existing Ptemplate consumers",
    },
)

...

metrics.Registry.MustRegister(PtemplateGauge)
```

And we then initialise this package in the `init()` function of main.go

```golang
func init() {
...
monitoring.RegisterMetrics()
}
```

### Reporting Metrics

To then gather the metrics we use that Prometheus object anywhere else in the code we need to publish the metric

```golang
monitoring.PtemplateGauge.Set(float64(len(podNames)))
```

### Useful Metrics

Useful metrics to gather depends very much on the implementation details of the driver the operator is managing.

## Links

[Monitoring and Observability](https://sdk.operatorframework.io/docs/building-operators/golang/advanced-topics/#monitoring-and-observability)

[Prometheus golang client docs](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
