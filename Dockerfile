FROM golang:1.12

ADD https://github.com/golang/dep/releases/download/v0.5.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR /$GOPATH/src/github.com/wez470/echo
COPY . ./
RUN dep ensure --vendor-only
EXPOSE 8080

ENTRYPOINT ["go", "run", "."]
