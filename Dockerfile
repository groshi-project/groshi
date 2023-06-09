FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY  groshi.go ./
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/bin/groshi

EXPOSE 8080

ENTRYPOINT ["/usr/bin/groshi"]
