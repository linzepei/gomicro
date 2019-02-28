FROM alpine:3.2
ADD conf /conf
ADD html /html
ADD ihomeweb /ihomeweb
WORKDIR /
ENTRYPOINT [ "/ihomeweb" ]
EXPOSE 8888