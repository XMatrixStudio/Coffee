FROM scratch
COPY config /config
ADD coffee /
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/
CMD ["/coffee"]
EXPOSE 30003


