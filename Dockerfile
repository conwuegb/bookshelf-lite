FROM golang:latest AS build

# Copy source code into the container to WORKDIR
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.mod .

# Disable C Go in order to produce a statically linked binary
ENV CGO_ENABLED=0

# Download module dependencies
RUN go mod tidy -v

# Build the binary `swagger`
RUN go build -a -installsuffix cgo -o swagger .

# Create a runtime stage with port 8080 exposed using TCP
FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8080/tcp
ENTRYPOINT ["./swagger"]
