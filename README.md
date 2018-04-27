# fiche-golang
fiche-golang is a simple golang version of [fiche](https://github.com/solusipse/fiche).

# Client-side usage
for example, use public server

```bash
echo "I will always love you" | nc p.fht.im 9999
```
you could get an url to your paste as a response. e.g.:

```
https://p.fht.im/BoUI
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
  -dir string
        directory (default "/tmp")
  -host string
        bind ip (default "127.0.0.1")
  -port string
        http listen port (default "9999")
  -prefix string
        prefix of saved file,eg: https://p.fht.im/ (default "http://127.0.0.1/")
```
for example, I want to run a public fiche server with prefix https://p.fht.im and store data in /data/fiche_data/, command below.
```
./paste_server --prefix https://p.fht.im --dir /data/fiche_data --host "0.0.0.0"
```

## Example nginx config
Add a built-in server for fiche-golang is simple, try to add it yourself or use nginx to server the file. here is an example.
```nginx
server {
    listen 80;
    server_name mysite.com www.mysite.com;
    charset utf-8;

    location / {
            root /data/fiche_data/;
            index index.txt index.html;
    }
}
```
Fiche has no http server built-in, thus you need to setup one if you want to make files available through http.

## TODO
- [ ] Build executable for common platform
- [ ] Add ipv6 support.
- [ ] Add a systemd example
- [ ] Add a dockerfile
- [ ] Maybe add more options such as white-list etc..