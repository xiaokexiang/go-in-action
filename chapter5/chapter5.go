package chapter5

import (
	"fmt"
)

func Main() {
	test1()
	test2()
	test3()
	test4()
	test6()
}

func test6() {
	l := &Logger{}
	l.SetFlags(1)
	fmt.Printf("flags: %d", l.Flags())
}

// Logger 定义实体，但字段是当前包私有，无法导出的
type Logger struct {
	flags  int
	prefix string
}

// Flags 等同于java中字段的getFlags方法，go默认去掉get，且可导出，所以方法名大写
func (logger *Logger) Flags() int {
	return logger.flags
}

func (logger *Logger) SetFlags(flag int) {
	logger.flags = flag
}

func test4() {
	p := &Person{"jack", 20}
	person := p.modifyPerson          // 获取方法变量
	person()                          // 效果等同于p.modifyPerson()
	person2 := (*Person).modifyPerson // 方法表达式
	person2(p)                        // 调用方法表达式时需要传入接收者实参
	test5(person2)
}

func test5(f func(*Person)) {
	p := &Person{"navy", 18}
	f(p)
}

func (p *Person) modifyPerson() {
	fmt.Printf("person: %#v\n", p)
}

// 结构体内嵌作为接收者类型
func test3() {
	o := &Outer{Inner{1, 2}, "3"}
	o.outer()
	o.inner() // 作用等同于o.Inner.inner()
}

type Inner struct{ A, B int }

func (i Inner) inner() {
	fmt.Printf("Inner: A: %d, B: %d\n", i.A, i.B)
}

type Outer struct {
	Inner
	C string
}

func (o Outer) outer() {
	fmt.Printf("Outer: C: %s\n", o.C)
}

func test2() {
	m := &Values{
		"contextType": {
			"application/json",
		},
	}
	m.Add("contextType", "application/xml")
	fmt.Printf("m: %#v\n", m)
	fmt.Printf("m[\"contextType\"]: %s\n", m.Get("contextType"))
}

type Values map[string][]string

func (v *Values) Get(key string) []string {
	if vs := (*v)[key]; len(vs) > 0 {
		return vs
	}
	return []string{}
}

func (v *Values) Add(key, value string) {
	(*v)[key] = append((*v)[key], value)
}

func test1() {
	p := Person{"jack", 18}
	fmt.Printf("person: %p, person name: %s\n", &p, p.Name)
	p.printPerson()     // 基于值复制的传参
	(&p)._printPerson() // 基于指针的传参
	p._printPerson()    // 与(&p)效果相同，因为编译器会做隐式的转换

	var p1 *Person
	p1._printPerson() // 当指针类型作为接收者时，实参可以为nil，值类型时则不可以为nil
}

func (p Person) printPerson() {
	fmt.Printf("person: %p, person name: %s\n", &p, p.Name) // person的地址每次都不相同
}

func (p *Person) _printPerson() {
	if p == nil {
		fmt.Printf("person: %p\n", p)
		return
	}
	(*p).Name = "lucy"
	fmt.Printf("person: %p, person name: %v\n", p, (*p).Name)
}

type Person struct {
	Name string
	Age  int
}
