FROM golang:1.21 AS builder

COPY . /src
RUN mkdir -p   /data/conf
COPY ./configs/config.yaml /data/conf
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY --from=builder /data/conf /data/conf

WORKDIR /app

EXPOSE 8088
EXPOSE 9000
VOLUME /data/conf

CMD ["./minerva", "-conf", "/data/conf"]
