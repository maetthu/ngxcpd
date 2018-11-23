FROM alpine
LABEL maintainer="Matthias Blaser <git@mooch.ch>"
COPY dist/linux_amd64/ngxcpd /ngxcpd
ENTRYPOINT ["/ngxcpd"]