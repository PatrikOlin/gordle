FROM golang:1.17 AS build-env
ENV PORT 4040

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o gordle

FROM debian:buster

WORKDIR /

COPY --from=build-env /app/gordle .
COPY --from=build-env /app/_data.db .

EXPOSE $PORT

CMD ["./gordle"]
