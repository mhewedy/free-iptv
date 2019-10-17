
[![Build Status](https://travis-ci.org/mhewedy/free-iptv.svg?branch=master)](https://travis-ci.org/mhewedy/free-iptv)

To run we need to setup `tesseract` the ocr library:

1. start with running:
```bash
docker run --rm -it -v $(pwd):/work golang:latest
```
2. then inside the container run the commands from the dockerfile:
```bash
https://github.com/otiai10/gosseract/blob/master/Dockerfile
```
3. inside the container build and run the app:
```bash 
cd /work
go run *.go
```