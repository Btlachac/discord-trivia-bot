FROM golang:1.21 as builder

ENV CGO_ENABLED 0

WORKDIR /src

ADD . /src

RUN go build -o /app ./cmd/main.go

FROM scratch

COPY --from=builder /app /app
COPY --from=builder /src/migrations ./migrations

CMD ["/app"]