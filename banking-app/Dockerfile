# Build stage
FROM golang:1.22.5-alpine as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN if test -e .env; then echo 'found .env'; else mv sample.env .env; fi; \
    go build -v -o /dist/third-party

# Deployment stage
FROM alpine:3.20

RUN apk -U upgrade

WORKDIR /usr/src/app

COPY --from=build /usr/src/app ./

COPY --from=build /dist/third-party /usr/local/bin/third-party

CMD ["third-party"]