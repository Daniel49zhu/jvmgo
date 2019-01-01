- java命令

java命令有四种形式
```
java [-options] class [args]
java [-options] -jar jarfile [args]
javaw [-options] class [args]
javaw [-options] -jar jarfile [args]
```
可以向java命令传递三组参数：选项，主类名（或者JAR文件名）和main方法参数。
选项由`-`开头。通常，第一个非选项参数给出主类的完全限定名（fully qualified class name）。
但是如果用户提供了–jar选项，则第一个非选项参数表示JAR文件名，java命令必须从这个JAR文件中寻找主类。

选项可以分为两类，标准选项和非标准选项。标准选项较为稳定不会轻易变动，非标准选项以-X开头，另一部分
较高级的非标准选项以-XX开头
![选项](images/option.jpg "选项")
详细可以参考https://docs.oracle.com/javase/8/docs/technotes/tools/windows/java.html

[cmd.go](cmd.go)

go提供的os包里定义了一个Args便俩个，其存放了传递给命令行参数的全部参数，直接处理Args变量会需要大量处理，
而flag包可以保暖关注我们处理命令行选项。parseCmd方法处理了-?/-help,-cp/-classpath,-version三种选项，并将之后所有参数
封装进了cmd.class和cmd.args两个参数中

[main.go](main.go)

main函数显示调用ParseCommand()函数来解析命令行参数
```
 .\ch01.exe -version
 ouput:version 0.0.1
 
.\ch01.exe -cp foo/bar MyApp arg1 arg2
output:classpath:foo/bar class:MyApp args:[arg1 arg2]
  
.\ch01.exe -help
output:Usage: D:\GoLang\GoProject\bin\ch01.exe [-optipns] class [args...]
```

- 小结

第一章介绍了java命令的基本用法，编写了简单的命令行工具来处理-?/-help,-cp/-classpath,-version三种选项，
第二章将深入了解-classpath选项，来了解java虚拟机是如何加载class文件的
