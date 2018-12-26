FROM scratch
COPY config /config
VOLUME ["/UserData", "/Thumb"]
ADD coffee /
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/
CMD ["/coffee"]
EXPOSE 30070
