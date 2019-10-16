
To run we need to setup `tesseract` the ocr library:
https://github.com/otiai10/gosseract/blob/master/Dockerfile

1. start with running:
```bash
docker run --rm -it -v $(pwd):/work golang:latest
```
2. install `tesseract` as in the Dockerfile above:
```bash
apt-get install -y -qq libtesseract-dev libleptonica-dev
```
3. inside the container build and run the app:
```bash 
cd /work
go run main.go
```