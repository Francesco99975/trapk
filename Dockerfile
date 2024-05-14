FROM golang:1.21.5-alpine3.18 AS build

RUN apk --no-cache add gcc g++ make git

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN mv .prod.env .env

RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/trapk ./cmd/server/*.go

FROM alpine:3.18

RUN apk update && apk upgrade && apk --no-cache add ca-certificates

WORKDIR /go/bin

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/.env /go/bin/
COPY --from=build /go/src/app/data /go/bin/data
COPY --from=build /go/src/app/static /go/bin/static

EXPOSE 8099

ENTRYPOINT /go/bin/trapk --port 8099