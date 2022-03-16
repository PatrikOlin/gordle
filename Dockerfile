FROM golang:1.17-alpine AS build-env

ENV PORT 4040
ENV APP_NAME gordle

WORKDIR /app

COPY . ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o $APP_NAME *.go

FROM alpine:3.14

WORKDIR /

COPY --from=build-env /$APP_NAME .

EXPOSE $PORT

CMD ./$APP_NAME
