package main

import (
	"flag"
	"net/http"
	"github.com/mucolud/log"
	"os"
	"io"
	"io/ioutil"
)

var (
	shost = flag.String("shost", ":1317", "服务器地址")
	pkcname = flag.String("pkgname", "", "包名称")
)

func init() {
	flag.Parse()
}

func main() {
	resp, err := http.Get(*shost + "/?name=" + *pkcname)
	if err != nil {
		log.Fatal(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	f, err := os.Create("/tmp/pkg.tar.gz")
	if err != nil {
		log.Error(err)
		return
	}
	if resp.StatusCode == 500 {
		bs, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(bs))
	}

	io.Copy(f, resp.Body)
	f.Close()
	log.Info("下载ok,路径:", "/tmp/pkg.tar.gz")

}
