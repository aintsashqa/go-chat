FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

RUN mkdir -p templates

# Copy binary from build to main folder
RUN cp /build/main .
RUN cp /build/.env .
COPY ./templates ./templates

# Export necessary port
# EXPOSE ${SERVICE_PORT}

# Command to run when starting the container
CMD ["/dist/main"]
