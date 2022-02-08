package twig

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"
	"os/exec"
	"sync"
)

type TwigView struct {
	Cmd             *exec.Cmd
	OutPutReader    *bufio.Reader
	InputPipeReader *io.PipeReader
	InputPipeWriter *io.PipeWriter
	mt              *sync.Mutex
	index           int
}

func (this *TwigView) Init(index int) *TwigView {
	param := []string{
		Root + `/ext/twig/lib/php/index.php`, //脚本文件位置
		"--tplDir=" + Root + "/views/twig/",  //传递给php用于查找模板目录
		"--dev=1",                            //开发模式不实用缓存，修改后无需重新编译，生产环境需要清除缓存
		"--cacheDir=" + ConfCacheDir,
		"--header=" + ConfHeader, //stdin,stdout 值传递的开头和结尾
		"--footer=" + ConfFooter, //stdin,stdout 值传递的开头和结尾
	}
	cmd := exec.Command("php", param...)
	pr, pw := io.Pipe()
	cmd.Stdin = pr
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		println("Cmd StdoutPipe Error :" + err.Error())
	}
	if outPipe == nil {
		println("Cmd StdoutPipe is nil")
	}
	reader := bufio.NewReader(outPipe)
	err = cmd.Start()
	if err != nil {
		println("Cmd Start Error :" + err.Error())
	}
	tv := &TwigView{cmd, reader, pr, pw, &sync.Mutex{}, index}
	//println("php pipe start " , index)
	return tv
}

func RenderPipe(tpl string, data interface{}) string {
	if ConfDev {
		return RenderDev(tpl, data)
	}

	index := rand.Intn(ConfNumPipes)
	tv := TwigViewBox[index]
	tv.mt.Lock()
	defer tv.mt.Unlock()
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
	_, err = tv.InputPipeWriter.Write([]byte(tmp))
	if err != nil {
		println("Cmd StdinPipe Write Error :" + err.Error())
	}
	rst := ""
	for {
		s, err := tv.OutPutReader.ReadString('\n')
		if len(s) > len(startFlag) && s[0:len(startFlag)] == startFlag {
			rst = ""
		}
		rst += s
		if err != nil {
			println("err", err.Error())
			break
		}
		if err == io.EOF {
			println("== io.EOF")
			break
		}
		if len(s) > len(endFlag) && s[len(s)-len(endFlag)-1:len(s)-1] == endFlag {
			return rst[len(startFlag) : len(rst)-len(endFlag)-1]
		}
	}
	return rst
}
