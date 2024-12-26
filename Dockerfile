FROM golang:1.21-alpine
WORKDIR /app
RUN apk add --no-cache git build-base
COPY go.mod ./
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o chord-node
EXPOSE 8080
CMD ["./chord-node"]