package rtda

// jvm stack
//用经典的链表（linked list）数据结构来实现Java虚拟机栈，这样
//栈就可以按需使用内存空间，而且弹出的帧也可以及时被Go的垃
//圾收集器回收
type Stack struct {
	maxSize uint   //保存栈的容量
	size    uint   //栈的当前大小
	_top    *Frame //
}

func newStack(maxSize uint) *Stack {
	return &Stack{
		maxSize: maxSize,
	}
}

func (self *Stack) push(frame *Frame) {
	if self.size >= self.maxSize {
		panic("java.lang.StackOverflowError")
	}

	if self._top != nil {
		frame.lower = self._top
	}

	self._top = frame
	self.size++
}

func (self *Stack) pop() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	top := self._top
	self._top = top.lower
	top.lower = nil
	self.size--

	return top
}

func (self *Stack) top() *Frame {
	if self._top == nil {
		panic("jvm stack is empty!")
	}

	return self._top
}
