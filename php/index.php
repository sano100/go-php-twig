<?php
$options = getopt('', array("tplDir::", "dev::", "cacheDir::", "header::", "footer::"));

$cacheDir   = $options["cacheDir"];
$tplDir     = $options["tplDir"];
$debug      = $options["dev"]       ? true : false;
$startFlag  = $options["header"]    ? $options["header"] : "[[go_php_piper_start]]";
$endFlag    = $options["footer"]    ? $options["footer"] : "[[go_php_piper_finish]]";

require_once __DIR__ . "/vendor/autoload.php";
$loader = new \Twig\Loader\FilesystemLoader($tplDir);
if ($debug){
    $opt = array("debug" => true, "auto_reload" => true, "cache" => false);
}else{
    $opt = array("debug" => false, "auto_reload" => false, "cache" => $cacheDir);
}
$twig = new Twig_Environment($loader, $opt);
$input = "";
while ($stdin = fread(STDIN, 5120)) {
    if (substr($stdin, 0, strlen($startFlag)) == $startFlag) {
        $input = "";
    }
    $input .= $stdin;
    if (substr($stdin, strlen($endFlag) * -1) != $endFlag) {
        continue;
    }
    $len = strlen($input);
    //转为php数组
    try {
        $conf = json_decode(substr($input, strlen($startFlag), strlen($endFlag) * -1), true);
    } catch (Exception $exception) {
        fwrite(STDOUT, $startFlag.$exception->getMessage() .$endFlag."\n");
        $input = "";
        continue;
    }
    $template = $twig->load($conf["__tplPath"]);
    $content = $template->render($conf["data"]);
    //STDOUT 返回渲染后的数据
    fwrite(STDOUT, $startFlag.$content .$endFlag."\n");
}