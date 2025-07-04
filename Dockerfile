
FROM golang:1.24-alpine AS buildstage

WORKDIR /Sole-Spot

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk --no-cache add ca-certificates

RUN go build -o /Sole-Spot/sole-spot-app ./cmd1/main.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /Sole-Spot

COPY --from=buildstage /Sole-Spot/sole-spot-app .
COPY --from=buildstage /Sole-Spot/.env .

COPY --from=buildstage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=buildstage /Sole-Spot/templates ./templates/

CMD ["./sole-spot-app"]
