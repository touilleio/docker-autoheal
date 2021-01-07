FROM gcr.io/distroless/base

ARG arch=amd64

COPY docker-autoheal-${arch} /docker-autoheal

# Because of access to docker.sock, it's easier to run it as root...
#USER nobody

ENTRYPOINT ["/docker-autoheal"]
EXPOSE 8080
