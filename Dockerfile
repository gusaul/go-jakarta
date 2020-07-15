FROM golang:1.14-alpine AS builder

ENV GO111MODULE=on

WORKDIR /src

COPY . ./
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o product main.go


FROM alpine

COPY --from=builder /src/product ./
ENTRYPOINT ["./product"]