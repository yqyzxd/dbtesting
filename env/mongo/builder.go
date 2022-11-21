package mongo

import (
	"fmt"
	"github.com/yqyzxd/dbtesting"
)

const (
	mongoConnStr string = "mongodb://%s:%s@%s:%s/?authSource=admin&readPreference=primary&ssl=false"
)

type builder struct {
	config *dbtesting.Config
}

func NewBuilder(config *dbtesting.Config) *builder {
	return &builder{
		config: config,
	}
}

func (b *builder) BuildEnv() ([]string, error) {
	return nil, nil
}

func (b *builder) BuildURI(hostIP, hostPort string) (string, error) {
	return fmt.Sprintf(mongoConnStr, b.config.User, b.config.Password, hostIP, hostPort), nil
}
