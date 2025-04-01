FROM golang:1.24 as go-builder

ENV CGO_ENABLED 0

COPY . /src
RUN cd /src/cmd/app && go build

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Europe/Budapest

RUN cp /usr/share/zoneinfo/Europe/Budapest /etc/localtime

COPY --from=go-builder /src/cmd/app/app /app/
WORKDIR /app

CMD [ "./app" ]