package service

import (
	"io"
	"os"
	"fmt"
    "time"
    "bytes"
    "strings"
    "strconv"
	"syscall"
	"os/exec"
	"net/http"
    "io/ioutil"
	"os/signal"
    "path/filepath"
	"github.com/spf13/viper"
	"github.com/erikdubbelboer/gspt"
	log "github.com/sirupsen/logrus"
    "github.com/jingwu15/composer_pack/lib/json"
	gracehttp "github.com/jingwu15/composer_pack/lib/gracehttp"
	logchan "yundun/dispatcher_app_cli/lib/logchan"
)

var server *gracehttp.Server
var procTitleDefault string = "composer_pack_server"

func findProcess(procTitle string) ([]int, error) {
	var err error
	matches, err := filepath.Glob("/proc/*/cmdline")
	if err != nil {
		return nil, err
	}
	var pid int
	var pids = []int{}
	var tmp []string
	var body []byte
	for _, filename := range matches {
		body, err = ioutil.ReadFile(filename)
		if err == nil {
			if bytes.HasPrefix(body, []byte(procTitle)) {
				tmp = strings.Split(filename, "/")
				pid, _ = strconv.Atoi(tmp[2])
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}

//初始化日志
func InitLog() {
	go logchan.LogWrite()
	log.SetFormatter(&log.JSONFormatter{})
	//log.SetOutput(ioutil.Discard)
	log.SetLevel(log.DebugLevel)
	config := map[string]string{
		"error":      viper.GetString("log_error"),
		"info":       viper.GetString("log_info"),
		"writeDelay": "1",
		"cutType":    "day",
	}
	logChanHook := logchan.NewLogChanHook(config)
	log.AddHook(&logChanHook)
}

//处理信号，停止及重启
func handleSignals() {
	var sig os.Signal
	var signalChan = make(chan os.Signal, 100)
	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)
	for {
		sig = <-signalChan
		switch sig {
		case syscall.SIGTERM:
			log.Info("stop")
			server.Down()
			logchan.LogClose()
		case syscall.SIGUSR2:
			log.Info("restart")
			server.Restart()
			logchan.LogClose()
		default:
		}
	}
}

func Api_Up(w http.ResponseWriter, r *http.Request) {
    var err error
    reader, err := r.MultipartReader()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    packdir := viper.GetString("server.packdir")
    var keyname, filename string
    resps := map[string]map[string]string{}
    for {
        part, err := reader.NextPart()
        if err == io.EOF {
            break
        }
        resp := map[string]string{
            "code": "0",
            "msg": "failed",
            "filename": "",
        }
        keyname = part.FormName()
        filename = part.FileName()
        resp["filename"] = filename
        dst, err := os.Create(packdir + "/" + filename)
        if err != nil {
            resps[keyname] = resp
            continue
        }
        defer dst.Close()
        io.Copy(dst, part)

        resp["code"] = "1"
        resp["msg"] = "ok"
        resps[keyname] = resp
    }
    output, err := json.Encode(resps)
	fmt.Fprintln(w, string(output))
	return
}

func Run() {
	procTitle := procTitleDefault
	gspt.SetProcTitle(procTitle)
    fmt.Println(procTitle)

	go handleSignals()
	InitLog()

	http.HandleFunc("/up",     Api_Up)
    http.Handle("/down/", http.StripPrefix("/down/", http.FileServer(http.Dir(viper.GetString("server.packdir")))))
	hs := &http.Server{
		Addr:           viper.GetString("server.listen"),
		ReadTimeout:    time.Second * time.Duration(viper.GetInt("server.http_read_timeout")),
		WriteTimeout:   time.Second * time.Duration(viper.GetInt("server.http_write_timeout")),
		MaxHeaderBytes: 1 << 20,
	}
	hs.SetKeepAlivesEnabled(false)
    //err := hs.ListenAndServe()
    //fmt.Println(err)

	////改用 支持平滑重启的grace，因为涉及到信号一定要分为两行来写
	server = gracehttp.NewServer(hs)
    server.ListenAndServe()
}

func Start() {
	var err error
	cmdFile, _ := filepath.Abs(os.Args[0])
	cmd := "nohup " + cmdFile + " " + os.Args[1] + " run 2> /dev/null 1>/dev/null &"
	client := exec.Command("sh", "-c", cmd)
	err = client.Start()
	if err != nil {
		fmt.Println("composer_pack_server start error:")
		fmt.Println(err)
		return
	}
	err = client.Wait()
	if err != nil {
		fmt.Println("composer_pack_server start error:")
		fmt.Println(err)
		return
	}
	fmt.Println("composer_pack_server is started")
	return
}

func Restart() {
	procTitle := procTitleDefault
	pids, err := findProcess(procTitle)
	if err != nil {
		log.Error(err)
		return
	}
	for _, pid := range pids {
		syscall.Kill(pid, syscall.SIGUSR2)
	}
	fmt.Println("composer_pack_server is restarted")
}

func Stop() {
	procTitle := procTitleDefault
	pids, err := findProcess(procTitle)
	if err != nil {
		log.Error(err)
		return
	}
	for _, pid := range pids {
		syscall.Kill(pid, syscall.SIGTERM)
	}
	fmt.Println("composer_pack_server is stoped")
}
