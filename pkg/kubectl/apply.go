package cmdline

import (
	"fmt"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/ty"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var table = map[string]func(b []byte) error{
	"Pod":        applyPod,
	"Service":    applyService,
	"ReplicaSet": applyReplicaSet,
}

func ApplyCmdHandler(cmd *cobra.Command, args []string) {
	// 先看一下参数是不是文件路径
	result := checkFilePath(args)
	if !result {
		cmd.Usage()
		return
	}

	// 读取文件内容，先找到kind
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("failed to read yaml file")
		return
	}

	var tmp map[string]interface{}
	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		fmt.Println("failed to unmarshal yaml file")
		return
	}

	// kind不支持
	if tmp["kind"] == nil {
		fmt.Println("no kind field found")
		return
	}

	kind := tmp["kind"].(string)
	if table[kind] == nil {
		fmt.Println("kind not supported")
		return
	}

	// 根据kind跳转到相应的处理函数，相当于switch
	err = table[kind](data)
	if err != nil {
		fmt.Println(err)
	}
}

func applyPod(b []byte) error {
	var podDesc ty.PodDesc

	err := yaml.Unmarshal(b, &podDesc)
	if err != nil {
		fmt.Println("failed to unmarshal pod yaml")
		return err
	}

	// PostRequest()

	return nil
}

func applyService(b []byte) error {
	return nil
}

func applyReplicaSet(b []byte) error {
	return nil
}

func checkFilePath(args []string) bool {
	// 检查参数给出的文件路径是否存在

	if len(args) == 0 {
		return false
	}

	result, err := os.Stat(args[0])
	if err != nil {
		return false
	}

	if result.IsDir() {
		return false
	}

	return true
}
