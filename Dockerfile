FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY iap-token-generator /bin/

ENTRYPOINT ["/bin/iap-token-generator"]
