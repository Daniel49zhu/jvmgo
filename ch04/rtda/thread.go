package rtda

type Thread struct {
	pc    int // the address of the instruction currently being executed
	stack *Stack
	// todo
}

func NewThread() *Thread {
	return &Thread{
		stack: newStack(1024),
	}
}

func (self *Thread) PC() int {
	return self.pc
}

//虚拟机通过-Xss来指定虚拟机栈的大小
func (self *Thread) SetPC(pc int) {
	self.pc = pc
}

//入栈
func (self *Thread) PushFrame(frame *Frame) {
	self.stack.push(frame)
}

//出栈
func (self *Thread) PopFrame() *Frame {
	return self.stack.pop()
}

//返回当前帧
func (self *Thread) CurrentFrame() *Frame {
	return self.stack.top()
}
