FROM golang:1.24-bookworm as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -buildmode=c-shared -o cloudwatch-input.so plugin.go


FROM ghcr.io/selectel/fluent-bit:2025-04-22

COPY --from=build /src/cloudwatch-input.so /fluent-bit/plugins/

ENTRYPOINT ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/plugins/cloudwatch-input.so"]
CMD ["-c", "/fluent-bit/etc/fluent-bit.yaml"]
