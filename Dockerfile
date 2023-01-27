FROM golang:1.19.5-alpine3.17
RUN mkdir /app
WORKDIR /app

#Copy all the files inside container work dir
COPY . .

COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /build cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]

