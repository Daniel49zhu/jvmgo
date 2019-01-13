本章在前两章的基础上会是显示一个简单的解释器，并且实现大约150条指令，在之后会不断改进这个解释器，
让它可以执行更多的指令。

- 字节码和指令集

    Java虚拟机顾名思义没事一台虚拟的机器，而字节码(bytecode)就是运行在这台机器上的
    机器码。已知每个Java文件都会被便一场一个class文件，类或接口的方法信息就放在class文件的
    method_info结构中。如果方法不是抽象的也不是native方法，就会被编译器编译成字节码（即使方法是空的
    编译器也会生成一条return语句），存放在method_info的Code属性中。
    
    字节码中存放编译后的Java虚拟机指令。每条指令都以一个单字节的操作码（opcode）开头。由于只使用一字节
    表示字节码，因此Java虚拟机最后只能支持256（2^8）条指令，截止到第八版已经定义了205条指令，分别是
    0到202，254和255。这205条指令称为Java虚拟机的指令集(instruction set)。和汇编语言类似，每个
    字节码有一个助记符（mnemonic）。比如0x00助记符是nop。
    
    Java虚拟机使用的是变长指令，操作码后面可以跟零字节或多字节的操作数（operand）。如果把指令想象
    成函数的话，操作数就是它的参数。为了让编码后的字节码更加紧凑，很多操作吗本身就是隐含了
    操作数，比如把常数0推入操作数栈的指令是iconst_0。下面通过例子来观察虚拟机指令。
    
    ![getstatic](images/operand.jpg "getstatic指令")
    getstatic指令的操作码是0xB2，助记符是getstatic，操作数0x0002，代表常量池里的第二个常量。
    第四章中我们知道操作数栈和常量池只存放了数据的值，并不记录数据类型。所以指令必须知道自己在
    操作什么类型的数据。因此iadd指令就是对iint进行加法，dstore就是把操作数栈顶的double值弹放入巨变
    变量表中；areturn从方法中返回引用值。也就是说，如果某类指令可以操作不同类型的变量，
    则助记符的第一个字母表示变量类型。其对应关系如下所示。
    
    ![对应关系](images/relation.jpg "对应关系")
    
    Java虚拟机规范把205条指令按用途分成了11类，分别是：常量（constants）指令、加载（loads）指令、
    存储（stores）指令、操作数栈（stack）指令、数学（math）指令、转换（conversions）指令、比较（comparisions）指令
    、控制（control）指令、引用（refrences）指令、扩展（extended）指令、保留（reserved）指令。
    
    保留指令一共有三条，其中一条留给调试器，用于实现断点，操作码是202(0xCA),助记符是breakpoint。另外两条
    留给Java虚拟机实现内部使用，操作码分别是254和255，助记符是impdep1和impdep2.这三条指令不允许
    出现在class文件中。
    
    本章会实现11类中的9类，在本章讨论native方法时会用到保留指令impdep1指令，引用指令分布在第6，7，8，10章中。为了方便管理，
    我们将新建instructions目录并有10个子目录。
    
- 指令和指令解码

    Java虚拟机规范介绍了Java虚拟机的大致逻辑，如下
    ```
      do {
        atomically calculate pc and fetch opcode at pc;
        if (operands) fetch operands;
        execute the action for the opcode;
      } while (there is more to do);
    ```
    大致逻辑包含:计算pc，指令解码，指令执行，用go来实现代码大致如下：
    ```
      for {
        pc := calculatePC()
        opcode := bytecode[pc]
        inst := createInst(opcode)
        inst.fetchoperands(bytecode)
        inst.execute()
     }
    ```
    
    - Instruction接口
    
    [instruction.go](instructions/base/instruction.go),FetchOperands()方法从字节码中提取操作数，
    Execute()方法执行指令逻辑。有很多指令的操作数都是类似的，为了避免重复代码，按照操作数
    类型定义一些结构体，并实现了FetchOperands()方法。相当于Java中的抽象类。
    
    NoOperandsInstruction表示没有操作数的指令，所以其对应的FetchOperands()方法实现也是空的。
    
    BranchInstruction结构体表示跳转指令，Offset字段存放跳转偏移量。FetchOperands()方法从中字节码中
    读取一个uint16证书，转成int后赋给Offset字段。
                    
    存储和加载类指令需要根据索引存取局部变量表，索引由单子接操作数给出。这类指令抽象成Index8Instruction结构体，
    用Index字段表示局部变量表索引。FetchOperands()方法从字节码中读取一个int8整数，转成uint后赋给Index字段。
    
    有一些指令需要访问运行时常量池，常量池索引由两字节操作数给出。把这类指令抽象成Index16Instruction结构体，
    用Index字段表示常量池索引。Fetchoperands()方法从字节码中读取一个uint16整数，转成uint后赋给Index字段。
    
    - BytecodeReader
    
    [bytecode_reader.go](instructions/base/bytecode_reader.go)中定义了BytecodeReader结构体，code字段存放字节码，
    pc字段记录读取到了哪个字节。为了避免每次解码指令都新创建一个BytecodeReader为实例，定义一个Reset()方法。
    
    还需要定义两个方法：ReadInt32s()和SkipPadding()。这两个方法只有tableswitch和lookupswitch指令使用，介绍这两条指令
    时再给出代码。
    
    在接下来的9个小节中，将要按照分类依次实现150条指令，虽然数目众多，但是指令很多类似，比图iload、lload、fload、dload和
    aload这5条，除了操作的数据不同，代码几乎相同。
    
    
