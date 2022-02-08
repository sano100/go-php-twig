package twig

var TwigViewBox map[int]*TwigView

var ConfPhpFile = "php"
var ConfPhpDir string
var ConfTwigDir string
var ConfCacheDir string
var ConfHeader = "[[go_php_piper_start]]"
var ConfFooter = "[[go_php_piper_finish]]"
var ConfDev = true
var ConfNumPipes = 10 //开启几个php进行渲染工作,并不是多多益善，自己测试性能决定开几个
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
