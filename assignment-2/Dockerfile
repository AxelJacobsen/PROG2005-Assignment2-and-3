
FROM golang:1.17 AS builder

LABEL maintainer="aejacobs@stud.ntnu.no"
LABEL stage=builder

WORKDIR /go/src/app

# Copy relevant folders into container
COPY ./endpoints /go/src/app/endpoints
COPY ./variables /go/src/app/variables

COPY ./go.sum /go/src/app/go.sum
COPY ./go.mod /go/src/app/go.mod
COPY ./alpha3.json /go/src/app/alpha3.json
COPY ./main.go /go/src/app/main.go
COPY ./ass2-service-account.json /go/src/app/ass2-service-account.json

# Build command
RUN go build -o main .

#exposing port
EXPOSE 8080

# Instantiate binary
CMD ["./main"]
