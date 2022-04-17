FROM --platform=linux/amd64  golang:1.17-stretch as build-stage

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=0

RUN go build -o /app/build/svc-notificator ./main.go

FROM --platform=linux/amd64 alpine

COPY --from=build-stage /app/build/svc-notificator /bin/svc-notificator

EXPOSE 8000

ENTRYPOINT ["/bin/svc-notificator"]

