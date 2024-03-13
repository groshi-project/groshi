FROM golang:1.22 as builder

WORKDIR /groshi-build

COPY go.mod go.sum ./
COPY main.go ./
COPY ./internal ./internal

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -x -o ./groshi

FROM alpine:latest as runner
COPY --from=builder /groshi-build/groshi /usr/bin/groshi
EXPOSE 8080
ENTRYPOINT ["/usr/bin/groshi"]
