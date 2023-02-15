FROM golang:1.15.7-buster

RUN echo $PWD

COPY . /prod/
COPY cons /cons/
COPY prod_cons.sh /

RUN mkdir -p /tmp
# RUN touch /tmp/echo.sock


RUN chmod +x /prod_cons.sh

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV GOPATH=""

ENTRYPOINT ["/prod_cons.sh"]
