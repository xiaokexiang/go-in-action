package main

import (
	"fmt"
	"sync"
	"time"
)

/*
channel无缓冲,发送会阻塞等待接受
*/
func chanWithoutCache() {
	c := make(chan int)
	go func() {
		defer fmt.Println("goroutine exit ...")
		fmt.Println("goroutine ...")
		time.Sleep(1 * time.Second)
		c <- 1
	}()

	num := <-c // 无缓冲阻塞等待
	fmt.Printf("c: %d\n", num)
}

/*
结果:
向chan发送数据: 0, chan len(): 1, chan cap():3
向chan发送数据: 1, chan len(): 2, chan cap():3
向chan发送数据: 2, chan len(): 3, chan cap():3
从chan接收数据: 0
从chan接收数据: 1
从chan接收数据: 2
从chan接收数据: 3
main goroutine exit ...

解析:
goroutine 阻塞在 c <- 3（因为缓冲满）
main 阻塞在 <-c（因为缓冲空）
运行时直接让 3 从发送方传递到接收方(如果发送方和接收方同时阻塞（一个想发，一个想收），Go 运行时会直接绕过缓冲，在两者间建立临时传输通道)，然后：
接收方 (<-c) 先恢复执行，打印接收日志
发送方 (c <- 3) 后恢复执行，打印发送日志
*/
func chanWithCache() {
	c := make(chan int, 3) // 缓存为3的chan
	go func() {
		defer fmt.Printf("goroutine exit ...\n") // 不会打印: i=3的时候接收方会立刻接受数据并打印`main goroutine exit ...`后退出,子goroutine没有机会执行
		for i := 0; i < 4; i++ {
			c <- i
			fmt.Printf("向chan发送数据: %d, chan len(): %d, chan cap():%d\n", i, len(c), cap(c))
		}
	}()
	time.Sleep(2 * time.Second)
	for i := 0; i < 4; i++ {
		num := <-c // i=3的时候阻塞在这里
		fmt.Printf("从chan接收数据: %d\n", num)
	}
	fmt.Printf("main goroutine exit ...\n")
}

/*
如果chan没有close,main goroutine空等待就会提示dead lock问题，
*/
func chanWithClose() {
	c := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c) // 没有数据发送了就关闭
	}()

	for {
		if v, ok := <-c; ok { // true表示chan还有数据
			fmt.Printf("receive from c: %d\n", v)
		} else {
			break
		}
	}
	fmt.Println("exit ...")
}

/*
range迭代从channel中读取数据，如果chan关闭会退出range，如果没有会阻塞
*/
func chanWithRange() {
	c := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			c <- i
		}
		close(c)
	}()

	for data := range c {
		fmt.Printf("receive from c: %d\n", data)
	}
	fmt.Println("exit ...")
}

/*
监听多个chan读写状态
*/
func chanWithSelect() {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 6; i++ {
			fmt.Println(<-c) // c是否可读
			time.Sleep(500 * time.Millisecond)
		}
		quit <- 0
	}()
	x, y := 0, 1
	for {
		select {
		case c <- x: // c是否可写
			x = y
			y = x + y
		case <-quit:
			fmt.Printf("exit ...\n")
			return
		}
	}
}

/*
chan<- 只写
<-chan 只读
chan 可读可写
*/
func chanOnlyReadOrWrite() {
	c := convertToReadOnlyChan()
	for d := range c {
		fmt.Printf("receive from c: %d\n", d)
	}
	fmt.Println("main exit ... ")
}

func convertToReadOnlyChan() <-chan int {
	c := make(chan int, 10)
	go func() {
		defer close(c)
		for i := 0; i < 6; i++ {
			c <- i
		}
	}()
	return c // 返回只读chan无法再次写入
}

/*
sync.WaitGroup 同步机制
var wg sync.WaitGroup
wg.Add(1)

	go func(){
	  wg.Done()
	}()

wg.Wait()
*/
func chanWithWaitGroup() {
	var wg sync.WaitGroup
	c := make(chan<- int, 10) // 只写chan
	wg.Add(1)
	go func() {
		for i := 0; i < 10; i++ {
			c <- i
			fmt.Printf("send i to send-only chan: %d\n", i)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("main goroutine exit ...")

	a := 10
	a++
	a--
}

func chanTest() {
	var wg sync.WaitGroup
	c := make(chan int)
	wg.Add(1)
	go func() {
		defer close(c)
		for i := 0; i < 10; i++ {
			c <- i
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for d := range c {
			fmt.Printf("receive d: %d\n", d)
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("main exit ...")
}

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

func main() {
	//chanWithoutCache()
	//chanWithCache()
	//chanWithClose()
	//chanWithRange()
	//chanWithSelect()
	//chanWithWaitGroup()
	//chanOnlyReadOrWrite()
	//chanTest()
	test1()
}
