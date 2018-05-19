FROM golang:latest

RUN mkdir -p ${GOPATH}/src/github.com/SimoneStefani/thesis-algorithms
ADD . ${GOPATH}/src/github.com/SimoneStefani/thesis-algorithms/
WORKDIR ${GOPATH}/src/github.com/SimoneStefani/thesis-algorithms

RUN pwd
RUN go build -o thesis .
