## 5. 接口

接口类型是对其他类型行为的概括和抽象。相比java的显示声明，go的独特之处在于它是`隐式实现`，对于一个具体的类型，无须声明它实现了哪些接口，只需要提供接口所必须的方法即可。

```go
// interface{}效果等同于java中的Object
func test1(args ...interface{}) {
	fmt.Println(args)
}
// any等同于interface{}
func test2(args ...any) {
	fmt.Println(args)
}
```

> 1.go中提供了一种接口类型，即`interface{}`，等价于java的Object。1.18版本后可以使用any替换，效果是等价的。

```go
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

type Iphone struct {
}

type Android struct {
}

func (i Iphone) Call() { // 方法名、参数与返回值都必须相同
	fmt.Println("i'm iphone!")
}

func (a Android) Call() {
	fmt.Println("i'm android!")
}

func test3() {
	var phone Phone // 声明接口类型的变量
	phone = new(Iphone)
	phone.Call()

	phone = new(Android)
	phone.Call()
}
```

> 1. 接口的动态值：接口的类型信息 & 接口值的具体值或底层对象的副本。
> 2. go的接口是隐式实现，一个类型只需要实现接口所定义的`所有方法`，无需显示声明。
> 3. 接口是值类型，可以作为`函数的参数和返回值`，也可以`赋值给其他接口变量`。接口类型的零值是`nil`。
> 4. 接口变量在运行时存放的是`实际的对象和值`，通过类型断言可以将接口变量转换为具体类型，并访问其属性和方法。

### 5.1 接口断言

类型断言是一个作用在接口值上的操作，即`x.(T)`其中x是一个接口类型的表达式，T是一个类型（称为断言类型）。即`x的动态类型是否是T`。

```go
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

	if s, ok := p.(Say); ok { // ok返回bool值，表示p的类型是否是Say
		fmt.Printf("%#v\n", s)
	} else {
		fmt.Println("panic: interface conversion")
	}
}
```