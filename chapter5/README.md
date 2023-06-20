## 4. 方法

#### 4.1 定义

在go中没有类的概念，一般用struct来代替类的操作，但struct中只包含字段属性，所以go提供了一种特殊的函数，其通过`作用在某个接收者上面来实现与其关联`。

```go
// func (接收者 接收者类型) 方法名(参数列表) 返回值类型 {}

func main() {
    p := Person{"jack", 18}
	fmt.Printf("person: %p, person name: %s\n", &p, p.Name)
	p.printPerson()     // 基于值复制的传参
	(&p)._printPerson() // 基于指针的传参
	p._printPerson()    // 与(&p)效果相同，因为编译器会做隐式的转换

	var p1 *Person
	p1._printPerson() // 当指针类型作为接收者时，类型可以为nil，值类型不可以为nil
}

func (p Person) printPerson() { // person作为接收者类型
	fmt.Printf("person: %p, person name: %s\n", &p, p.Name)
}
// 允许nil作为接收者
func (p *Person) _printPerson() { // 指针作为接收者类型
	if p == nil {
		fmt.Printf("person: %p", p)
		return
	}
	(*p).Name = "lucy"
	fmt.Printf("person: %p, person name: %v\n", p, (*p).Name)
}

type Person struct {
	Name string
	Age  int
}
```

> 1. 相比函数，方法与其的区别在于多了接收者，即`把方法绑定到这个接收者对应的类型上`。
> 2. 接收者可以是任何类型，但不能是接口，因为接口是抽象定义，方法确实具体实现。
> 3. 允许出现接收者类型不同，但是方法名相同的方法。
> 4. 接收者类型是指针时使用的是`引用传递`，其他则是`值传递`，针对同一个结构体的方法的接收者类型需要统一。
> 5. 只有当接收者类型为指针时，可以将`nil作为实参`传入。
> 6. 当实参是T类型的变量而形参的接收者是*T类型时，`编译器会隐式的获取变量的地址`。

#### 4.2 结构体内嵌的类型

```go
type Inner struct{ A, B int }

func (i Inner) inner() {
	fmt.Printf("Inner: A: %d, B: %d\n", i.A, i.B)
}

func (o Outer) outer() {
	fmt.Printf("Outer: C: %s\n", o.C)
}

type Outer struct {
	Inner
	C string
}

func test3() {
	o := &Outer{Inner{1, 2}, "3"}
	o.outer()
	o.inner() // 作用等同于o.Inner.inner()
}
```

> 通过结构体的嵌套，可以无感的调用内嵌的结构体的方法。

#### 4.2 方法变量和方法表达式

```go
type Person struct {
	Name string
	Age  int
}

func (p *Person) modifyPerson() { // 指针作为实参的接收者类型
	fmt.Printf("person: %#v\n", p) // 打印接收者参数
}

func test5(f func(*Person)) { // 方法表达式作为形参需要指定类型
    p := &Person{"lucky", 18}
    f(p)
}

func main() {
    p := &Person{"jack", 20}
    person := p.modifyPerson // 获取方法变量，只能在方法所属的实例p上调用
    person() // 调用方法 效果等同于 p.modifyPerson() 不需要指定接收者
    
    person2 := (*Person).modifyPerson // 获取方法表达式
    person2(p) // 调用方法表达式时需要传入接收者实参
    
    test5(person2) // 将方法表达式作为函数参数传递
}
```

> 1. 方法变量：接`将方法赋值给一个变量`，不需要指定接收者，只能在方法所属类型的实例上调用，不能在其他类型的实例上调用。
> 2. 方法表达式：方法表达式的语法是`T.methodName`，需要`显式指定接收者类型`作为第一个参数。可以作为普通函数一样传递。

#### 4.3 封装

如果变量或者方法不是通过对象访问到的，这称作`封装的变量或者方法（又称数据隐藏）`。

```go
// 在另一个包中调用Logger的方法
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
```

> 1. go中的封装效果等同于Java中的类字段属性设为私有，但是对外提供getter和setter用来修改字段属性。
> 2. go中的getter方法是默认去掉get的，比如getAge在go中就是Age，注意需要大写（因为需要对其他包导出）。