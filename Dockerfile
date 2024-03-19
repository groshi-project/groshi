FROM golang:1.22 as builder

WORKDIR /groshi-build

# copy sources:
COPY go.mod go.sum ./
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./docs ./docs
COPY main.go ./

# setup go env:
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

# install dependencies:
RUN --mount=type=cache,target=/gomod-cache go mod download

# build groshi:
ENV GOOS=linux
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache go build -x -o ./groshi

FROM alpine:latest as runner
COPY --from=builder /groshi-build/groshi /usr/bin/groshi
EXPOSE 8080
ENTRYPOINT ["/usr/bin/groshi"]