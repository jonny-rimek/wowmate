FROM golang:1.13

RUN go get -u -d github.com/magefile/mage
WORKDIR $GOPATH/src/github.com/magefile/mage
RUN go run bootstrap.go

FROM node:13.8
RUN npm i -g aws-cdk