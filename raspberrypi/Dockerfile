FROM golang:1.16-buster AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN go build -o /bmp280
EXPOSE 9110
ENTRYPOINT ["/bmp280"]