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
	//pointer()
	//deferFunc()
	arrOrSliceFunc()
}
func arrOrSliceFunc() {
	arrFunc()
	mySlice := []int{1, 2, 3}
	sliceFunc(mySlice)
	fmt.Println("--------------------------->")
	for _, v := range mySlice {
		fmt.Println("v: ", v)
	}
}

func arrFunc() {
	var arr [3]int
	arr2 := [3]int{1, 2, 3}
	fmt.Printf("arr type: %T\n", arr) // [3]int
	for i := 0; i < len(arr); i++ {
		fmt.Println("i: ", arr[i])
	}
	for i := 0; i < len(arr2); i++ {
		fmt.Println("i: ", arr2[i])
	}
	for i, v := range arr2 {
		fmt.Println("index: value: ", i, v)
	}
}

/*
1. slice又称动态数组[]int,相比数组容量不固定,数组创建时需要指定容量.
2. slice作为参数传递时, slice是引用传递而数组是值传递
3. make([]int, 3, 5) 开辟3个空间，slice的cap是5，size是3，当数量超过cap，会再次扩容一块区域，大小为5（开辟的地址不一定连续）
4. 如果不指定cap，默认等于size
*/
func sliceFunc(mySlice []int) {
	for _, v := range mySlice {
		fmt.Println("v: ", v)
	}
	mySlice[0] = 100 // slice切片会被修改
	fmt.Println("--------------------------->")
	var slice1 []int         // 默认不分配空间
	slice2 := make([]int, 3) // 开辟空间大小为3,默认值为0
	if slice1 == nil {
		fmt.Println("slice是一个空切片")
	} else {
		fmt.Println("slice是有空间的")
	}

	if slice2 == nil {
		fmt.Println("slice是一个空切片")
	} else {
		fmt.Println("slice是有空间的")
	}
	fmt.Println("--------------------------->")
	slice3 := make([]int, 2, 4)
	fmt.Printf("size: %d, cap: %d, slice: %v\n", len(slice3), cap(slice3), slice3)

	slice3 = append(slice3, 1)
	fmt.Printf("size: %d, cap: %d, slice: %v\n", len(slice3), cap(slice3), slice3)

	slice3 = append(slice3, 2)
	fmt.Printf("size: %d, cap: %d, slice: %v\n", len(slice3), cap(slice3), slice3)

	slice3 = append(slice3, 3)
	fmt.Printf("size: %d, cap: %d, slice: %v\n", len(slice3), cap(slice3), slice3) // cap -> cap * 2

	fmt.Println("--------------------------->")

	slice4 := []int{1, 2, 3}
	slice5 := slice4[0:2] // [0,2)
	slice5 = slice4[:2]   // [0,2)
	slice5 = slice4[1:]   // [1,...)
	slice5 = slice4[:]    // 全部
	slice4[0] = 100
	fmt.Println(slice4)
	fmt.Println(slice5) // s5和s4都会改变，他们指向的是同一数据

	slice6 := make([]int, 3)
	copy(slice6, slice5) // copy slice5指向的内存的数据，slice6指向新的地址空间
	slice4[0] = 101
	fmt.Println(slice6)
}

/*
1. 多个defer的顺序是先进后出
2. return的顺序早于defer
return -> defer2 -> defer1
*/
func deferFunc() int {
	defer fmt.Println("defer1")
	defer fmt.Println("defer2")
	return returnFunc()
}

func returnFunc() int {
	fmt.Println("return called ...")
	return 0
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
