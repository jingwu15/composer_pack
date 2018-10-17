package service

import (
	"io"
	"os"
    "fmt"
    "path"
    "bytes"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "mime/multipart"
	"github.com/spf13/viper"
)

func up(api string, filepath string) string {
    bodyBuffer := &bytes.Buffer{}
    bodyWriter := multipart.NewWriter(bodyBuffer)

    fileWriter, _ := bodyWriter.CreateFormFile("gzfile", path.Base(filepath))

    file, _ := os.Open(filepath)
    defer file.Close()

    io.Copy(fileWriter, file)

    contentType := bodyWriter.FormDataContentType()
    bodyWriter.Close()

    resp, _ := http.Post(api, contentType, bodyBuffer)
    defer resp.Body.Close()

    resp_body, _ := ioutil.ReadAll(resp.Body)

    return string(resp_body)
}

func Up() {
	packdir := viper.GetString("packdir")
    filenames, err := filepath.Glob(packdir + "/*.tar.gz")
    if err != nil {
        fmt.Println("扫描文件失败！")
    }

    flag := 0
    for _,filename := range filenames {
        result := up(viper.GetString("api"), filename)
        fmt.Println(filename, result)
        flag = 1
    }
    if flag == 0 {
        fmt.Println("未找文件！")
    }
}

