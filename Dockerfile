##########################
# builder image definition
##########################
FROM golang:1.12-alpine as builder

RUN apk update && apk add git

WORKDIR /ues
COPY . .

RUN go build -i -o ues github.com/nathanows/ues/server

RUN mkdir -p dist && cp -rfv ues dist

##########################
# runtime image definition
##########################
FROM alpine:3.9

RUN adduser -D -h /home/appuser appuser appuser
USER appuser
WORKDIR /home/appuser

COPY --from=builder /ues/dist/. .

CMD ["sh", "-c", "./ues"]
