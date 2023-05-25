package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("=============> terminal <=============")
	terminal()
	fmt.Println("=============> forLoop <=============")
	forLoop()
	fmt.Println("=============> duplicate <=============")
	//duplicate()
	fmt.Println("=============> fetch_all <=============")
	fetchAll()
	fmt.Println("=============> web_server <=============")
	//webServer()
	fmt.Println("=============> declaration <=============")
	declaration()
	fmt.Println("=============> pointer <=============")
	pointer()
}

/**
  读取terminal参数，自动将入参转为slice（动态容量的有序数组）
  go build main.go
  ./main Mon Tue Wed Thu Fri Sat Sun
*/
func terminal() {
	fmt.Println(os.Args[0])   // 第一位为命令本身: ./main
	fmt.Println(os.Args[1:3]) // 读取命令行第1、2个参数 半开区间（左闭右开）: [Mon Tue]
	fmt.Println(os.Args[1:])  // 读取命令行第一个之后的所有参数
}

/**
for initialization; condition; post {
}
tips:
1. 大括号是必须的，且左大括号必须和post在同一行
2. initialization：最先执行，可以是变量申明赋值语句或一个函数调用。initialization和post都可以省略，此时分号可以省略，效果等同while循环
3. 若condition也没有，效果等同于for{}，可以通过break和continue打断循环。
*/
func forLoop() {
	for i := 0; i < len(os.Args); i++ {
		fmt.Println(os.Args[i])
	}

	/*
		1. range方法返回slice的index和对应index的值，假设我们不需要使用index，但是go不允许出现不使用的临时变量。
		2. 使用_（空标识符）：空标识符可以用在任何语法需要变量名称但是程序逻辑不需要的地方
	*/
	for _, arg := range os.Args[3:] {
		fmt.Println(arg)
	}
	for i, arg := range os.Args[3:] {
		fmt.Println(i, arg)
	}
}

func duplicate() {
	counts := make(map[string]int)      // make构建map
	input := bufio.NewScanner(os.Stdin) // 读取命令行输入
	for input.Scan() {                  // 读到新行返回true
		counts[input.Text()]++
		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\t%s\n", n, line) // %d 数值 %s 字符串
			}
		}
	}
}

func fetchAll() {
	start := time.Now()
	ch := make(chan string)
	arr := [3]string{"https://cn.bing.com", "https://baidu.com", "https://sogou.com"}
	for _, i := range arr {
		go fetch(i, ch) // 启动goroutine
	}
	for range arr {
		fmt.Println(<-ch) // 从通道接受
	}
	fmt.Printf("%.2f elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintln(err)
		return
	}
	body, err := io.Copy(ioutil.Discard, resp.Body) // 丢弃响应内容， copy会返回响应内容大小
	_ = resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintln(err)
		return
	}
	sec := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", sec, body, url)
}

func webServer() {
	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe("localhost:18080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path) // 将url作为响应输出
}

const commonConstant = "common"

func declaration() {
	var a = commonConstant
	var b string // 会默认初始化为""
	var c string = "c"
	var d = "d"
	e := "e"          // 短变量声明
	x, y := true, 2.3 // 一次定义多个变量
	fmt.Println(a)    // common
	fmt.Println(b)    // ""
	fmt.Println(c)    // c
	fmt.Println(d)    // d
	fmt.Println(e)    // e
	fmt.Println(x)
	fmt.Println(y)
	c, d = d, c // 交换c和d的值
	fmt.Println(c)
	fmt.Println(d)

}

func pointer() {
	a := "1"
	p := &a         // &z表示获取一个指向整型变量的指针,类型是整型指针(*int)
	fmt.Println(*p) // *p表示p指向的变量
	*p = "2"
	fmt.Println(a) // 2

	var x, y int
	fmt.Println(&x == &x, &x == &y, &x == nil) // true false false

	q := 1
	incr(&q)              // &q所指向的值加1
	fmt.Println(q)        // 2
	fmt.Println(&q)       //指针
	fmt.Println(incr(&q)) // 3
}

func incr(p *int) int {
	*p++ // 递增p所指向的值,p本身不变(p是一个指针)
	return *p
}
