#########################
# Builder image
#########################

FROM docker.io/alpine:latest

RUN echo "Building container..." \
 && touch /tmp/hello-world
