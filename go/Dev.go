package twig

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func RenderDev(tpl string, data interface{}) string {
	param := []string{
		Root + `/ext/twig/lib/php/index.php`, //脚本文件位置
		"--tplDir=" + Root + "/views/twig/",  //传递给php用于查找模板目录
		"--dev=1",                            //开发模式不实用缓存，修改后无需重新编译，生产环境需要清除缓存
		"--cacheDir=" + ConfCacheDir,
		"--header=" + ConfHeader, //stdin,stdout 值传递的开头和结尾
		"--footer=" + ConfFooter, //stdin,stdout 值传递的开头和结尾
	}
	cmd := exec.Command("php", param...)
	var config = make(map[string]interface{})
	config["__tplPath"] = tpl
	config["data"] = data
	content, err := json.Marshal(config)
	if err != nil {
		println("json error")
	}
	//判断输入内容,依靠特定的结尾来判断本次输入是否完成
	startFlag := ConfHeader
	endFlag := ConfFooter
	tmp := startFlag + string(content) + endFlag
	stdin := bytes.NewBuffer([]byte(tmp))
	cmd.Stdin = stdin

	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out[len(startFlag) : len(out)-len(endFlag)-1])
}
