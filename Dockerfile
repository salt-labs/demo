#########################
# Builder image
#########################

FROM docker.io/rust:1.60 as BUILDER

ENV USER root

WORKDIR /build

ADD demo .

RUN rustup \
		target \
		install \
		x86_64-unknown-linux-musl

RUN cargo \
		build \
		--release \
		--target=x86_64-unknown-linux-musl

RUN mkdir /build/artifacts \
 && cp /build/target/x86_64-unknown-linux-musl/release/demo /build/artifacts/demo

#########################
# Runner
#########################

FROM scratch as RUNNER

WORKDIR /workdir

COPY --from=BUILDER /build/artifacts/demo /bin/entrypoint

ENV PATH=/bin

ENTRYPOINT ["entrypoint"]
