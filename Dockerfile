FROM golang:latest as build

RUN mkdir /app
COPY . /app

WORKDIR /app
RUN go build -o /token-renew

FROM alpine:latest

ENV PATH "$PATH:/usr/local/bin/"

RUN apk add --no-cache libc6-compat 

COPY --from=build /token-renew /token-renew
COPY ./sops.bin /sops

RUN install /token-renew /usr/local/bin/token-renew
RUN install /sops /usr/local/bin/sops

ENTRYPOINT ["/usr/local/bin/token-renew"]

