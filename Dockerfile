FROM golang:1.26.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ecommerce ./cmd/main.go

FROM alpine

WORKDIR /root/

COPY --from=build ./ecommerce ./

ARG PORT=8080
EXPOSE ${PORT}

CMD ["./ecommerce"]
