FROM --platform=$BUILDPLATFORM golang:1.22.5-bookworm AS build
WORKDIR /src

RUN apt-get update && apt-get install -y gcc g++-riscv64-linux-gnu

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=./go.sum,target=./go.sum \
    --mount=type=bind,source=./go.mod,target=./go.mod \
    go mod download -x

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 GOARCH=riscv64 GOOS=linux CC=riscv64-linux-gnu-gcc go build -o /bin/app ./cmd/dapp

# runtime stage: produces final image that will be executed
FROM ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
RUN <<EOF
set -e
apt-get update
apt-get install -y --no-install-recommends \
  busybox-static=1:1.30.1-7ubuntu3
rm -rf /var/lib/apt/lists/* /var/log/* /var/cache/*
useradd --create-home --user-group dapp
EOF

ARG MACHINE_EMULATOR_TOOLS_VERSION=0.16.1
ADD https://github.com/cartesi/machine-emulator-tools/releases/download/v${MACHINE_EMULATOR_TOOLS_VERSION}/machine-emulator-tools-v${MACHINE_EMULATOR_TOOLS_VERSION}.deb /
RUN dpkg -i /machine-emulator-tools-v${MACHINE_EMULATOR_TOOLS_VERSION}.deb \
  && rm /machine-emulator-tools-v${MACHINE_EMULATOR_TOOLS_VERSION}.deb

ENV PATH="/opt/cartesi/bin:${PATH}"

WORKDIR /var/opt/cartesi-app
COPY --from=build /bin/app app

# Change ownership of the application directory to the dapp user
RUN chown -R dapp:dapp .

ENV ROLLUP_HTTP_SERVER_URL="http://127.0.0.1:5004"

ENTRYPOINT ["rollup-init"]

CMD ["/var/opt/cartesi-app/app"]