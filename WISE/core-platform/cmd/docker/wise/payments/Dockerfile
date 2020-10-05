FROM alpine:3.10.3

LABEL maintainer="Jerry Warren<jerry@wise.us>"
LABEL Description="wise company payments container"

RUN apk update && \
    # Add ca certificates for secure SSL
    apk add ca-certificates openssl && \
    # add wise_user so not running as root
    adduser -S -D -H -h /app wise_user && \
    # install self-signed ssl certificates
    install -o wise_user -g nogroup -m 0700 -d /app/ssl && \
      openssl req -subj '/CN=localhost' \
                  -x509 -newkey rsa:4096 \
                  -nodes -keyout /app/ssl/key.pem \
                  -out /app/ssl/cert.pem \
                  -days 3650 && \
      chown -R wise_user: /app/ssl/* && \
      chmod 0400 /app/ssl/key.pem && \
    # clean up apk
    apk del openssl && rm -rf /var/cache/apk/*

# add a bit of security here and run as user, instead of root
USER wise_user

# Copy over the app files
COPY --chown=wise_user:nogroup \
     ./main \
     ./web-payments \
     /app/

WORKDIR /app

CMD ["./main"]
