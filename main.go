package main

import (
	"fmt"
	"os"
	"net"
	"io/ioutil"
	"math/rand"
	"time"
	"flag"
	"github.com/couchbase/goutils/logging"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func readAndSave(conn *net.TCPConn, dir string, prefix string) {
	defer conn.Close()
	logging.Infof("Connection from %s\n", conn.RemoteAddr());
	data, error2 := ioutil.ReadAll(conn)
	if error2 != nil {
		fmt.Printf("Cannot write: %s\n", error2);
		os.Exit(1);
	}
	randstr := randomString(4)
	directory := dir + "/" + randstr
	directoryExists := true
	for directoryExists {
		_, err := os.Stat(directory)
		directoryExists = (err == nil)
	}
	error2 = os.Mkdir(directory, 775)
	if error2 != nil {
		conn.Write([]byte("failed create directory,please connect the admin"))
		logging.Errorf("failed create directory")
		return
	}
	f, error2 := os.Create(directory + "/index.txt")
	logging.Infof("save %s for %s", directory, conn.RemoteAddr())
	if error2 != nil {
		conn.Write([]byte("failed create directory"))
		logging.Errorf("failed create file, please connect the admin")
		return
	}
	f.Write(data)
	f.Close()
	conn.Write([]byte(prefix + randstr + "\n"))
}

func main() {
	port := flag.String("port", "9999", "http listen port")
	host := flag.String("host", "127.0.0.1", "bind ip")
	dir := flag.String("dir", "/tmp", "directory")
	prefix := flag.String("prefix", "http://127.0.0.1/", "prefix of saved file,eg: https://p.fht.im/")
	flag.Parse()
	addr, error := net.ResolveTCPAddr("tcp", *host + ":" + *port);
	if error != nil {
		fmt.Printf("Cannot parse \"%s\": %s\n", port, error);
		os.Exit(1);
	}
	listener, error := net.ListenTCP("tcp", addr);
	if error != nil {
		fmt.Printf("Cannot listen: %s\n", error);
		os.Exit(1);
	}
	logging.Infof("listen on %s:%s", *host, *port)

	if !strings.HasSuffix(*prefix, "/") {
		*prefix = *prefix + "/"
	}
	for { // ever...
		conn, error := listener.AcceptTCP();
		if error != nil {
			logging.Errorf("Cannot accept: %s\n", error);
			os.Exit(1);
		}

		go readAndSave(conn, *dir, *prefix);
	}
}
