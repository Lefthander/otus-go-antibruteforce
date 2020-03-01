FROM golang:1:14 AS pre-build
COPY . /abf
WORKDIR /abf
RUN make build-server

FROM alpine:3.7
COPY --from=pre-build /abf/bin/server/ abf/server/
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
