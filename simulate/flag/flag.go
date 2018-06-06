package flag

import (
	"flag"
	"log"
	"os"
)

// Flags 测试程序的 flag
var Flags struct {
	DBPath string
}

func init() {
	gopath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Panicln("GOPATH 未设置")
	}
	flag.StringVar(&Flags.DBPath, "dbpath", gopath, "db path")
	flag.Parse()
}
