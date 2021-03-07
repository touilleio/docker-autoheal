Docker autoheal
====

This tool monitors and restarts unhealthy docker containers.
This functionality was proposed to be included with the addition of `HEALTHCHECK`, however didn't go through.
This tool is a workaround till there is native support for `--restart-on-unhealthy` or similar ([https://github.com/docker/docker/pull/22719](https://github.com/docker/docker/pull/22719)).

# Disclaimer

This `docker-autoheal` is a rewrite of [the excellent docker-autoheal](https://github.com/willfarrell/docker-autoheal)
from [Will Farrell](https://twitter.com/willfarrell), but in [golang](https://golang.org/).
It is fully compliant with `willfarrell/docker-autoheal`, plus few goodies such as metrics (see below)

## Compatibility

- `AUTOHEAL_INTERVAL`, `AUTOHEAL_START_PERIOD`, `AUTOHEAL_DEFAULT_STOP_TIMEOUT` and `CURL_TIMEOUT`
are duration and are treated as seconds if no unit provided. For more precision `5` can be written `5s`.
- `DOCKER_SOCK` env variable can be written to more standard `DOCKER_HOST`. `DOCKER_SOCK` can still be used.

# Configuration

## Environment variable to configure docker-autoheal

| Env var name | Default | Description |
|--------------|---------|-------------|
| AUTOHEAL_CONTAINER_LABEL | autoheal | Specify the name of the label optin'ing to autoheal process. A special value `all` means all containers are watched |
| AUTOHEAL_INTERVAL | 5s | Interval at which containers health is checked |
| AUTOHEAL_START_PERIOD | 0s | Warmup time before running the first check |
| AUTOHEAL_DEFAULT_STOP_TIMEOUT | 10s | Time to give to containers before shutting down |
| DOCKER_HOST (or DOCKER_SOCK) | /var/run/docker.sock | Path/URI of docker socket |
| CURL_TIMEOUT | 30s | Timeout when interacting with docker |

## Additional configuration for the containers

| Label name | Description |
|------------|-------------|
| autoheal   | if `true` or `AUTOHEAL_CONTAINER_LABEL=all`, means this container is watched |
| autoheal.stop.timeout | Per containers override for stop timeout seconds during restart, as a duration, e.g. `120s` |

# Metrics

| Metrics name | Type | Description |
|--------------|------|-------------|
| check_count | counter | Count how many times containers have been checked |
| check_failure_count | counter | Count how many failures happened trying to check containers |
| restart_count | counter | Count how many containers have been restated |
| restart_failure_count | counter | Count how many restart failures happened |
