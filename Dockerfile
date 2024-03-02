FROM golang:1.21.5 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:3.19

WORKDIR /app

COPY --from=build /app/myapp .

EXPOSE 4500
ENTRYPOINT ["/app/myapp"]
