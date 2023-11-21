FROM golang:alpine as build

WORKDIR /app
COPY . /app

RUN go build -o /token-renew

FROM alpine:latest

ENV PATH "$PATH:/usr/local/bin/"

COPY --from=build /token-renew /token-renew
COPY ./sops.bin /sops

RUN install /token-renew /usr/local/bin/token-renew && install /sops /usr/local/bin/sops

ENTRYPOINT ["/usr/local/bin/token-renew"]

