FROM alpine:latest
RUN apk add netcat-openbsd
COPY ./netcat/netcat_echo.sh netcat_echo.sh
RUN chmod +x /netcat_echo.sh
ENTRYPOINT ["/bin/sh"]