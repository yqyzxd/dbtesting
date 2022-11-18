/*
Package dbtesting is used for unit test which gives you a clean db environment
for mysql usage:
	func TestMain(m *testing.M) {
		os.Exit(dbtesting.RunDBInDocker(m, &dbtesting.Config{
			Image:         "mysql:5.6",
			User:          "root",
			Password:      "123456",
			DBName:        "mu_test_db",
			DB:            dbtesting.Mysql,
			ContainerPort: "3306/tcp",
		}))
	}
	func NewEngine() (*xorm.Engine, error) {
		if dbtesting.ConnURI == "" {
			return nil, fmt.Errorf("conn uri is nil")
		}
		return xorm.NewEngine("mysql", dbtesting.ConnURI)
	}


for mongo usage:
		func TestMain(m *testing.M) {
				os.Exit(dbtesting.RunDBInDocker(m, &dbtesting.Config{
					Image:         "mongo",
					User:          "admin",
					Password:      "123456",
					DB:            dbtesting.Mongo,
					ContainerPort: "27017/tcp",
				}))
			}
		func NewClient(c context.Context) (*mongo.Client, error) {
			if dbtesting.ConnURI == "" {
				return nil, fmt.Errorf("conn uri is nil")
			}
			return mongo.Connect(c, options.Client().ApplyURI(dbtesting.ConnURI))
		}
*/
package dbtesting

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/yqyzxd/dbtesting/waiter"
	"strings"
	"testing"
)

var ConnURI string

type Config struct {
	Image         string
	DBName        string
	User          string
	Password      string
	ContainerPort string
	DB            DB
}
type DB int

const (
	Mysql DB = iota
	Mongo
)
const (
	mysqlConnStr string = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	mongoConnStr string = "mongodb://%s:%s@%s:%s/?authSource=admin&readPreference=primary&ssl=false"
)

// RunDBInDocker Docker中运行测试数据库
func RunDBInDocker(m *testing.M, config *Config) int {

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	context := context.Background()
	envMap := map[string]string{
		"MYSQL_ROOT_PASSWORD": config.Password,
		"MYSQL_DATABASE":      config.DBName,
	}
	var env []string
	for envKey, envVar := range envMap {
		env = append(env, envKey+"="+envVar)
	}
	containerBody, err := cli.ContainerCreate(context, &container.Config{
		Image: config.Image,
		ExposedPorts: nat.PortSet{
			nat.Port(config.ContainerPort): {},
		},
		Env: env,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{ //将容器端口 映射到以下的系统端口
			nat.Port(config.ContainerPort): []nat.PortBinding{
				{
					HostIP:   "127.0.0.1", //只接受本地请求，如果是0.0.0.0则是接收所有请求
					HostPort: "0",         //27018这里的端口写0 的话是会自动寻找空闲端口。写固定端口那就是指定的端口
				},
			},
		},
	}, nil, nil, "")

	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(context, containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	defer func() {
		cli.ContainerRemove(context, containerBody.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("container removed")
	}()

	inspectJson, err := cli.ContainerInspect(context, containerBody.ID)
	if err != nil {
		panic(err)
	}
	hostPortBinding := inspectJson.NetworkSettings.Ports[nat.Port(config.ContainerPort)][0]
	switch config.DB {
	case Mysql:
		ConnURI = fmt.Sprintf(mysqlConnStr, config.User, config.Password, hostPortBinding.HostIP, hostPortBinding.HostPort, config.DBName)
	case Mongo:
		ConnURI = fmt.Sprintf(mongoConnStr, config.User, config.Password, hostPortBinding.HostIP, hostPortBinding.HostPort)
	}
	port := strings.ReplaceAll(config.ContainerPort, "/tcp", "")
	waiter.ForLog(port, cli, containerBody.ID).Wait(context)

	fmt.Printf("listening at %+v\n", hostPortBinding)

	fmt.Println("container started")

	return m.Run()
}
