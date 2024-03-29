## 3. 函数

### 3.1 函数声明

```go
/*
func func_name(parameter-list) (result-list){}
*/
func x0(a, b int, c string) (x, y int, z string) {
	return // 等同于return x,y,z 称为裸返回（按照返回值顺序返回）
}
```

> 1. 函数的参数传递是按`值传递`的，函数接收到是实参的副本。
> 1. 存在多个相同类型的参数时，可以使用简写 x,y int。
> 2. 可以给返回值进行声明，声明为`局部变量`，并会根据类型`初始化为零值。`
> 3. 若已经给返回值进行变量声明，那么函数中可以直接使用return返回。
> 3. 当两个函数拥有相同的形参列表和返回列表时，那么这两个函数的`类型或签名是相同的`。

### 3.2 异常构建

```go
func main() {
	if err := waitForServer("https://github1.com", 10*time.Second, 2*time.Second); err != nil {
		fmt.Printf("waitForServer: %s\n", err)
	}
}
/*
url: 请求的url
timeout： 全部请求超时间隔
requestTimeout： 每次请求的间隔
*/
func waitForServer(url string, timeout time.Duration, requestTimeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := http.Client{
		Timeout: requestTimeout, // http client设置超时时间
	}
	for i := 0; time.Now().Before(deadline); i++ {
		_, err := client.Head(url) // 只获取响应头信息，不获取响应体
		if err == nil {
			return nil
		}
		fmt.Printf("server not responding (%s), retrying...\n", err)
		time.Sleep(time.Second << i) // 指数级增长
	}
	//return fmt.Errorf("server %s failed to respond after %s\n", url, timeout) // 创建格式化错误信息
	return fmt.Errorf("server %s failed to respond after %s\n", url, timeout) // 创建格式化错误信息
}
```

> 1. `error`是内置的接口类型，当其值为非空值，意味着失败，为控制意味着成功。
> 2. go语言中使用`普通的值`来不是异常来报告错误，使用fmt.Errorf()构建新的error错误值。

### 3.3 匿名函数

> 函数字面量就像函数声明，但在func关键字后面没有函数的名称，它是一个表达式，它的值被称为匿名函数。

```go
func main() {
    // 匿名函数作为参数
    s := strings.Map(func(r rune) rune {
	    return r + 1
    }, "HAL-9000")
    fmt.Printf("map: %#v\n", s)
    a := anonymousFunc() // 获取匿名函数的引用，因为返回的是func() int类型
    fmt.Println(a()) // 此时匿名函数a中的局部变量x=1
    fmt.Println(a()) // 此时匿名函数a中的局部变量x=2
    func(x int) { fmt.Printf("x: %d\n", x) }(1) // 构建匿名函数并调用
}

/*
匿名函数作为返回值类型
*/
func anonymousFunc() func() int {
	var x int // 能够被内层的匿名函数访问和更新
	return func() int {
		x++
		return x * x
	}
}
```

#### 3.3.1 迭代变量捕获

```go
/*
演示go的循环迭代变量捕获问题
*/
func test1() {
	var funcs []func()
	num := []int{1, 2, 3, 4}
	for _, num := range nums {
		// num := num                     // 避免迭代变量捕获问题
		funcs = append(funcs, func() { // 构建匿名函数的slice，for循环结束后num的值为4，num共用一个内存地址，每次都是更新值
			fmt.Printf("num: %d\n", num) // 匿名函数中记录的num不是值，是内存地址
		})
	}
	for _, f := range funcs {
		f() // 执行匿名函数时，num的值随着迭代结束已经变为4
	}
}
```

> 1. 迭代的变量num使用的`是同一个内存地址`，每次循环都是更新内存地址对应的值。
> 2. 在第一个循环中生成的所有匿名函数都`共享相同的循环变量num（逃逸到堆中）`，且num记录的是`迭代变量的内存地址`，而不是值。
> 3. 在匿名函数执行时，num中存储的值等于最后一次迭代的值，所以打印的都是相同的数字4。
> 4. 只需使用`局部变量拷贝循环变量的值`，这样在匿名函数中局部变量num记录的就是`某一时刻循环变量的值`，而不是循环变量的地址。

### 3.3.4 变长函数

变长函数在被调用的时候可以有`可变的参数个数`，其本质就是一个某种类型的slice。(java中的可变参数其实是个数组)

```go
func test(a ...int) {
    fmt.Printf("type: %T", a) // []int
}

func main() {
    test(1,2,3)
    a := []int{1,2,3}
    test(a...)
}
```

> 对于已经存在的一个slice，通过添加`...`实现对变长函数的调用。

### 3.3.5 延迟函数

defer用于延迟执行一个函数调用，该函数会在`当前函数返回之前`被调用执行，无论该函数是正常返回还是发生异常。存在多个defer语句时，调用的顺序是`后进先出`。常用于`关闭文件、资源释放、释放锁及跟踪函数执行`。

```go
func test3() {
	resp, err := http.Get("https://baidu.com")
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	defer resp.Body.Close() // 关闭资源
}

func test4() {
	fmt.Println("step1")
	defer fmt.Println("step2")
	defer fmt.Println("step3")
    // step1 -> step3 -> step2
}

func test5() {
	defer test6("test5")() // 会获取函数但会等到return前再执行
	time.Sleep(2 * time.Second)
	fmt.Println("test5 do something...")
}

func test6(msg string) func() {
	now := time.Now() // 会在一开始获取函数值时就计算
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s(%s)", msg, time.Since(now))//time.Since(now)会在test5返回时才会计算
	}
}
```

> 1. 执行`defer test6("test5")()`时，会先获取函数值（即test6中返回的匿名函数）但不会立马执行，会等到test5()方法返回前再执行。
> 2. defer函数执行时，其函数中的变量（即time.Since(now)）才会被确定，并在函数返回前保持不变。

### 3.3.6 宕机与恢复

当程序遇到无法处理的错误或异常情况时，可以使用 `panic` 函数引发 panic。Panic会导致程序立即停止执行，并开始执行调用栈的展开过程，在其展开过程中执行defer函数，最后程序终止。

`recover()`是一个内建函数，能够处理panic异常，其只能在defer函数中使用，用于捕获和处理panic异常。

```go
func test7() {
    defer fmt.Println("defer print")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	panic("Something went wrong")                // 宕机并输出日志信息和栈转储信息到stdout
	fmt.Println("This line will not be printed") // 不会被执行
    // Recovered from panic: Something went wrong -> defer print 
}
```

> 1. 在使用recover()函数时，Panic不会导致程序停止执行，会被recover()捕获和处理。
> 2. recover()只能在defer函数中执行，存在多个defer函数时仍然按照`后进先出`的原则调用。
