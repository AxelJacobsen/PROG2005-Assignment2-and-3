
FROM golang:1.17 AS builder

LABEL maintainer="aejacobs@stud.ntnu.no"
LABEL stage=builder

WORKDIR /go/src/app

# Copy relevant folders into container
COPY ./cases /go/src/app/code/endpoints/cases
COPY ./notification /go/src/app/code/endpoints/notification
COPY ./policy /go/src/app/code/endpoints/policy
COPY ./public /go/src/app/code/endpoints/public
COPY ./inner /go/src/app/code/endpoints/sharedFuncs/inner
COPY ./webhookFuncs /go/src/app/code/endpoints/sharedFuncs/webhookFuncs
COPY ./status /go/src/app/code/code/endpoints/status
COPY ./constants /go/src/app/code/constants

COPY ./go.sum /go/src/app/go.sum
COPY ./go.mod /go/src/app/go.mod
COPY ./alpha3.json /go/src/app/code/constants/alpha3.json
COPY ./main.go /go/src/app/main.go
COPY ./ass3-service-account.json /go/src/app/ass3-service-account.json

# Build command
RUN go build -o main .

#exposing port
EXPOSE 8080

# Instantiate binary
CMD ["./main"]
