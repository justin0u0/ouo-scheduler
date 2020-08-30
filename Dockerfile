FROM debian:stretch-slim

WORKDIR /

COPY _output/bin/ouo-scheduler /usr/local/bin

CMD ["ouo-scheduler"]

