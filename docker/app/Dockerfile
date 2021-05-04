FROM golang:1.15.2-alpine
RUN apk update && apk add git

RUN mkdir /go/src/Gacha

WORKDIR /go/src/Gacha

ADD ../.. /go/src/Gacha
EXPOSE 8090