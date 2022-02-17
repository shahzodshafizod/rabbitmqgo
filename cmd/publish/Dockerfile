FROM golang:alpine AS builder

WORKDIR /build

RUN apk --no-cache add ca-certificates

COPY ../../ .

ENV GO111MODULE=on
RUN go mod tidy
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o publish cmd/publish/main.go

FROM scratch

COPY --from=builder /build/publish /

ENTRYPOINT ["/publish"]
