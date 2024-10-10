FROM golang:1.22.5 AS builder

WORKDIR /app

COPY . .


RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -gcflags=all="-N -l" -o /dongle

FROM debian:latest
COPY --from=builder /dongle /dongle
COPY ./pay/cert /cert
ENV TZ=Asia/Shanghai \
  DEBIAN_FRONTEND=noninteractive

ENV ENVIRONMENT=production

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
  && echo ${TZ} > /etc/timezone \
  && dpkg-reconfigure --frontend noninteractive tzdata \
  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "/dongle" ]