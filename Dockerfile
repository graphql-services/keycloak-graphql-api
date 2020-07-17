FROM golang as builder
WORKDIR /go/src/github.com/graphql-services/id
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /tmp/app server.go

FROM jakubknejzlik/wait-for as wait-for

FROM alpine:3.5

WORKDIR /app

COPY --from=builder /tmp/app /usr/local/bin/app
# COPY --from=builder /go/src/github.com/graphql-services/id/schema.graphql /app/schema.graphql

# RUN apk --update add docker
RUN apk --update add ca-certificates

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN chmod +x /usr/local/bin/app

ENTRYPOINT []
CMD [ "/bin/sh", "-c", "app server" ]