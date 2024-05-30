package version_test

import (
	"fmt"
	"testing"

	"github.com/yanshicheng/ikubeops-gin-demo/version"
)

func TestVersion(t *testing.T) {
	version.IkubeopsGoVersion = "go1.22.3"
	version.IkubeopsCommit = "123456"
	version.IkubeopsBranch = "master"
	version.IkubeopsBuildTime = "2020-01-01 00:00:00"
	version.IkubeopsTag = "v1.1.1"
	fmt.Println(version.FullTagVersion())
	fmt.Println(version.ShortTagVersion())
}
