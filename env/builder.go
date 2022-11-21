package env

import (
	"github.com/yqyzxd/dbtesting/config"
	"github.com/yqyzxd/dbtesting/env/mongo"
	"github.com/yqyzxd/dbtesting/env/mysql"
)

type Builder struct {
	config.Builder
}

func NewBuilder(cfg *config.Config) *Builder {
	var envBuilder config.Builder

	switch cfg.DB {
	case config.Mysql:
		envBuilder = mysql.NewBuilder(cfg)
	case config.Mongo:
		envBuilder = mongo.NewBuilder(cfg)
	}
	return &Builder{
		envBuilder,
	}
}

func (b *Builder) BuildEnv() ([]string, error) {
	return b.Builder.BuildEnv()
}
func (b *Builder) BuildURI(host, port string) (string, error) {
	return b.Builder.BuildURI(host, port)
}
