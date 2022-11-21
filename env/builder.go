package env

import (
	"github.com/yqyzxd/dbtesting"
	"github.com/yqyzxd/dbtesting/env/mongo"
	"github.com/yqyzxd/dbtesting/env/mysql"
)

type Builder struct {
	dbtesting.Builder
}

func NewBuilder(config *dbtesting.Config) *Builder {
	var envBuilder dbtesting.Builder
	switch config.DB {
	case dbtesting.Mysql:
		envBuilder = mysql.NewBuilder(config)
	case dbtesting.Mongo:
		envBuilder = mongo.NewBuilder(config)
	}
	return &Builder{
		envBuilder,
	}
}

func (b *Builder) BuildEnv() ([]string, error) {
	return b.BuildEnv()
}
func (b *Builder) BuildURI(host, port string) (string, error) {
	return b.BuildURI(host, port)
}
