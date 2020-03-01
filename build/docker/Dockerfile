FROM golang:1.14 as build-env
ENV NAME "abf"
WORKDIR /${NAME}
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env AS build
ENV NAME "abf"
WORKDIR /opt/${NAME}
COPY . .
RUN CGO_ENABLED=0 go build -o bin/${NAME}-srv

FROM alpine
ENV NAME "abf"
WORKDIR /${NAME}
COPY --from=build ${NAME}/bin/${NAME} ./${NAME}
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
EXPOSE 8999
CMD ./${NAME} abf-srv
