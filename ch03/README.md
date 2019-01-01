- class文件

    本章将详细讨论class文件的格式，编写代码来解析二进制数据。
    构成class文件的基本数据单位是字节，可以把整个class文件当作是一个字节流来处理。这些数据在class文件中以大端（big-endian）方式
    存储，为了描述class文件格式，虚拟机规范定义了u1，u2，u4三种类型来表示1，2，4字节的无符号整数，分别对应Go语言
    中的uint8.uiny16，uint32类型。
    
    相同类型的多条数据一般按表（table）的形式存储在class文件中。表由表头和表项（item）构成，表头是u2或u4的整数。假设表头是n，后面就会
    跟着n个表项数据。
    
    java虚拟机规范使用一种类似C的结构体语法来描述class文件的格式。整个class文件被描述为一个ClassFile结构，代码如下
```
ClassFile {
    u4 magic;
    u2 minor_version;
    u2 major_version;
    u2 constant_pool_count;
    cp_info constant_pool[constant_pool_count-1];
    u2 access_flags;
    u2 this_class;
    u2 super_class;
    u2 interfaces_count;
    u2 interfaces[interfaces_count];
    u2 fields_count;
    field_info fields[fields_count];
    u2 methods_count;
    method_info methods[methods_count];
    u2 attributes_count;
    attribute_info attributes[attributes_count];
}
```
  JDK提供了一个命令行工具javap，可以用它来反编译class文件，下面定义一个ClassFileTest类来作为模板解析，代码如下
```
    public class ClassFileTest {
        public static final boolean FLAG = true;
        public static final byte BYTE = 123;
        public static final char X = 'X';
        public static final short SHORT = 12345;
        public static final int INT = 123456789;
        public static final long LONG = 12345678901L;
        public static final float PI = 3.14f;
        public static final double E = 2.71828;
        public static void main(String[] args) throws RuntimeException {
            System.out.println("Hello, World!");
        }
}
```
    