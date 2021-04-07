FROM golang:latest
ENV GOPROXY https://goproxy.cn,direct
WORKDIR .
COPY . /home/PreTest
RUN go build .
EXPOSE 80
ENTRYPOINT ["./"]