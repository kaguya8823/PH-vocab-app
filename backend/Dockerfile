FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:3.22

WORKDIR /app
COPY --from=build /app/server ./server

EXPOSE 3002

CMD ["./server"]
