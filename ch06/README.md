第4章初步实现了线程私有的运行时数据区，在此基础上第5章实现了一个简单的解释器和150多条指令。
这些指令主要是操作局部变量表和操作数栈，进行数学运算，比较运算和跳转控制等。本章将实现
线程共享的运行时数据区，包括方法区和运行时常量池。

第2章实现了类路径，可以找到class文件，并把数据加载到内存中。第3章实现了class文件解析，可以把
class数据解析成一个ClassFile结构体。本章将进一步处理ClassFile'结构体，把它加以转换，放进方法区
以供后续使用。本章还会初步讨论类和对象的设计，实现一个简单的类加载器，并且实现类和对象相关的部分指令。

- 方法区

    方法区(Method Area)是运行时数据取的一块逻辑区域，由多个线程共享。方法区主要存放从class
    文件获取的类信息。此外，类变量也存放在方法区中。当Java虚拟机第一次使用某个类时，它会搜索
    类路径，找到相应的class文件，然后读取并解析class文件，把相关信息放进方法区。至于方法区
    到底位于何处，是固定大小还是动态调整，是否参与垃圾回收，以及如何在方法区内存放类数据等，
    Java虚拟机规范并没有明确规定。需要放进方法区的信息有
    
    - 类信息
    
    使用结构体来表示将要放进方法区内的类。在rtda/heap下新建[class.go](rtda/heap/class.go)，在其中
    定义Class结构体，代码如下：
    ```
    type Class struct {
    	accessFlags       uint16
    	name              string // thisClassName
    	superClassName    string
    	interfaceNames    []string
    	constantPool      *ConstantPool
    	fields            []*Field
    	methods           []*Method
    	loader            *ClassLoader
    	superClass        *Class
    	interfaces        []*Class
    	instanceSlotCount uint
    	staticSlotCount   uint
    	staticVars        Slots
    }
    ```
    accessFalg是类的访问标志，总共16bit。字段和方法也有访问标志，但具体标志位的含义可能有所不同。
    根据Java虚拟机规范，把哥哥比特位的含义统一定义在heap/[access_flags.go](rtda/heap/access_flags.go)里。
    name、superClassName和interfaceNames字段分别存放类名、超类名和接口名。注意这些都是完全限定名，形式类似
    `java/lang/Object`的形式。constantPool字段存放运行时常量池指针，filed和methods字段分别村官方字段表和
    方法表。运行时常量池将在稍后详细介绍。