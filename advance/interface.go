package main

import "fmt"

/*
1. 完全实现接口内所有方法才是实现类，go是隐形实现接口
2. interface就是指向子类的指针，所以需要给内存地址。即 var a Animal = &Dog{Color{}}
*/

type Animal interface {
	Sleep()
	GetColor() string
}

type Color struct {
	color string
}

type Dog struct {
	Color
}

func (d *Dog) Sleep() {
	fmt.Printf("Dog: %v sleeping ...\n", d)
}

func (d *Dog) GetColor() string {
	return d.color
}

type Cat struct {
	Color
}

func (c *Cat) GetColor() string {
	return c.color
}

func (c *Cat) Sleep() {
	fmt.Printf("Cat: %v sleeping ...\n", c)
}

func ShowAnimal(a Animal) {
	fmt.Printf("animal type: %T, color: %s\n", a, a.GetColor())
}

func ShowAnimal2(a interface{}) {
	if v, ok := a.(*Dog); ok { // 注意*Dog和Dog不是一种类型
		fmt.Printf("type is Dog: %v\n", v)
	} else {
		fmt.Printf("type is not Dog\n")
	}
}

func main() {
	var d Animal = &Dog{Color{color: "red"}}
	d.Sleep()
	fmt.Printf("Dog color: %s\n", d.GetColor())
	ShowAnimal(d)
	ShowAnimal2(d)
	var c Animal = &Cat{Color{color: "blue"}}
	c.Sleep()
	c.GetColor()
	fmt.Printf("Cat color: %s\n", c.GetColor())
	ShowAnimal(c)
	ShowAnimal2(c)
}
