FROM golang:1.22.2-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o merchant_bank_payment_go_api /build/cmd/app

EXPOSE 8000

CMD ["/build/merchant_bank_payment_go_api"]
