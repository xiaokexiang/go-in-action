package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p *Person) Walk() {
	fmt.Printf("person: %s walking ...\n", p.Name)
}

type Man struct {
	Person // 继承Person结构
	Sex    string
}

func (m *Man) Say() {
	fmt.Printf("man: %s saying: i'm %s\n", m.Name, m.Sex)
}

func (m *Man) Walk() {
	fmt.Printf("man: %s walking: i'm %s\n", m.Name, m.Sex)
}

func (m *Man) ModifyName(name string) {
	m.Person.Name = name
}

func main() {
	m := Man{
		Person{
			Name: "zs",
			Age:  18,
		},
		"male",
	}
	m.Walk()
	m.Say()
	m.ModifyName("ls") // 方法是传递ptr所以可以成功修改
	m.Say()
}
