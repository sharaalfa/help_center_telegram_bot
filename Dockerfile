FROM golang:1.23-alpine AS builder


ADD . /app

ADD /internal/boundary/*.go /app/internal/boundary/
ADD /internal/app/*.go /app/internal/app/
ADD /internal/config/*.go /app/internal/config/
ADD /internal/gate/mongodb/*.go /app/internal/gate/mongodb/
ADD /internal/gate/postgres/*.go /app/internal/gate/postgres/
ADD /internal/gate/redis/*.go /app/internal/gate/redis/
ADD /pkg/logger/*.go /app/pkg/logger/
ADD /pkg/models/*.go /app/pkg/models/
ADD /config/*.yaml /app/config/
ADD /cmd/*.go /app

WORKDIR /app
RUN go mod download
RUN go build -o main .

EXPOSE 12345

# Запускаем приложение
CMD ["./main", "--config=./config/config_local.yaml"]