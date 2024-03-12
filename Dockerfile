FROM golang:1.21.5-alpine AS build

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o app cmd/main.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=build /app/app .

ENTRYPOINT ["./app"]