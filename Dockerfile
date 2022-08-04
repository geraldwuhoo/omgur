# Builder
FROM docker.io/golang:1.19-alpine3.16 AS builder

WORKDIR /build
COPY go.mod go.sum /build/
ENV CGO_ENABLED=0 \
    GOOS=linux
RUN go mod download
COPY . /build/
RUN go build -ldflags "-w -s" -trimpath ./cmd/omgur

# Clean image
FROM docker.io/alpine:3.16
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/omgur .

USER nobody:nobody

CMD ["./omgur"]
EXPOSE 8080
