package main

import (
	"flag"
	"net/http"
	"os/exec"
	"os"
	"time"
	"bytes"
	"errors"
	"io/ioutil"
)

var (
	host = flag.String("host", ":1317", "监听地址")
)

func init() {
	flag.Parse()
}

func download(pkgname string) (string, error) {
	buf := bytes.NewBufferString("")
	tmpPath := time.Now().Format("20060102150405")
	cmd := exec.Command("go", "get", "-v", "-u", pkgname)
	cmd.Env = []string{"GOPATH=/tmp/" + tmpPath, "GOROOT=" + os.Getenv("GOROOT"), "PATH=" + os.Getenv("PATH")}
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		return "", errors.New(buf.String())
	}

	fname := "/tmp/" + tmpPath + ".tar.gz"
	buf.Reset()
	cmd = exec.Command("tar", "-vczf", fname, "/tmp/" + tmpPath + "/")
	cmd.Stdout = buf
	cmd.Stderr = buf
	err = cmd.Run()
	if err != nil {
		return "", errors.New(buf.String())
	}

	return fname, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pkname := r.URL.Query().Get("name")
		if pkname == "" {
			w.WriteHeader(500)
			w.Write([]byte("包名为空"))
			return
		}

		fname, err := download(pkname)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		content, err := ioutil.ReadFile(fname)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(content)

	})
	http.ListenAndServe(*host, nil)
}
