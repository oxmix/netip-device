FROM docker.io/library/golang:1.23-alpine AS builder
ARG GO_PRIVATE
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o device .

FROM docker.io/library/alpine:3.20
LABEL description="https://cloudnetip.com/wiki"
ARG VERSION
ENV VERSION=$VERSION
ARG VERSION_HASH
ENV VERSION_HASH=$VERSION_HASH
# gcompat (glibc) for nvidia-smi
RUN apk --no-cache add gcompat mdadm smartmontools

COPY --from=builder /app/device .

ENTRYPOINT ["./device"]
