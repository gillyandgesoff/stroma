FROM debian:wheezy
ADD ./stroma /usr/bin/stroma
ENTRYPOINT exec stroma -maxconn=10
EXPOSE 3000
