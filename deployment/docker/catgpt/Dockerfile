# Use the official Golang image as a base image
FROM golang:1.21.3 AS builder

# Set the working directory inside the container
WORKDIR /go/src/github.com/NickTaporuk/devops-training

COPY ./catgpt .

RUN go build -o catgpt main.go && chmod +x catgpt

######## Start a new stage from scratch #######
FROM gcr.io/distroless/static-debian12:latest-amd64

WORKDIR /root

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/NickTaporuk/devops-training /bin/catgpt

EXPOSE 8080

CMD ["catgpt"]
