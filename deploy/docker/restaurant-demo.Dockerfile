FROM golang:1.25.10-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY restaurant-demo ./restaurant-demo
COPY shared ./shared

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/restaurant-demo ./restaurant-demo/cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/restaurant-demo /app/restaurant-demo

USER app

EXPOSE 9091

ENTRYPOINT ["/app/restaurant-demo"]
