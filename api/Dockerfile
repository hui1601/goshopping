FROM golang:1.22.2-bullseye
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/app
CMD ["/app/app"]