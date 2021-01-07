Docker autoheal
====

This tool monitors and restarts unhealthy docker containers.
This functionality was proposed to be included with the addition of `HEALTHCHECK`, however didn't go through.
This tool is a workaround till there is native support for `--restart-on-unhealthy` or similar ([https://github.com/docker/docker/pull/22719](https://github.com/docker/docker/pull/22719)).

# Disclaimer 
This `docker-autoheal` is a rewrite of [the excellent docker-autoheal](https://github.com/willfarrell/docker-autoheal)
from [Will Farrell](https://twitter.com/willfarrell), but in [golang](https://golang.org/).
It is fully compliant with `willfarrell/docker-autoheal`.

# Configuration




ENV Defaults

AUTOHEAL_CONTAINER_LABEL=autoheal
AUTOHEAL_INTERVAL=5   # check every 5 seconds
AUTOHEAL_START_PERIOD=0   # wait 0 seconds before first health check
AUTOHEAL_DEFAULT_STOP_TIMEOUT=10   # Docker waits max 10 seconds (the Docker default) for a container to stop before killing during restarts (container overridable via label, see below)
DOCKER_SOCK=/var/run/docker.sock   # Unix socket for curl requests to Docker API
CURL_TIMEOUT=30     # --max-time seconds for curl requests to Docker API

Optional Container Labels

autoheal.stop.timeout=20        # Per containers override for stop timeout seconds during restart



