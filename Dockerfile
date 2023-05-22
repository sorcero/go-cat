FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd/go-cat && go build -o /go-cat


FROM alpine
COPY --from=builder /go-cat /bin/go-cat
ENTRYPOINT ["/bin/go-cat"]

