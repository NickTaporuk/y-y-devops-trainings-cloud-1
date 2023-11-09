# Use the official Golang image as a base image
FROM golang:1.21.3 AS builder

# Force the go compiler to use modules
ENV GO111MODULE=on

# Set the working directory inside the container
WORKDIR /go/src/github.com/NickTaporuk/devops-training

COPY ./catgpt .

RUN GOOS=linux GOARCH=amd64 go build -a -o catgpt main.go && ls -lah

######### Start a new stage from scratch #######
FROM alpine:3.7

#WORKDIR /root

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/NickTaporuk/devops-training/catgpt /bin/catgpt

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* \

EXPOSE 8080
EXPOSE 9090

#CMD ["ls -lah"]
CMD ["catgpt"]
