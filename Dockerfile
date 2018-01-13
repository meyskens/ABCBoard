FROM debian:sid

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y pkg-config libwebkit2gtk-4.0-dev libgtk-3-dev libasound-dev curl git

RUN curl -sL https://deb.nodesource.com/setup_9.x | bash - && apt-get install -y nodejs

ENV GO_VERSION 1.9.2
RUN curl -fSL "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar xzC /usr/local
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

COPY ./ /go/src/github.com/meyskens/abcboard/
WORKDIR /go/src/github.com/meyskens/abcboard/


RUN go get
RUN cd ./frontend && npm install && cd ../
RUN ./build.sh