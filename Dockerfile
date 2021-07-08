FROM golang:1.16

ENV GOPROXY https://goproxy.cn,direct

WORKDIR $GOPATH/src/GinBlog
COPY . $GOPATH/src/GinBlog
RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -o go-gin-blog .

EXPOSE 9099
ENTRYPOINT ["./go-gin-blog"]
