package chapter4

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Main() {
	x, y, z := x0(1, 2, "3")
	fmt.Printf("x: %d, y: %d, z: %s\n", x, y, z)
	x1(x0(1, 2, "3")) // 10是实参，y是形参
	//if err := waitForServer("https://github1.com", 10*time.Second, 2*time.Second); err != nil {
	//	fmt.Printf("waitForServer: %s\n", err)
	//}
	s := strings.Map(func(r rune) rune {
		return r + 1
	}, "HAL-9000")
	fmt.Printf("map: %#v\n", s)
	a := anonymousFunc() // 获取匿名函数的引用
	fmt.Println(a())
	func(x int) { fmt.Printf("x: %d\n", x) }(1) // 构建匿名函数并调用
	topoSort()
	test1()
	values := []int{1, 2, 3} // 定义slice
	test2(values...)         // 将slice作为可变参数
	test3()
	test4()
	test5()
	test7()
}

/*
func func_name(parameter-list) (result-list){
}
1. 存在多个相同类型的参数时，可以使用简写 x,y int
2. 可以给返回值进行声明，声明为`局部变量`，并会根据类型初始化为零值
3. 若已经给返回值进行变量声明，那么函数中可以直接使用return返回
*/
func x0(a, b int, c string) (x, y int, z string) {
	return a, b, c
}

func x1(a, b int, c string) (z string) { // z被声明为返回值的局部变量
	fmt.Printf("z: %v\n", &z) // 每次函数被调用时z的内存地址都不相同
	z = strconv.Itoa(a) + strconv.Itoa(b) + c
	return // 效果等同于 return c
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

/*
函数的返回值类型是 func() int
通过匿名函数来创建没有函数名的函数，匿名函数可以在函数内部定义，并且可以访问其外部作用于的变量
*/
func anonymousFunc() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

/*
构建拓扑图
*/
func topoSort() {
	prereqs := map[string][]string{
		"algorithms": {"data structures"},
		"calculus":   {"linear algebra"},
		"compilers": {
			"data structures",
			"formal languages",
			"computer organization",
		},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}
	fmt.Printf("result: %#v\n", strings.Join(doTopoSort(prereqs), ","))
}

func doTopoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) { // 创建匿名函数
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}

/*
演示go的循环迭代变量捕获问题
迭代的变量使用的内存地址是同一个，每次循环更新对应地址的值
*/
func test1() {
	var funcs []func()
	nums := []int{1, 2, 3, 4}
	for _, num := range nums {
		num := num                     // 避免迭代变量捕获问题
		funcs = append(funcs, func() { // 构建匿名函数的slice，for循环结束后num的值为4，num共用一个内存地址，每次都是更新值
			fmt.Printf("num: %d\n", num) // 匿名函数中记录的num不是值，是内存地址
		})
	}
	for _, f := range funcs {
		f() // 执行匿名函数时，num的值随着迭代结束已经变为4
	}
}

func test2(a ...int) {
	fmt.Println(a)
	fmt.Fprintf(os.Stdout, "type: %T\n", a) // []int
}

func test3() {
	resp, err := http.Get("https://baidu.com")
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	defer resp.Body.Close()
}

func test4() {
	fmt.Println("step1")
	defer fmt.Println("step2")
	defer fmt.Println("step3")
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
		log.Printf("exit %s(%s)", msg, time.Since(now)) // time.Since(now)会在test5返回时才会计算
	}
}

func test7() {
	defer fmt.Println("defer print")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	panic("Something went wrong")                // 宕机并输出日志信息和栈转储信息到stdout
	fmt.Println("This line will not be printed") // 不会被执行
}
