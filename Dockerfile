FROM golang:1.17-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o vsratoslav

FROM ubuntu

WORKDIR /app

COPY --from=builder /app/vsratoslav .
COPY /public /app

CMD ["./vsratoslav --help"]