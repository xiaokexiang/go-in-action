package main // main包决定main函数启动入口
import (
	"fmt"
	pkg1 "go-in-action-2025/package1" // 给包设置别名
	pkg2 "go-in-action-2025/package2"
	_ "go-in-action-2025/package3" // 不使用包中的方法,只使用init\
	. "go-in-action-2025/package4" // 导入包下所有的公有方法,不推荐
)

func main() {
	//fmt.Println("hello go")
	//export()
	//constant()
	//baseFunction()
	//importInit()
	pointer()
}

/*
java和go的基本数据类型都是`值复制`传递
特性	        Java 引用类型	            Go 指针类型
传递内容	    对象引用的副本（指针的副本）	指针的副本（地址的副本）
修改对象属性	✅ 影响原始对象	        ✅ 影响原始对象
重新赋值引用	❌ 不影响原引用	         ❌ 不影响原指针
底层本质	    仍是值传递（传递引用副本）	仍是值传递（传递地址副本）

&: 对变量取址
*: 对指针取值
*/
func pointer() {
	a := 1
	b := 1
	changeValue(a)
	fmt.Println("a= ", a)
	changeValuePointer(&b)
	fmt.Println("b= ", b)

	fmt.Println("before -> a,b= ", a, b)
	swap(&a, &b)
	fmt.Println("before -> a,b= ", a, b)

	// 二级指针...N级指针
	var c *int = &a
	var d **int = &c
	fmt.Println("c,d= ", c, d)
}

func changeValuePointer(p *int) {
	*p = 10 // *p表示p指向的地址 p表示自己在内存中的值
}

func changeValue(p int) {
	p = 10
}

func swap(a, b *int) {
	tmp := *a
	*a = *b
	*b = tmp
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
	Echo() // 不需要.的方式调用
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
