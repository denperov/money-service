FROM alpine:3.10.2
RUN  apk add --no-cache bash curl jq
WORKDIR /app
COPY /test/accounts/api_test.sh ./
ENTRYPOINT ["./api_test.sh"]
