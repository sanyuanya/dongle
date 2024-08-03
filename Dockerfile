FROM golang:1.22.5 as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /dongle

FROM debian:latest
COPY --from=builder /dongle /dongle

ENTRYPOINT [ "/dongle" ]