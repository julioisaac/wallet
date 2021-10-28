FROM golang:1.17.2 as builder
RUN export GIT_SSL_NO_VERIFY=1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add curl
WORKDIR /root/
COPY --from=builder /app .
EXPOSE 8000
CMD ["./app"]