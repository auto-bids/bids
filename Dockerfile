FROM golang:1.21.8 AS build


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:3.19

WORKDIR /app

COPY --from=build /app/myapp .

EXPOSE 4900
ENTRYPOINT ["/app/myapp"]
