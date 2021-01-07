package main

import (
	"github.com/sqooba/go-common/logging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDockerEventsListener(t *testing.T) {

	config := envConfig{
		LogLevel:       "debug",
		Interval:       10 * time.Second,
		ContainerLabel: autohealAllContainersValue,
	}
	err := logging.SetLogLevel(log, config.LogLevel)

	dockerHandler, err := newDockerHandler(config)
	assert.Nil(t, err)

	err = dockerHandler.handle()
	assert.Nil(t, err)
}
