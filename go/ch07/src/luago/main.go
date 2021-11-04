package main

import (
	"fmt"
	. "luago/api"
	"luago/state"
	"io/ioutil"
	"os"
	"luago/binchunk"
	. "luago/vm"
)

func main() {
	if len(os.Args) > 1 {
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {panic(err)}
		proto := binchunk.Undump(data)
		luaMain(proto)
	}
}

func luaMain(proto *binchunk.Prototype) {
	// 获取运行输入函数需要用到的寄存器数量
	nRegs := int(proto.MaxStackSize)
	// 栈空间要分配多一点，因为指令实现函数也需要空间
	ls := state.New(nRegs + 8, proto)
	// 在栈里预留出寄存器空间，剩下的流量给指令实现函数
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		// 取出指令
		inst := Instruction(ls.Fetch())
		if inst.Opcode() != OP_RETURN {
			// 执行指令
			inst.Execute(ls)
			// 打印指令和栈信息
			fmt.Printf("[%02d] %s", pc + 1, inst.OpName())
			printStack(ls)
		} else {
			// 遇到return，返回
			break
		}
	}
}

func printStack(ls LuaState)  {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:	fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:	fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:	fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}