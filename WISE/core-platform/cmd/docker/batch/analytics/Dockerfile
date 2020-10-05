FROM alpine:3.10.3

LABEL maintainer="Jerry Warren<jerry@wise.us>"
LABEL Description="wise company account batch analytics container"

# add a bit of security here and run as user, instead of root
RUN apk update && \
    # Add ca certificates for secure SSL
    apk add ca-certificates && \
    # Add timezone information
    apk add --no-cache tzdata && \
    rm -rf /var/cache/apk/* && \
    # add wise_user so not running as root
    adduser -S -D -H -h /app wise_user

USER wise_user

ADD ./main /app/
WORKDIR /app

CMD ["./main"]