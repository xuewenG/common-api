FROM golang:1.25.0-alpine AS builder
ARG VERSION
ARG COMMIT_ID
WORKDIR /app
RUN apk add --no-cache build-base tzdata
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build VERSION=${VERSION} COMMIT_ID=${COMMIT_ID}

FROM alpine
WORKDIR /app
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/dist/common-api /app/common-api
RUN chmod +x /app/common-api
CMD ["/app/common-api"]
