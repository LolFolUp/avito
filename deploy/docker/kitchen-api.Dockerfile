FROM golang:1.25.10-alpine AS go-base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

FROM go-base AS goose-builder

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/goose github.com/pressly/goose/v3/cmd/goose

FROM go-base AS kitchen-builder

COPY kitchen ./kitchen
COPY shared ./shared

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/kitchen-api ./kitchen/cmd

FROM alpine:3.22 AS kitchen-api

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=kitchen-builder /out/kitchen-api /app/kitchen-api

USER app

EXPOSE 8080

ENTRYPOINT ["/app/kitchen-api"]

FROM alpine:3.22 AS migrate

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=goose-builder /out/goose /app/goose
COPY db/migrations /migrations

USER app

ENTRYPOINT ["/app/goose"]
