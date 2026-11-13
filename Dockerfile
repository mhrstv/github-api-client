FROM golang:1.27 AS build
LABEL authors="mhrstv"
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /github-api-client ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /github-api-client /app/github-api-client
COPY usernames.txt /app/usernames.txt

ENTRYPOINT ["/app/github-api-client"]
CMD ["usernames.txt"]