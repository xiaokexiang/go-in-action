## The Go Programming Language
### 程序结果
#### 名称

- 函数、变量、常量、类型、语句标签和包都遵循: `字母或下划线开头，后面跟随任意数量的字符，数字和下划线，并且区分大小写`。
- 若实体在函数中声明，那么实体只在函数局部有效。如果声明在函数外，将对包内所有源文件可见。
- 实体的第一个字母的`大小写`决定其可见性是否跨包，大写表示对包外是可见和可访问的。

#### 声明

```go
package main

const commonConstant = "common"

func main() {
	var a = commonConstant
	var b string // 会默认初始化为""
	var c string = "c"
	var d = "d"
	e := "e"                  // 短变量声明
	x, y := true, 2.3 // 一次定义多个变量
	fmt.Println(a)            // common
	fmt.Println(b)            // ""
	fmt.Println(c)            // c
	fmt.Println(d)            // d
	fmt.Println(e)            // e
	fmt.Println(x)
	fmt.Println(y)
	c, d = d, c // 交换c和d的值
	fmt.Println(c)
	fmt.Println(d)
}
```

> 1. `:=`表示声明，类似`A a = new A()`。`=` 表示赋值，类似`b = a`。

#### 指针

指针的值是一个`变量的地址`。使用指针可以在无须知道变量名称的情况下，简洁读取或更新变量的值。

```go
package main

func pointer() {
	a := "1"
	p := &a         // &z表示获取一个指向整型变量的指针,类型是整型指针(*int)
	fmt.Println(*p) // *p表示p指向的变量
	*p = "2"
	fmt.Println(a) // 2

	var x, y int
	fmt.Println(&x == &x, &x == &y, &x == nil) // true false false

	q := 1
	incr(&q)       // &q所指向的值加1
	fmt.Println(q) // 2
	fmt.Println(&q) //指针
	fmt.Println(incr(&q)) // 3
}

func incr(p *int) int {
	*p++ // 递增p所指向的值,p本身不变(p是一个指针)
	return *p
}
```

> 1. &z表示获取一个指向整型变量z的指针,类型是整型指针(*int)
> 2. `*p`表示指针p指向的变量，`*p`代表一个变量。而p则代表指针，是一个0x开头的地址。