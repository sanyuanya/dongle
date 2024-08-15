FROM golang:1.22.5 as builder

WORKDIR /app

COPY . .

COPY ./pay/cert /cert
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /dongle

FROM debian:latest
COPY --from=builder /dongle /dongle

ENV TZ=Asia/Shanghai \
  DEBIAN_FRONTEND=noninteractive

ENV ENVIRONMENT=production

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
  && echo ${TZ} > /etc/timezone \
  && dpkg-reconfigure --frontend noninteractive tzdata \
  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT [ "/dongle" ]