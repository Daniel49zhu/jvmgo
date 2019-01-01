package classpath

import (
	"os"
	"strings"
)

//存放了路径分隔符';'
const pathListSeparator = string(os.PathListSeparator)

type Entry interface {
	//负责寻找和加载class文件，参数是class文件的相对路径，路径之间用斜线（/）分割，文件名
	//有.class后缀。比如要读取java.lang.Object类，传入的参数应该是java/lang/Object.class，
	//返回值是读取的字节数据，最终定位到class文件的Entry以及Error信息
	readClass(className string) ([]byte, Entry, error)
	//用于返回变量的字符串表示
	String() string
}

//根据参数不同创建不同的Entry实例，
func newEntry(path string) Entry {
	if strings.Contains(path, pathListSeparator) {
		return newCompositeEntry(path)
	}

	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}

	if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") || strings.HasSuffix(path, ".ZIP") {

		return newZipEntry(path)
	}

	return newDirEntry(path)
}
