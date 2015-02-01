FROM gliderlabs/alpine:3.1
RUN apk-install espeak opus lame flac && \
    apk del libstdc++
RUN cd /tmp && \
    wget http://downloads.xiph.org/releases/opus/opus-tools-0.1.9.tar.gz && \
    tar xzf opus-tools-0.1.9.tar.gz && \
    cd opus-tools-0.1.9/ && \
    apk-install build-base flac-dev opus-dev libogg-dev && \
    ./configure && \
    make && \
    make install && \
    rm -rf /tmp/* && \
    apk del build-base flac-dev opus-dev libogg-dev
COPY server.go /tmp/
RUN cd /tmp && \
    apk-install go && \
    go build -o /server server.go && \
    apk del go && \
    rm /tmp/server.go
EXPOSE 8080
CMD ["/server"]
