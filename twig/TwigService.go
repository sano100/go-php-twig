package twig

var TwigViewBox map[int]*TwigView

//--因为配置过环境变量这里使用php,否则请使用绝对路径
var ConfPhpFile = "php"

//--配置生成index.php脚本以及composer.json文件的位置
var ConfPhpDir string

//--配置模板文件*.twig的根目录
var ConfTwigDir string

//--配置模板缓存的位置
var ConfCacheDir string

//--配置是否为开发模式，
//----true:开发模式下仅开启1个php管道进程，twig模板文件修改后实时生效
//----false:生产环境下，根据ConfNumPipes启用多个php线程，并且模板使用缓存，
//----修改twig文件后需要先清空缓存才能生效,配合opcache可能有更佳性能表现
var ConfDev bool

// go 与 php 通信是，总是以ConfHeader开头，以ConfFooter修改，
// 没有特殊原因无需配置，有默认值，ConfHeader和ConfFooter应该由良好的辨识度，且不能相同
var ConfHeader = "[[go_php_piper_start]]"
var ConfFooter = "[[go_php_piper_finish]]"

//--配置开启几个php线程进行模板渲染，
//--todo 因为是随机分配，低并发情况可能会出现都挤在同一个线程等待的情况
//--开启几个php进行渲染工作,并不是多多益善，自己测试性能决定开几个
var ConfNumPipes = 10

func Init() {
	if ConfDev {
		return
	}
	TwigViewBox = make(map[int]*TwigView)
	for i := 0; i < ConfNumPipes; i++ {
		tv := new(TwigView)
		TwigViewBox[i] = tv.Init(i)
	}
}

//执行渲染
//--tpl:模板相对于模板根目录的相对路径，
//--data：模板数据
func Render(tpl string, data map[string]interface{}) string {
	//根据开发环境决定使用开发模式渲染还是生产模式
	if ConfDev {
		return RenderDev(tpl, data)
	} else {
		return RenderPipe(tpl, data)
	}
}
