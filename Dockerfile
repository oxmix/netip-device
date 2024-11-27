FROM docker.io/library/golang:1.22-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o device .

FROM docker.io/library/alpine:3.18
LABEL description="https://oxmix.net"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
# gcompat (glibc) for nvidia-smi
RUN apk --no-cache add gcompat mdadm smartmontools

COPY --from=builder /app/device .

ENTRYPOINT ["./device"]
