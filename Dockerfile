# based on https://github.com/otiai10/gosseract/blob/master/Dockerfile
FROM golang:latest

RUN apt-get update -qq

RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

#ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr

RUN apt-get install -y -qq \
  tesseract-ocr-eng

COPY . /app
RUN cd /app && go build

ENTRYPOINT ["/app/free-iptv"]