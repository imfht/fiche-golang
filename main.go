package main

import (
	"flag"
	"fmt"
	"github.com/couchbase/goutils/logging"
	"github.com/pkg/errors"
	"gopkg.in/h2non/filetype.v1"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const MaxUploadSize = 4096 * 1024
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
	return err == nil
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

func readAll(r io.Reader) (b []byte, err error) {
	buf := make([]byte, 0, MaxUploadSize) // buffer size=max size
	tmp := make([]byte, MaxUploadSize/8)  // using small tmo buffer for demonstrating
	exitByEof := false
	totalByte := 0
	for totalByte < MaxUploadSize {
		n, err := r.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			} else {
				exitByEof = true
			}
			break
		}
		//fmt.Println("got", n, "bytes.")
		totalByte += n
		buf = append(buf, tmp[:n]...)
	}
	if !exitByEof {
		return buf, errors.New("File Too Large!")
	}
	return buf, nil // only consider file too large Exception
}

func readAndSave(conn *net.TCPConn, dir string, prefix string) {
	defer conn.Close()
	logging.Infof("Connection from %s\n", conn.RemoteAddr())
	data, error2 := readAll(conn) // max 10M
	//ioutil.ReadAll()
	if error2 != nil {
		//logging.Errorf("Error when read data from socket: %s\n", error2)
		conn.Write([]byte("File Too Large. MAX 4096 kb \n")) // FIXME seems I repeat me again. error2 contains "File Too Large".
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
	conn.Write([]byte(prefix + fileName + "\n"))
}

var prefix string

func init() {
	prefix = os.Getenv("prefix")
	if len(prefix) == 0 {
		prefix = "http://127.0.0.1/"
	}
}
func main() {
	port := flag.String("port", "9999", "http listen port")
	host := flag.String("host", "0.0.0.0", "bind ip")
	flag.Parse()
	addr, error := net.ResolveTCPAddr("tcp", *host+":"+*port)
	if error != nil {
		fmt.Printf("Cannot parse \"%s\": %s\n", port, error)
		os.Exit(1)
	}
	listener, error := net.ListenTCP("tcp", addr)
	if error != nil {
		fmt.Printf("Cannot listen: %s\n", error)
		os.Exit(1)
	}
	logging.Infof("listen on %s:%s", *host, *port)

	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}
	var dir = "./data/"
	var fs = http.FileServer(http.Dir("./data"))
	go http.ListenAndServe("0.0.0.0:80", fs)
	for { // ever...
		conn, error := listener.AcceptTCP()
		if error != nil {
			logging.Errorf("Cannot accept: %s\n", error)
			os.Exit(1)
		}
		go readAndSave(conn, dir, prefix)
	}
}
