#########################
# Builder image
#########################

FROM docker.io/golang:latest as BUILDER

ENV USER root

WORKDIR /build

ADD demo .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata

#########################
# Runner
#########################

FROM scratch as RUNNER

WORKDIR /workdir

COPY --from=BUILDER /build/demo /bin/entrypoint

ENV PATH=/bin

ENTRYPOINT ["entrypoint"]
