FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

RUN go mod download

COPY . .

RUN go build -o login ./cmd/login

FROM alpine

WORKDIR /app

COPY --from=builder /build/login login
COPY ./cmd/login/migrate migrate/

CMD ["/app/login"]