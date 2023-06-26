package chapter6

import (
	"flag"
	"fmt"
)

func Main() {
	test1(1, "2", map[string]string{"1": "1"}, []int{1, 2, 3}, struct {
		Name string
	}{"jack"})
	test2(1, "2", map[string]string{"1": "1"}, []int{1, 2, 3}, struct {
		Name string
	}{"jack"})
	test3()
	//test4()
	test5()
}

// interface{}效果等同于java中的Object
func test1(args ...interface{}) {
	fmt.Println(args)
}

// any等同于interface{}
func test2(args ...any) {
	fmt.Println(args)
}

func test3() {
	var phone Phone
	phone = new(Iphone)
	phone.Call()

	phone = new(Android)
	phone.Call()

}

func test4() {
	s := flag.String("s", "str", "string")
	n := flag.String("i", "0", "number")
	flag.Parse()
	fmt.Printf("string: %s, number: %s\n", *s, *n)
}

func test5() {
	var p Phone
	p = Iphone{}
	i := p.(Iphone) // 接口断言：如果p的动态类型是Iphone，那么会返回p的动态值
	i.getIphone()   // 调用Iphone类型独有的方法（非实现的接口方法）

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("s0 print error: %s\n", r)
		}
	}()
	s0 := p.(Say) // 会产生panic
	s0.doSay("hello")

	if s, ok := p.(Say); ok {
		fmt.Printf("%#v\n", s)
	} else {
		fmt.Println("panic: interface conversion")
	}
}

type Object interface {
	equals(source any, target any) bool // 此处不需要定义接收者类型，与给struct定义方法不同
	Say                                 // 接口嵌套接口
}

type Say interface {
	doSay(something string) (word string)
}

type Phone interface {
	Call()
}

type Iphone struct { // 理解成java中实现接口的类
}

type Android struct {
}

func (i Iphone) Call() { // 理解成类中实现接口的方法
	fmt.Printf("%#v\n", i) // 存储的是实际的对象类型和值，而不是接口的类型
	fmt.Println("i'm iphone!")
}

func (a Android) Call() {
	fmt.Printf("%#v\n", a) // 存储的是实际的对象类型和值
	fmt.Println("i'm android!")
}

func (i Iphone) getIphone() Phone { // 接口作为返回值类型
	return Iphone{}
}

func (a Android) getAndroid() Phone {
	return Android{}
}

func doCall(p Phone) { // 接口作为函数的参数
	p.Call()
}
