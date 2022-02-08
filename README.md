# go-php-twig
允许你在golang项目中使用php的twig模板引擎，享受php的强大吧！哈哈哈！

### 注意
1.首先你需要准备php环境和composer环境

2.获取项目

`go get github.com/sano100/go-php-twig`

###安装使用
`

 	//1.设置配置参数
 	//--因为配置过环境变量这里使用php,否则请使用绝对路径
 	twig.ConfPhpFile = "php"
 	
 	//--配置生成index.php脚本以及composer.json文件的位置
 	twig.ConfPhpDir = "D:/go_www/test/php"
 	
 	//--配置模板缓存的位置
 	twig.ConfCacheDir = "D:/go_www/test/var/tmp"
 	
 	//--配置模板文件*.twig的根目录
 	twig.ConfTwigDir = "D:/go_www/test/views"
 	
 	//--配置是否为开发模式，
 	//----true:开发模式下仅开启1个php管道进程，twig模板文件修改后实时生效
 	//----false:生产环境下，根据ConfNumPipes启用多个php线程，并且模板使用缓存，
 	//----修改twig文件后需要先清空缓存才能生效,配合opcache可能有更佳性能表现
 	twig.ConfDev = true
 	
 	//--配置开启几个php线程进行模板渲染，
 	//--todo 因为是随机分配，低并发情况可能会出现都挤在同一个线程等待的情况
 	twig.ConfNumPipes = 3
 	
 	// go 与 php 通信是，总是以ConfHeader开头，以ConfFooter结束，
 	// 没有特殊原因无需配置，有默认值，ConfHeader和ConfFooter应该由良好的辨识度，且不能相同
 	//twig.ConfHeader = "[[go_php_piper_start]]"
 	//twig.ConfFooter = "[[go_php_piper_finish]]"
 
 	//创建运行所需的php代码
 	//--需要提前准备php环境和composer环境，
 	//--即 php -v 和 composer i 命令均能正常运行；
 	twig.CheckPHPFileOrCreate()//此方法安装后运行一次即可删除
 	//--CheckPHPFileOrCreate命令结束后还应切换到ConfPhpDir配置的目录执行composer i命令安装相关依赖
 
 	//初始化的部分，程序运行生命周期只运行一次
 	twig.Init()
 
 	//模板数据渲染
 	var data = make(map[string]string)
 	data["name"] = "this is name"
 	//模板路径说明：基于ConfTwigPath，“/”表示根目录，请使用绝对路径，相对路径会出错
 	//--使用模板嵌套时也一样{% extends "/main.twig" %}，请使用绝对路径
 	rst := twig.Render("/index.twig", data)
 	println(rst)
 `