- 常量指令

    常量指令把常量推入操作数栈顶。常量可以来自三个地方：隐含在操作码里、操作数和运行时常量池。常量指令共有21条。本节实现
    其中18条，另外三条ldc系列指令用于从运行时常量池加载常量将在第6章介绍。
    
    - nop指令
    
    [nop.go](instructions/constants/nop.go) 最简单的一条指令，什么也不做
    
    - const系列指令
    
    这一系列指令把隐含在操作码中的常量值推入操作数栈顶。在instruction\constants目录下创建[const.go](instructions/constants/const.go)文件，
    在其中定义15条指令。account_null指令把null引用推入操作数栈顶
    ```
        func (self *ACONST_NULL) Execute(frame *rtda.Frame) {
            frame.OperandStack().PushRef(nil)
        }
    ```
    dconst_0指令把double型0推入操作数栈顶
    ```
        func (self *DCONST_0) Execute(frame *rtda.Frame) {
            frame.OperandStack().PushDouble(0.0)
        }
    ```
    iconst_m1指令把int型-1推入操作数栈顶
    ```
        func (self *ICONST_M1) Execute(frame *rtda.Frame) {
            frame.OperandStack().PushInt(-1)
        }
    ```
    
    - bipush和sipush指令
    
    bipush指令从操作数中获取一个byte型整数，扩展成int型，然后推入栈顶。sipush指令从操作数中获取一个short型整数扩展成int型，然后推入栈顶。
    在instructions\constants下新建[ipush.go]
    
- 加载指令
    
    加载指令从局部变量表获取变量，然后推入操作数栈顶。加载指令共33条，按照所操作的变量类型分为6类，aload系列
    操作引用类型变量、dload系列操作double类型变量、fload系列操作float变量、iload系列操作int变量、lload系列操作long变量、xaload操
    作数组。本节实现其中的25条，数组相关在第8章实现。下面以iload为例介绍加载指令
    
    在instructions/loads目录下新建[iload.go](instructions/loads/iload.go),在其中定义五条指令,iload的索引来自操作数，其余四条
    来自操作码中，代表要从局部变量表中获取第几个变量
    ```
    type ILOAD struct{ base.Index8Instruction }
    type ILOAD_0 struct{ base.NoOperandsInstruction }
    type ILOAD_1 struct{ base.NoOperandsInstruction }
    type ILOAD_2 struct{ base.NoOperandsInstruction }
    type ILOAD_3 struct{ base.NoOperandsInstruction }
    ```
    
- 存储指令
    
    存储指令是把变量从操作数栈顶弹出，放入局部变量表中，正好和加载指令相反。存储指令也分为6类，以lstore系列
    指令为例进行介绍，在instructions/stores目录下创建[lstore.go](instructions/stores/lstore.go)，在其中包含了5
    条指令，如下
    ```
    type LSTORE struct{ base.Index8Instruction }
    type LSTORE_0 struct{ base.NoOperandsInstruction }
    type LSTORE_1 struct{ base.NoOperandsInstruction }
    type LSTORE_2 struct{ base.NoOperandsInstruction }
    type LSTORE_3 struct{ base.NoOperandsInstruction }
    ```
    lstore指令的索引来自操作数，其余四条的索引隐含在操作码中