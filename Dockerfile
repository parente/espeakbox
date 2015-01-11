FROM progrium/busybox
RUN opkg-install --force-depends espeak lame opus-tools
EXPOSE 8080
COPY server /server
CMD ["/server"]
