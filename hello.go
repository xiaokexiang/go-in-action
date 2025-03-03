package main // main包决定main函数启动入口
import (
	"fmt"
	pkg1 "go-in-action-2025/package1" // 给包设置别名
	pkg2 "go-in-action-2025/package2"
	_ "go-in-action-2025/package3" // 不使用包中的方法,只使用init
)

func main() {
	//fmt.Println("hello go")
	//export()
	//constant()
	//baseFunction()
	importInit()
}

/*
程序执行顺序
main()
  |
import pkg1 -> pkg1
  |             |
const       import pkg2 -> pkg2
  |             |           |
 var          const      import pkg3 -> ...
  |             |           |
init()         var        const
  |             |           |
main()        init()       var
  |                         |
exit                      init()
*/

func importInit() {
	pkg1.Echo()
	pkg2.Echo()
}

func init() {
	fmt.Println("main init")
}

var g = 100

// := 声明且赋值
func export() {
	var a int
	fmt.Println("a= ", a)
	fmt.Printf("type of a = %T\n", a)

	var b int = 100
	fmt.Println("b= ", b)
	fmt.Printf("type of b = %T\n", b)

	var c = 100
	fmt.Println("c= ", c)
	fmt.Printf("type of c = %T\n", c)

	d := 100
	fmt.Println("d= ", d)
	fmt.Printf("type of d = %T\n", d)

	fmt.Println("g= ", g)
	fmt.Printf("type of g = %T\n", g)

	var (
		g1 = 100
		g2 = "a"
	)
	fmt.Println("g1= ", g1, ", g2= ", g2)
	h1, h2 := 100, "1"
	fmt.Println("h1= ", h1, ", h2= ", h2)

}

func constant() {
	const l = 10 // read only
	const (      // 每行的iota都加1，iota默认为0
		A = iota + 1
		B
		C
	)
	fmt.Println("A=", A)
	fmt.Println("B=", B)
	fmt.Println("C=", C)
}
func baseFunction() {
	fmt.Println(function(1, false))
	fmt.Println(function2(0, 1))
	fmt.Println(function3(0, 1))
	fmt.Println(function4(0, 1))
	fmt.Println(function5(0, 1))
}

func function(a int, b bool) bool {
	c := true
	return c
}
func function2(a int, b int) (int, string) {
	return 0, ""
}
func function3(a int, b int) (c int, d string) {
	return 0, "hello"
}
func function4(a, b int) (c, d int) {
	return 0, 1
}
func function5(a, b int) (c, d int) { // 可以理解成c,d已经初始化了
	return
}
