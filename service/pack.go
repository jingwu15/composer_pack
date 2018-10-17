package service

import (
	"io"
	"os"
    "fmt"
    "strings"
    "crypto/md5"
	"archive/tar"
	"compress/gzip"
	"github.com/spf13/viper"
)

//压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, file := range files {
		err := compress(file, "", tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, tw *tar.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压 tar.gz
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func Pack() {
    var err error
    viper.Get("pros")
	pros := viper.GetStringSlice("pros")
	packdir := viper.GetString("packdir")

    var md5text, filename string
    var file, md5file *os.File
    for _,pro := range pros {
        files := []*os.File{}
        filename = pro + "/composer.json"

        //md5
        md5file, err = os.Open(filename)
        if(err != nil) {
            fmt.Println(filename, "不存在！")
            os.Exit(-1)
        }
        finfo, _ := md5file.Stat()
        fsize := finfo.Size()
        raw := make([]byte, fsize)
        _, _ = md5file.Read(raw)
        md5text = fmt.Sprintf("%x", md5.Sum(raw))

        for _,name := range []string{"composer.json", "composer.lock", "vendor"} {
            filename = pro + "/" + name
            file, err = os.Open(filename)
            if(err != nil) {
                fmt.Println(filename, "不存在！")
                os.Exit(-1)
            }
            files = append(files, file)
        }
        err = Compress(files, packdir + "/" + md5text + ".tar.gz")
        if(err != nil) {
        }
    }
}

