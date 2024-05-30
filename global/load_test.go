package global_test

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/settings"
	"testing"

	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/test"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	fmt.Printf("%#v\n", global.C)
	var m global.ConfigMode = "file"
	defaultCf := "config/config.yaml"
	err := settings.LoadConfig(m, defaultCf)
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", global.C)
}

func TestLoadConfigFromEnv(t *testing.T) {
	should := assert.New(t)
	test.EnvLoad()
	var m global.ConfigMode = "env"
	fmt.Printf("%#v\n", global.C)
	defaultCf := "config/config.yaml"
	err := settings.LoadConfig(m, defaultCf)
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", global.C)
	should.Equal("8080", global.C.App.GrpcPort)
}
