### App调度 Cli

```
crontabd 用于生成计划任务, 并写入消息队列
```

#### 1. 安装
```
# 1. 安装信赖包
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/Shopify/sarama
go get github.com/buger/jsonparser
go get github.com/garyburd/redigo/redis
go get github.com/json-iterator/go
go get github.com/sirupsen/logrus
go get github.com/erikdubbelboer/gspt
go get github.com/oschwald/maxminddb-golang
go get github.com/jinzhu/gorm
go get github.com/beanstalkd/go-beanstalk
或
sh ./install.sh

# 2. 构建项目
go build -o ./bin/composer_pack github.com/jingwu15/composer_pack/composer_pack
```

#### 2. 使用
```
[jingwu@local composer_pack]$ ./bin/composer_pack
pack composer file and up to file server

Usage:
  ./composer_pack [command]

Available Commands:
  check       check the composer update
  help        Help about any command
  pack        pack composer.json,composer.lock,vendor to md5.tar.gz
  server      the file server
  up          up the *.tar.gz file to the file server
  version     composer_pack version

Flags:
  -c, --config string   JSON format configuration file (default "./composer_pack.json")
  -h, --help            help for ./composer_pack

Use "./composer_pack [command] --help" for more information about a command.
```

#### 3. 测试上传
./bin/composer_pack server start

curl --form "gzfile=@./composer_pack.json" http://127.0.0.1:8083/up

