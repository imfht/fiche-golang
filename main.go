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
	"gopkg.in/h2non/filetype.v1"
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

func getFileExtension(buf []byte) string {
	kind, unknown := filetype.Match(buf)
	if unknown != nil || kind.Extension == "unknown" {
		return ".txt"
	}
	return "." + kind.Extension
}
func DirExists(directory string) bool {
	_, err := os.Stat(directory)
	return (err == nil)
}

func MakeRandomDirectory(base_dir string, perm os.FileMode) string {
	randstr := randomString(4)
	directory := base_dir + "/" + randstr
	for DirExists(directory) {
		randstr = randomString(4)
		directory = base_dir + "/" + randstr // TODO read readString length from command line
	}
	err := os.Mkdir(directory, perm)
	if err != nil {
		logging.Errorf("failed to generate new file %s", err)
		panic("shit, I'm going to died now....")
	}
	return randstr
}

func readAndSave(conn *net.TCPConn, dir string, prefix string) {
	defer conn.Close()
	logging.Infof("Connection from %s\n", conn.RemoteAddr())
	data, error2 := ioutil.ReadAll(conn)
	if error2 != nil {
		logging.Errorf("Error when read data from socket: %s\n", error2)
		return
	}
	directory := MakeRandomDirectory(dir, 0777)
	extension := getFileExtension(data)
	fileName := directory + "/index" + extension
	f, error2 := os.Create(dir + "/" + fileName)
	logging.Infof("save %s for %s", fileName, conn.RemoteAddr())
	if error2 != nil {
		conn.Write([]byte("failed create file, please connect admin"))
		logging.Errorf("failed create file, %s", error2)
		return
	}
	f.Write(data)
	f.Close()
	if extension == ".txt" {
		conn.Write([]byte(prefix + directory + "\n")) // nginx will auto read index.txt when visit the directory. short url is more friendly.
	} else {
		conn.Write([]byte(prefix + fileName + "\n"))
	}
}

func main() {
	port := flag.String("port", "9999", "http listen port")
	host := flag.String("host", "127.0.0.1", "bind ip")
	dir := flag.String("dir", "/tmp/", "directory")
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
