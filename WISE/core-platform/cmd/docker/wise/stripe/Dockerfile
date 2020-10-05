FROM alpine:3.10.3

LABEL maintainer="Jerry Warren<jerry@wise.us>"
LABEL Description="wise company money request container"

RUN apk update && \
    # Add ca certificates for secure SSL
    apk add ca-certificates && \
    rm -rf /var/cache/apk/* && \
    # add wise_user so not running as root
    adduser -S -D -H -h /app wise_user

# add a bit of security here and run as user, instead of root
USER wise_user

# Copy over the app files
COPY ./main /app/

WORKDIR /app

CMD ["./main"]
