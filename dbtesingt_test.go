package dbtesting

import (
	"github.com/yqyzxd/dbtesting/config"
	"os"
	"testing"
)

func TestRunDBInDocker(t *testing.T) {

}

func TestMain(m *testing.M) {
	os.Exit(RunDBInDocker(m, &config.Config{
		DB: config.Mysql,
	}))
}
