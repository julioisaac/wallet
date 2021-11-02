FROM golang:1.17.2 as builder
RUN export GIT_SSL_NO_VERIFY=1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest

ENV API_ENDPOINT='https://api.coingecko.com/api/v3/simple/price'
ENV MONGODB_URI='mongodb://mongo:27017'
ENV OAP_SKY_WALKING='oap:11800'

RUN apk --no-cache add curl
EXPOSE 8000
WORKDIR /root/
COPY --from=builder /app .
CMD ["./app"]