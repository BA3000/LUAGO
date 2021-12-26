package main

import (
	"io/ioutil"
	"luago/state"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		// 读取二进制chunk
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil { panic(err) }
		ls := state.New()
		// 加载到栈顶
		ls.Load(data, os.Args[1], "b")
		// 运行主函数
		ls.Call(0, 0)
	}
}
