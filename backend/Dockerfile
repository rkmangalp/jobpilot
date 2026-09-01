FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY backend ./backend
RUN go build -o /jobpilot-api ./backend/cmd/api
FROM alpine:3.20
COPY --from=build /jobpilot-api /jobpilot-api
EXPOSE 8080
ENTRYPOINT ["/jobpilot-api"]
