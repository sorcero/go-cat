FROM alpine
ENTRYPOINT ["/go-cat"]
COPY go-cat /
