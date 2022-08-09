package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const Version = "V 0.1.0"

type Config struct {
	Namespace    string
	Group        string
	DataId       string
	Password     string
	Username     string
	Path         string
	FileName     string
	Ip           string
	Port         uint64
	PrintVersion bool
}

var config Config
var wg sync.WaitGroup

func init() {
	flag.StringVar(&config.Path, "path", "", "Local Storage Path")
	flag.StringVar(&config.FileName, "filename", "", "Local Storage filename")
	flag.StringVar(&config.Username, "username", "nacos", "the username to authenticate as")
	flag.StringVar(&config.Password, "password", "nacos", "the password to authenticate with")
	flag.StringVar(&config.DataId, "dataId", "", "the dataId in nacos")
	flag.StringVar(&config.Group, "group", "DEFAULT_GROUP", "the group in nacos")
	flag.StringVar(&config.Namespace, "namespace", "", "the namespace in nacos")
	flag.StringVar(&config.Ip, "ip", "127.0.0.1", "the ip in nacos")
	flag.Uint64Var(&config.Port, "port", 8848, "the port in nacos")
	flag.BoolVar(&config.PrintVersion, "version", false, "print version and exit")
}

func watch(client config_client.IConfigClient) {
	//Listen config change,key=dataId+group+namespaceId.
	err := client.ListenConfig(vo.ConfigParam{
		DataId: config.DataId,
		Group:  config.Group,
		OnChange: func(namespace, group, dataId, data string) {
			createFile(data)
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})
	if err != nil {
		panic(err)
	}
}

func getConfig(client config_client.IConfigClient) (content string, err error) {
	content, err = client.GetConfig(vo.ConfigParam{
		DataId: config.DataId,
		Group:  config.Group,
	})
	return
}

func createFile(content string) {
	filename := config.Path + config.FileName
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		// 创建文件失败处理
		panic(err)
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		// 写入失败处理
		panic(err)
	}
}

func main() {
	flag.Parse()

	if config.PrintVersion {
		fmt.Printf("go-nacos %s (Go Version: %s)\n", Version, runtime.Version())
		os.Exit(0)
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(
			config.Ip,
			config.Port,
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos")),
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(config.Namespace),
		constant.WithTimeoutMs(10000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername(config.Username),
		constant.WithPassword(config.Password),
		constant.WithLogLevel("warn"),
	)

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	content, err := getConfig(client)
	if err != nil {
		panic(err)
	}
	createFile(content)

	wg.Add(1)
	go watch(client)
	wg.Wait()
}
