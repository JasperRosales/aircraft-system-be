# ---------- Build stage ----------
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/api


# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/app /app/app

EXPOSE 8080

USER nonroot:nonroot
CMD ["/app/app"]

