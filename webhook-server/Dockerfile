FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webhook

FROM gcr.io/distroless/base-debian12

WORKDIR /

# Copy cert directory with correct directory permissions
COPY --from=builder --chown=nonroot:nonroot --chmod=0555 /app/certs /certs

# Fix file permissions explicitly
COPY --from=builder --chown=nonroot:nonroot --chmod=0444 /app/certs/tls.crt /certs/tls.crt
COPY --from=builder --chown=nonroot:nonroot --chmod=0400 /app/certs/tls.key /certs/tls.key

COPY --from=builder /app/webhook /webhook

USER nonroot:nonroot
EXPOSE 8443
ENTRYPOINT ["/webhook"]
