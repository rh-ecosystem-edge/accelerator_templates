package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	k8sdeviceplugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type exampleDevicePlugin struct {
	k8sdeviceplugin.UnimplementedDevicePluginServer
}

var logger logr.Logger

const (
	DEVICENAMEPREFIX = "ptemplate"
	RESOURCENAME     = "example.com/" + DEVICENAMEPREFIX
)

// GetDevicePluginOptions always returns an empty response.
func (p *exampleDevicePlugin) GetDevicePluginOptions(_ context.Context, _ *k8sdeviceplugin.Empty) (*k8sdeviceplugin.DevicePluginOptions, error) {
	return &k8sdeviceplugin.DevicePluginOptions{}, nil
}

func (p *exampleDevicePlugin) ListAndWatch(empty *k8sdeviceplugin.Empty, stream k8sdeviceplugin.DevicePlugin_ListAndWatchServer) error {
	// Send an initial list of available devices.

	devPathGlob := fmt.Sprintf("/dev/%s-*", DEVICENAMEPREFIX)

	ptdevices, err := filepath.Glob(devPathGlob)
	if err != nil {
		klog.Error(err, "could not get devices", "error", err)
		os.Exit(1)
	}
	logger.Info("devices found", "matching", devPathGlob, "ptdevices", ptdevices)
	deviceCount := len(ptdevices)
	resp := &k8sdeviceplugin.ListAndWatchResponse{
		Devices: make([]*k8sdeviceplugin.Device, deviceCount),
	}

	for i, dev := range ptdevices {
		devminor := strings.SplitN(dev, "-", 2)[1]

		resp.Devices[i] = &k8sdeviceplugin.Device{
			ID:     devminor,
			Health: k8sdeviceplugin.Healthy,
		}
	}

	if err := stream.Send(resp); err != nil {
		return status.Errorf(codes.Unknown, "failed to send response: %v", err)
	}

	logger.Info("successfully reported node device to the kubelet", "number devices", deviceCount)

	// Wait for the stream to be closed or cancelled.
	<-stream.Context().Done()

	return nil
}

func (p *exampleDevicePlugin) Allocate(ctx context.Context, req *k8sdeviceplugin.AllocateRequest) (*k8sdeviceplugin.AllocateResponse, error) {
	// Check that the requested devices are available.
	for _, containerReq := range req.ContainerRequests {
		for _, id := range containerReq.DevicesIDs {
			_, err := os.Stat("/dev/" + DEVICENAMEPREFIX + "-" + id)
			if err != nil {
				logger.Error(fmt.Errorf("requested device not present on the node"), "missing device", " device id ", id)
				return nil, status.Errorf(codes.NotFound, "requested device %s is not available", id)
			}
		}
	}

	// Return the allocated devices.
	resp := &k8sdeviceplugin.AllocateResponse{
		ContainerResponses: []*k8sdeviceplugin.ContainerAllocateResponse{},
	}
	for _, req := range req.ContainerRequests {
		containerResp := &k8sdeviceplugin.ContainerAllocateResponse{
			Envs:        make(map[string]string, len(req.DevicesIDs)),
			Annotations: make(map[string]string, len(req.DevicesIDs)),
		}

		for _, id := range req.DevicesIDs {
			devPath := fmt.Sprintf("/dev/%s-%s", DEVICENAMEPREFIX, id)
			envVariable := fmt.Sprintf("%s_%s", DEVICENAMEPREFIX, id)

			containerResp.Envs[envVariable] = devPath
			logger.Info("set environment variable", "env name", envVariable, "value", containerResp.Envs[envVariable])

			annotationKey := devPath
			containerResp.Annotations[annotationKey] = ""
			logger.Info("set annotation", "annotation name", annotationKey, "value", containerResp.Annotations[annotationKey])

			containerResp.Devices = append(containerResp.Devices, &k8sdeviceplugin.DeviceSpec{
				HostPath:      devPath,
				ContainerPath: devPath,
				Permissions:   "rw",
			})
			logger.Info("setting device file", "hostPath", devPath, "containerPath", devPath)
		}
		resp.ContainerResponses = append(resp.ContainerResponses, containerResp)
	}
	return resp, nil
}

func main() {
	logger = klogr.New().WithName("device-plugin")

	var configFile string
	flag.StringVar(&configFile, "config", "", "The path to the configuration file.")
	flag.Parse()

	socketName := fmt.Sprintf("%s.sock", DEVICENAMEPREFIX)
	pluginSocketPath := fmt.Sprintf("/var/lib/kubelet/device-plugins/%s", socketName)

	// Create a listener for the gRPC server.
	listener, err := net.Listen("unix", pluginSocketPath)
	if err != nil {
		os.Remove(pluginSocketPath)
		logger.Error(err, "failed to listen on the socket", err)
		os.Exit(1)
	}

	// Create a new gRPC server and register our device plugin with it.
	server := grpc.NewServer()
	k8sdeviceplugin.RegisterDevicePluginServer(server, &exampleDevicePlugin{})

	logger.Info("grpc server created, callbacks registered, listening for commands", "unix socket", socketName)

	var g run.Group

	g.Add(
		func() error {
			if err := server.Serve(listener); err != nil {
				return fmt.Errorf("gRPC server exited unexpectedly: %v", err)
			}
			return nil
		},
		func(error) {
			server.Stop()
		},
	)

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(
		func() error {
			kubeletSock := "/var/lib/kubelet/device-plugins/kubelet.sock"
			conn, err := grpc.Dial(kubeletSock, grpc.WithInsecure(), grpc.WithBlock(),
				grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
					return net.DialTimeout("unix", addr, timeout)
				}))
			if err != nil {
				return fmt.Errorf("failed to dial grpc: %v", err)
			}

			client := k8sdeviceplugin.NewRegistrationClient(conn)
			request := &k8sdeviceplugin.RegisterRequest{
				Version:      k8sdeviceplugin.Version,
				Endpoint:     socketName,
				ResourceName: RESOURCENAME,
			}
			if _, err = client.Register(context.Background(), request); err != nil {
				return fmt.Errorf("failed to register to kubelet: %v", err)
			}

			logger.Info("plugin registered with kubelet", "ResourceName", RESOURCENAME)

			conn.Close()
			<-ctx.Done()
			return nil
		},
		func(error) {
			os.Remove(pluginSocketPath)
			cancel()
		},
	)

	// Exit gracefully on SIGINT and SIGTERM.
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	g.Add(
		func() error {
			for {
				select {
				case <-term:
					return nil
				case <-ctx.Done():
					return nil

				}
			}
		},
		func(error) {
			cancel()
		},
	)

	// Wait for the server to exit.
	g.Run()
}
