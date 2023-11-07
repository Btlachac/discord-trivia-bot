FROM golang:1.21 as builder

WORKDIR /src

ADD . /src

RUN go build -o /app ./cmd

FROM scratch
COPY --from=builder /app /app
CMD ["/app"]