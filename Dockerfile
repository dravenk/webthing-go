FROM golang

LABEL MAINTAINER='longxianwen@outlook.com'

WORKDIR /go/src/github.com/dravenk/webthing-go
COPY . .

RUN go install examples/single-thing/single-thing.go

EXPOSE 8888
CMD ["single-thing"]