package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"sync"
	"time"
)

type dockerHandler struct {
	config               envConfig
	client               *client.Client
	lock                 sync.Locker
	restartContainerChan chan []restartContainer
}

const (
	dockerLabelStopTimeout     = "autoheal.stop.timeout"
	autohealAllContainersValue = "all"
	autohealManagedValue       = "true"
)

type restartContainer struct {
	containerId string
	killTimeout time.Duration
}

func newDockerHandler(config envConfig) (*dockerHandler, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithTimeout(config.CurlTimeout))
	if err != nil {
		return nil, err
	}

	_, err = dockerClient.Ping(context.Background())
	if err != nil {
		log.Debugf("Got an error while pinging the docker engine %v, err = %v", dockerClient.DaemonHost(), err)
		return nil, err
	}

	handler := &dockerHandler{
		config: config,
		client: dockerClient,
	}

	return handler, nil
}

func (h *dockerHandler) handle() error {

	time.Sleep(h.config.StartPeriod)
	err := h.getAndRestartUnhealthyContainers()
	if err != nil {
		log.Warnf("Got an exception while trying to restart unhealthy container, err = %v", err)
	}
	for range time.Tick(h.config.Interval) {
		err := h.getAndRestartUnhealthyContainers()
		if err != nil {
			log.Warnf("Got an exception while grying to restart unhealthy container, err = %v", err)
		}
	}

	return nil
}

func (h *dockerHandler) getAndRestartUnhealthyContainers() error {
	containersToRestart, err := h.getUnhealthyContainers()
	if err != nil {
		log.Debugf("Got an exception while getting unhealthy containers, err = %v", err)
		return err
	}
	log.Debugf("Unhealthy containers to restart: %v", containersToRestart)
	for _, c := range containersToRestart {
		counterContainerRestart.Inc()
		err = h.restartContainer(c)
		if err != nil {
			counterContainerRestartFailure.Inc()
			return err
		}
	}
	return nil
}

func (h *dockerHandler) getUnhealthyContainers() ([]restartContainer, error) {

	containerToRestart := make([]restartContainer, 0, 0)
	counterContainerCheck.Inc()
	f := filters.NewArgs()
	f.Add("health", "unhealthy")
	if h.config.ContainerLabel != "all" {
		f.Add("label", h.config.ContainerLabel+"=true")
	}
	containers, err := h.client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		log.Warnf("Got an error while listing the containers, err = %v", err)
		counterContainerCheckFailure.Inc()
		return containerToRestart, err
	}
	for _, c := range containers {

		isManaged, killTimeout, err := h.getManagedContainerAndKillTimeoutFromLabels(c.Labels)
		if err != nil {
			log.Warnf("Got an exception while processing container labels %v, err = %v", c.Labels, err)
			counterContainerCheckFailure.Inc()
			continue
		}
		if isManaged {
			log.Debugf("Container %s is managed", c.ID)
			containerToRestart = append(containerToRestart, restartContainer{containerId: c.ID, killTimeout: killTimeout})
		}
	}

	return containerToRestart, nil
}

func (h *dockerHandler) getManagedContainerAndKillTimeoutFromLabels(labels map[string]string) (isManaged bool, killTimeout time.Duration, err error) {

	if h.config.ContainerLabel == autohealAllContainersValue {
		isManaged = true
	} else if v, ok := labels[h.config.ContainerLabel]; ok && v == autohealManagedValue {
		isManaged = true
	}

	if v, ok := labels[dockerLabelStopTimeout]; ok {
		var kt time.Duration
		kt, err = time.ParseDuration(v)
		if err != nil {
			return
		}
		killTimeout = kt
	}
	return
}

func (h *dockerHandler) restartContainer(container restartContainer) error {
	return h.client.ContainerRestart(context.Background(), container.containerId, &container.killTimeout)
}
