package pprof

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"os/exec"
	"runtime/pprof"
	"strconv"
	"strings"
)

type ProfileType string

//ProfileContext pprof相关，比如配置等
type ProfileContext struct {
	exposeType ProfileType
	httpPort   int
}

const (
	TypeNoExpose   ProfileType = ""
	TypeExposeFile ProfileType = "file"
	TypeExposeHttp ProfileType = "http"
	TypeExposeSvg  ProfileType = "svg"
)

//custom profile
var libProfile *pprof.Profile
var context ProfileContext

//Init 初始化pprof配置，启动http服务
func Init(profName string, exposeType string, httpPort int) {
	context.exposeType = ProfileType(exposeType)
	context.httpPort = httpPort
	if context.exposeType == TypeNoExpose {
		return
	}
	//profName := "my_experiment_thing"
	libProfile = pprof.Lookup(profName)
	if libProfile == nil {
		libProfile = pprof.NewProfile(profName)
	}

	if context.exposeType == TypeExposeFile {
		ff, err := os.Create(profName + ".pprof")
		if err != nil {
			log.Fatal(err)
		}
		libProfile.WriteTo(ff, 1)

		var cpuProfFile = "cpuProfFile"
		if cpuProfFile != "" {
			f, err := os.Create(cpuProfFile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	if context.exposeType == TypeExposeHttp {
		serverMux := mux.NewRouter()
		startServer(serverMux, profName)
	}

	if context.exposeType == TypeExposeSvg {
		serverMux := mux.NewRouter()
		serverMux.HandleFunc("/debug/pprofsvg/", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/block", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/goroutine", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/heap", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/mutex", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/threadcreate", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/cpuprofile", http.HandlerFunc(svgPprof))
		serverMux.HandleFunc("/debug/pprofsvg/"+profName, http.HandlerFunc(svgPprof))
		startServer(serverMux, profName)
	}
}

func startServer(serverMux *mux.Router, profName string) {
	serverMux.HandleFunc("/debug/pprof/", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/block", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/goroutine", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/heap", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/mutex", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/threadcreate", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/cpuprofile", http.HandlerFunc(httppprof.Index))
	serverMux.HandleFunc("/debug/pprof/"+profName, http.HandlerFunc(httppprof.Index))

	go func() {
		logEntry := logrus.WithFields(logrus.Fields{
			"type": context.exposeType,
			"port": context.httpPort,
		})
		logEntry.Infoln("pprof http started")
		if err := http.ListenAndServe(":"+strconv.Itoa(context.httpPort), serverMux); err != nil {
			log.Fatal("failed to start stress server", err)
		}
	}()
}

//对svg路径的请求做一次包装，输出通过命令返回的svg图片
func svgPprof(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	k := req.Form["debug"]
	if len(k) > 0 && k[0] == "1" {
		fmt.Println("test")
		out := execPprof(req.URL.Path)
		res.Write(out)
	} else {
		httppprof.Index(res, req)
	}
}

func execPprof(url string) []byte {
	url = strings.Replace(url, "pprofsvg", "pprof", -1)
	subProcess := exec.Command("go", "tool", "pprof", "-svg", "http://localhost:"+strconv.Itoa(context.httpPort)+""+url+"?debug=")
	out, err := subProcess.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	str := string(out)
	index := strings.Index(str, "<") //获取svg图片部分
	if index >= 0 {
		str = str[index:]
	} else {
		str = ""
	}
	return []byte(str)
}

func AddClient(clientID uint64) {
	if libProfile != nil {
		libProfile.Add(clientID, 1)
		updateFile()
	}
}
func RemoveClient(clientID uint64) {
	if libProfile != nil {
		libProfile.Remove(clientID)
		updateFile()
	}
}

func updateFile() {
	if context.exposeType == TypeExposeFile {
		ff, err := os.OpenFile(libProfile.Name()+".pprof", os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		libProfile.WriteTo(ff, 1)
	}
}
