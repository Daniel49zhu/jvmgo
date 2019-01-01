- 搜索class文件

    当我们在加载一个类之前，首先需要加载的就是它的超类（java.lang.Object），同时因为main方法中的参数数组，我们还需要准备好
java.lang.String和java.lang.String[]类等等，本章将讨论java虚拟机是从哪里寻找到这些类的。

- 类路径
    
    java的虚拟规范没有规定虚拟机如何寻找类，因此不同的虚拟机会采用不同的方案。Oracle的虚拟机是根据
    类路径（class path）来搜索类，按照搜索的先后顺序，类路径可以分为：
    - 启动类路径（bootstrap classpath）：默认为jre/lib目录，java标准库（rt.jar）位于该路径，可以通过-Xbootclasspath来指定启动类路径
    - 扩展类路径（extension classpath）：默认为jre/lib/extmulu，使用java扩展机制的类位于这个路径
    - 用户类路径（user classpath）：默认是当前目录，可以通过-classpath/-cp来指定，也能用来指定一个具体的jar文件或zip文件；
    可以通过分号来指定多个，JDK6开始允许使用`*`来指定一个目录下的所有jar文件
```
    java -cp path\to\classes ...
    java -cp path\to\lib1.jar ...
    java -cp path\to\lib2.zip ...
    java -cp path\to\classes;lib\a.jar;lib\b.jar;lib\c.zip ...
    java -cp classes;lib\* ...
```
    
