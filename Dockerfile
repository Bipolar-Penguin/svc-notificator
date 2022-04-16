FROM golang:1.17-stretch as build-stage

WORKDIR /app

COPY . /app

RUN go build -o /app/build/svc-notificator ./main.go

FROM alpine

COPY --from=build-stage /app/build/svc-notificator /bin/svc-notificator

EXPOSE 8000

ENTRYPOINT ["/bin/svc-notificator"]

