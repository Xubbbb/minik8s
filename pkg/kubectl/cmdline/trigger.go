package cmdline

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/server"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

var MqConn *minik8s_message.MQConnection

var TriggerFuncTable = map[string]func(string, string) error{
	"Function": triggerFunction,
	"Workflow": triggerWorkflow,
}

func TriggerCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		cmd.Usage()
		return
	}

	// 先获取kind
	kind := args[0]
	triggerFunc, ok := TriggerFuncTable[kind]
	if !ok {
		fmt.Println("kind not supported")
		return
	}

	string1 := args[1]
	string2 := args[2]

	err := triggerFunc(string1, string2)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	// fmt.Println("result is ", result)
}

func triggerFunction(functionName string, paramFile string) error {
	// 先把参数从文件里读出来
	bytes, err := os.ReadFile(paramFile)
	if err != nil {
		fmt.Println("failed to read bytes from file")
		return err
	}

	paramStr := string(bytes)

	// 构建请求
	var desc struct {
		FunctionName string `json:"functionName"`
		Params       string `json:"params"`
	}
	desc.FunctionName = functionName
	desc.Params = paramStr
	jsonData, _ := json.Marshal(desc)
	result, err := httpRequest.PostRequest(server.RootURL+server.TriggerServerlessFunction, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result is", result)
	return nil
}

func triggerWorkflow(workflowFile string, paramFile string) error {
	// 先把参数从文件里读出来
	workfileFileContent, err := os.ReadFile(workflowFile)
	if err != nil {
		fmt.Println("failed to read bytes from workflowFile")
		return err
	}
	paramFileContent, err := os.ReadFile(paramFile)
	if err != nil {
		fmt.Println("failed to read bytes from paramFile")
		return err
	}

	// 构建请求
	var desc struct {
		WorkfileFileContent string `json:"workflowFileContent"`
		ParamFileContent    string `json:"paramFileContent"`
		QueueName           string `json:"queueName"`
	}
	desc.WorkfileFileContent = string(workfileFileContent)
	desc.ParamFileContent = string(paramFileContent)
	jsonData, _ := json.Marshal(desc)
	// 申请一个queue，一块发过去
	// 既然是一个临时的就直接UUID前八位作为队列名了
	uuid, _ := idgenerate.GenerateID()
	desc.QueueName = uuid[:8]

	ch, err := MqConn.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return err
	}
	q, err := ch.QueueDeclare(desc.QueueName, true, true, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	done := make(chan bool)
	result, err := httpRequest.PostRequest(server.RootURL+server.TriggerServerlessWorkflow, jsonData)

	if err != nil {
		fmt.Println(err)
		return err
	}
	var wg sync.WaitGroup
	go MqConn.Subscribe(q.Name, func(d amqp.Delivery) {
		fmt.Println(string(d.Body))
		// TODO 判断是否停止监听
	}, done)
	wg.Wait()

	ch.QueueDelete(q.Name, false, false, false)

	fmt.Println("result is", result)
	return nil
}
