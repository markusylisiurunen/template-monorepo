FROM golang:1.18-buster as builder

WORKDIR /app/backend/service

COPY go.mod ./
COPY go.sum ./

COPY package/go /app/package/go

RUN go mod download

COPY pkg ./pkg

RUN mkdir -p bin && \
  GOOS=linux go build -mod=readonly -v -o ./bin/service github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg

FROM debian:buster-slim

WORKDIR /app

RUN set -x && \
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

COPY static ./static

COPY migrations ./migrations

COPY --from=builder /app/backend/service/bin/service ./bin/service

ENV SWIFTBEAVER_ENV=production

ENV PORT=80

CMD ["./bin/service"]
