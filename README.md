# fiche-golang
fiche-golang is a command line pastebin for sharing terminal output. inspired by [fiche](https://github.com/solusipse/fiche).

# docker start
```
sudo docker run -d --restart=unless-stopped -p 8080:80 -p 9999:9999 -e prefix=http://`curl -s https://api.ipify.org`:8080 imfht/fiche-golang:latest
```
tip: the docker image < 10mb. use docker is a good choise

# New Features
Old fiche only support txt file, fiche-golang support lots of file format and will auto add extension. for example:

## Fiche

```bash
➜  /tmp cat simple_jpg.jpg| nc termbin.com 9999
http://termbin.com/0zuw
➜  /tmp curl http://termbin.com/0zuw

����%
```
## Fiche-golang
```bash
➜  /tmp cat simple_jpg.jpg| nc termbin.co 9999
https://termbin.co/t4OZ/index.jpg
```
and click [https://termbin.co/t4OZ/index.jpg](https://termbin.co/t4OZ/index.jpg) you'll see the image.

# Client-side usage
for example, use public server

```bash
echo "I will always love you" | nc termbin.co 9999
```
you could get an url to your paste as a response. e.g.:

```
https://termbin.co/BoUI
```
# Server-side useage
## Installation
1. Clone
```
git clone https://github.com/imfht/fiche-golang
```
2. build
```
go build main.go
```

## Usage
```
Usage of ./paste_server:
  -prefix string
        prefix of saved file,eg: https://termbin.co/ (default "http://127.0.0.1/")
```

# Release
if you are using amd 64 Linux, you can download [https://github.com/imfht/fiche-golang/raw/master/filche-golang](https://github.com/imfht/fiche-golang/raw/master/filche-golang) and run directly.

## TODO
- [ ] Build executable for common platform
- [ ] Add ipv6 support.
- [ ] Add a systemd example
- [ ] Add a dockerfile
- [ ] Maybe add more options such as white-list etc..
- [x] Maybe I should limit file size to avoid memory leak.‰‰‰
