FROM scratch
COPY config /config
COPY Thumb /Thumb
COPY UserData /UserData
ADD coffee /
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/
CMD ["/coffee"]
EXPOSE 30003


