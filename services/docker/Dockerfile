FROM golang:1.13

RUN go get -u -d github.com/magefile/mage
WORKDIR $GOPATH/src/github.com/magefile/mage
RUN go run bootstrap.go

RUN curl -sL https://deb.nodesource.com/setup_12.x | bash -
RUN apt-get install -y nodejs
RUN npm i -g aws-cdk
RUN npm i -g typescript