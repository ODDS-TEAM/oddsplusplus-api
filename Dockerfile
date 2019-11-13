#Build State
# Base Image
FROM golang:1.11-alpine AS build-state

# Install Git
RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Run Test
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/api


# Deploy State
FROM alpine
RUN apk update && apk upgrade
RUN apk add curl
RUN apk add python
WORKDIR /app
RUN mkdir -p files/tavi50 && mkdir image && mkdir font
COPY --from=build-state /go/bin/api /app
