package main

import "fmt"

// 定义一个新的类型，这是int的一个别名
type helloInt int

type person struct {
	name string
	age  int
}

func (p *person) getName() string { // 注意这里要用指针类型
	return p.name
}

func (p *person) setName(n string) {
	p.name = n
}

func (p *person) getAge() int {
	return p.age
}

func (p *person) setAge(a int) {
	p.age = a
}

func main() {
	//typeAlias()
	//structFunc()
	structClass()
}

func typeAlias() {
	var i helloInt = 10
	fmt.Printf("value: %d, type of helloInt: %T\n", i, i) // main.helloInt
}

/*
结构体作为参数传递是值类型传递，是结构体的副本，不是指针传递
*/
func structFunc() {
	p := person{
		name: "zs",
		age:  18,
	}
	fmt.Printf("person: %v\n", p)
	modifyPerson(p)
	fmt.Printf("person: %v\n", p)
	modifyPersonPtr(&p)
	fmt.Printf("person: %v\n", p)
}

func modifyPerson(p person) {
	p.name = "ls"
}
func modifyPersonPtr(p *person) {
	p.name = "ls" // 等价于(*p).name = "ls"
}

func structClass() {
	p := person{
		name: "ls",
		age:  28,
	}
	p.setName("zs")
	p.setAge(18)
	fmt.Printf("person: %v\n", p)

	p2 := &person{
		name: "ls",
		age:  28,
	}
	p2.name = "ls_"
	fmt.Printf("person: %v\n", p2)
}
