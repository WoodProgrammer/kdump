FROM golang:1.19-alpine as builder
WORKDIR /app
RUN apk add --update curl gcc alpine-sdk libpcap-dev
RUN curl -k https://developers.cloudflare.com/cloudflare-one/static/documentation/connections/Cloudflare_CA.pem > cf.pem

RUN cat cf.pem >> /etc/ssl/certs/ca-certificates.crt

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /tcpdump_exporter .

CMD [ "/tcpdump_exporter" ]