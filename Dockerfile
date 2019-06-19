# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="Erwan ROUSSEL <erwan.roussel51@gmail.com>"

# Declare environmtents variable
ENV FLUME_DATABASE_ADMIN admin
ENV FLUME_DATABASE_SECRET this_is_a_secret_token  

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/flume-cloud-services/database

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["database"]