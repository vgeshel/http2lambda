FROM kiasaki/alpine-golang
MAINTAINER vgeshel@gmail.com

WORKDIR /gopath/src/app
ADD . /gopath/src/app/
RUN go get app

CMD []
ENTRYPOINT ["/gopath/bin/app"]
