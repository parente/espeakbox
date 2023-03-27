FROM alpine:3.17.2
RUN apk add --no-cache espeak opus lame flac wget && \
    apk del libstdc++
RUN cd /tmp && \
    wget https://downloads.xiph.org/releases/opus/opus-tools-0.1.9.tar.gz --no-check-certificate && \
tar xzf opus-tools-0.1.9.tar.gz && \
    cd opus-tools-0.1.9/ && \
    apk add --no-cache build-base flac-dev opus-dev libogg-dev && \
    ./configure && \
    make && \
    make install && \
    rm -rf /tmp/* && \
    apk del build-base flac-dev opus-dev libogg-dev
COPY server.go /tmp/
RUN cd /tmp && \
    apk add --no-cache go && \
    go build -o /server server.go && \
    apk del go && \
    rm /tmp/server.go
EXPOSE 8080
CMD ["/server"]
