package chapter7

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Main() {
	//test9("D:/UserData/解压/netcat-1.11/")
	//test10()
	//test11(2001, 1001, true)
	//test12()
	//test13()
	//test14()
	//test15()
	test16()
}

// goroutine：每一个并发执行的活动，一个程序启动时会有一个主 goroutine（参考java的main方法的线程）。
// 主goroutine结束了会导致其他的goroutine也一同结束
func test1() {
	concurrent(false)
	go concurrent(true)
	time.Sleep(time.Second * 5)
}

func concurrent(sleep bool) {
	if sleep {
		time.Sleep(time.Second * 2)
	}
	fmt.Printf("123\n")
}

func test2() {
	if listen, err := net.Listen("tcp", "localhost:8080"); err != nil {
		log.Fatal(err)
	} else {
		for {
			if accept, err := listen.Accept(); err != nil { // accept方法阻塞直到有连接请求进来，在没有goroutine的前提下多个客户端请求，后一个请求只能等待前面的请求关闭才能获取返回
				log.Println(err)
				continue
			} else {
				go handleConnection3(accept)
			}
		}
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	for {
		_, err := io.WriteString(connection, time.Now().Format("2006-01-02 03:04:05\n")) // 类似java的yyyy-MM-dd HH:mm:ss
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// 回声
func handleConnection2(connection net.Conn) {
	_, _ = io.Copy(connection, connection)
	_ = connection.Close()
}

func handleConnection3(connection net.Conn) {
	input := bufio.NewScanner(connection)
	for input.Scan() {
		echo(connection, input.Text(), 1*time.Second)
	}
	connection.Close()
}

func echo(c net.Conn, text string, duration time.Duration) {
	_, _ = fmt.Fprintln(c, "\t", strings.ToUpper(text))
	time.Sleep(duration)
	_, _ = fmt.Fprintln(c, "\t", text)
	time.Sleep(duration)
	_, _ = fmt.Fprintln(c, "\t", strings.ToLower(text))
}

/*
ch := make(chan any)  // 类型是chan int 与new一样，make返回的是引用，零值是nil
_ = make(chan int, 3) // 创建一个缓存为3的通道
ch <- 1024            // 发送数据到通道中
x := <-ch             // 接收数据并赋值给x
close(ch)             // 设置标志值表示已经发送完毕，关闭后的发送会导致panic，在一个已经关闭的通道上进行接收操作，将获取所有已经发送的值，直到通道为空。
*/
func test3() {
	if conn, err := net.Dial("tcp", "localhost:8080"); err != nil {
		log.Fatal(err)
	} else {
		ch := make(chan struct{})
		go func() {
			_, _ = io.Copy(os.Stdout, conn)
			fmt.Println("done")
			ch <- struct{}{}
		}()
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			log.Fatal(err)
		}
		_ = conn.Close()
		fmt.Println(2)
		<-ch
	}
}

/*
管道基础版：goroutine_a -> channel -> goroutine_b -> channel -> goroutine_c
*/
func test4() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		for x := 0; x < 10; x++ {
			ch1 <- x
		}
		close(ch1)
	}()

	go func() {
		for {
			// ok: true表示接收成功 false表示当前的接收操作在一个关闭的并且读完的通道上
			// 如果不关闭ch1，那么此处会一直等待数据接收，ch2也一直在等待数据发送，主goroutine的print无法执行导致死锁
			if x, ok := <-ch1; !ok {
				break
			} else {
				ch2 <- x << 1
			}
		}
		close(ch2)
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println(<-ch2)
	}
}

/*
管道提升版：使用range迭代channel直到接收完所有的值
*/
func test5() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for x := 0; x < 10; x++ {
			ch1 <- x
		}
		/*
		  close不是必须的，只有发送方的数据都发送完毕才需要关闭通道
		  垃圾回收器会根据通道是否可以被访问来决定是否回收。
		*/
		close(ch1)
	}()

	go func() {
		for x := range ch1 { // 通过range循环读取通道所发送的值，接收完后关闭循环
			ch2 <- x << 1
		}
		close(ch2)
	}()

	for x := range ch2 {
		fmt.Println(x)
	}
}

/*
管道高级版：使用单向管道避免误用。
<-chan ：只能接收数据的channel，不允许关闭通道（箭头从通道chan出来，即从通道接收数据）
chan<- : 只能发送数据到channel，允许关闭channel（箭头进入chan，即接收发送数据到通道）
*/
func test6() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func(in chan<- int) {
		for x := 0; x < 10; x++ {
			in <- x
		}
		close(in)
	}(ch1) // 隐式的将chan int转为了chan<- int类型
	// 注意out和in对应的类型
	go func(out chan<- int, in <-chan int) {
		for x := range in {
			out <- x << 1
		}
		close(out)
	}(ch2, ch1)

	func(in <-chan int) {
		for x := range in {
			fmt.Println(x)
		}
	}(ch2)
}

func onlyReceive(ch <-chan int) {
}

func onlySend(ch chan<- int) {
}

func test7() {
	ch1 := make(chan int)
	go func() {
		for x := 0; x < 5; x++ {
			time.Sleep(time.Second)
			ch1 <- x
		}
		close(ch1)
	}()
	/*for {
		fmt.Println("waiting")
		x := <-ch1
		fmt.Println(x)
	}*/
	/*for {
		if x, ok := <-ch1; ok {
			fmt.Println(x)
		} else {
			fmt.Printf("channel close, %d\n", x)
			time.Sleep(time.Second)
			//break
		}
	}*/
	/*for {
		if x, ok := <-ch1; ok {
			fmt.Println(x)
		} else {
			fmt.Printf("channel close, %d\n", x)
			time.Sleep(time.Second)
		}
	}*/
	for x := range ch1 {
		fmt.Println(x)
	}
}

func test8() {
	ch1 := make(chan string, 3) // 设置容量为3的有缓存通道
	/*
		1. 发送操作会在队列的尾部插入一个元素，接收操作会从队列的头部移除一个元素。
		2. 如果通道满了，发送会阻塞到另一个goroutine接收操作留出可用的空间。
		3. 执行接收操作的goroutine会阻塞到另一个goroutine在通道上发送数据。
	*/
	fmt.Printf("channel通道内元素个数: %d\n", len(ch1))
	fmt.Printf("channel缓冲区容量大小：%d\n", cap(ch1))
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "baidu.com"
	}()
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "cn.bing.com"
	}()
	go func() {
		time.Sleep(3 * time.Second)
		ch1 <- "google.com"
		close(ch1) // 如果不关闭会导致死锁（因为主goroutine一直在等待数据，但是此时已经没有其他数据向通道发送了）
	}()
	for x := range ch1 {
		fmt.Printf("fast website: %s\n", x)
	}
}

func test9(path string) {
	var files []string
	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	ch1 := make(chan string)
	for _, file := range files {
		go func(file string) {
			random := rand.Intn(5)
			time.Sleep(time.Duration(random) * time.Second)
			ch1 <- file
		}(file)
	}

	for _, file := range files { // 注意这里range的不是ch1，是files，如果是ch1不主动关闭会panic
		fmt.Printf("fileName: %s\n", file)
	}
}

/*
sync.WaitGroup Wait() Add() Done()
效果等同于java的countDownLatch
*/
func test10() {
	now := time.Now()
	var total int
	var wg sync.WaitGroup
	ch1 := make(chan int)
	for x := 0; x < 100; x++ {
		wg.Add(1) // 必须在goroutine外且开始之前
		go func(x int) {
			defer wg.Done()
			ch1 <- x
		}(x)
	}
	go func() {
		wg.Wait()
		close(ch1)
	}()
	for x := range ch1 {
		total += x
	}
	fmt.Printf("time: %s, total: %d", time.Since(now), total)
}

/*
模拟java的fork/join
*/
func test11(n int, i int, parallel bool) {
	start := time.Now()
	var total int
	if parallel {
		if n < i {
			log.Fatal("i must less than n!")
			return
		}
		var wg sync.WaitGroup
		ch1 := make(chan int)
		// 计算0-n的总和，先按照i切割，分别汇总最后相加总和
		for x := 0; x <= n/i; x++ {
			s, e := x*i+1, (x+1)*i
			if e > n {
				e = n
			}
			wg.Add(1)
			go func(a int, b int) {
				defer wg.Done()
				var t int
				for m := a; m <= b; m++ {
					t += m
				}
				ch1 <- t
			}(s, e)
		}
		go func() {
			wg.Wait()
			close(ch1)
		}()
		for x := range ch1 {
			total += x
		}
	} else {
		for x := 0; x <= n; x++ {
			total += x
		}
	}
	fmt.Printf("time: %s, total: %d\n", time.Since(start), total)
}

/*
每一个case指定一次通信

	   select {
		  case <- ch1:
		    // ...
		  case x := <- ch2:
		    // use x
		  case ch3 <- y:
		    // ...
		  default:
		    // ...
		}
*/
func test12() {
	start := time.Now()
	//tick := time.Tick(time.Second) // 获取定时器的只读通道
	ticker := time.NewTicker(time.Second)
	for countdown := 5; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-ticker.C
	}
	abort := make(chan struct{})
	go func() {
		_, _ = os.Stdin.Read(make([]byte, 1)) // 读取单个字节
		abort <- struct{}{}
	}()
	select { // 等待下面通道任意一个完成，若出现多个通道同时满足，那么select默认会随机选择一个
	case <-time.After(5 * time.Second): // 等待5s后执行（实现超时等待功能）
		ticker.Stop()
		fmt.Printf("Rocket launch: %s", time.Since(start))
	case <-abort:
		fmt.Println("Launch abort!")
		return
	}
}

// 打印偶数
func test13() {
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x := <-ch:
			fmt.Printf("x: %d\n", x)
		case ch <- i:
		}
	}
}

func test14() {
	start := time.Now()
	//tick := time.Tick(time.Second) // 获取定时器的只读通道
	ticker := time.NewTicker(time.Second)
	abort := make(chan struct{})
	go func() {
		_, _ = os.Stdin.Read(make([]byte, 1)) // 读取单个字节
		abort <- struct{}{}
	}()
	for countdown := 5; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select { // 等待下面通道任意一个完成，若出现多个通道同时满足，那么select默认会随机选择一个
		case <-ticker.C:
		case <-abort:
			fmt.Println("Launch abort!")
			return
		}
	}
	ticker.Stop()
	fmt.Printf("Rocket launch: %s", time.Since(start))
}

var done = make(chan struct{})

func test15() {
	roots := flag.Args() // ./main.exe C:/ D:/ E:/
	if len(roots) == 0 {
		roots = []string{"E:/Projects"}
	}
	fileSizes := make(chan int64)
	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go func(root string) {
			walkDir(root, &wg, fileSizes)
		}(root)
	}

	go func() {
		_, _ = os.Stdin.Read(make([]byte, 1)) // 用来接收信号中断goroutine
		close(done)
	}()

	go func() {
		wg.Wait()
		close(fileSizes)
	}()
	var nfiles, nbytes int64
	for {
		select {
		case <-done:
			fmt.Println("正在中断命令...")
			for range fileSizes {
				// do nothing 用来保证正在执行的goroutine执行完毕
			}
			fmt.Println("所有goroutine执行完毕...")
			return
		case size, ok := <-fileSizes:
			if ok {
				fmt.Printf("%d files %.1f MB\n", nfiles, float64(nbytes)/1e6) // 1e6 = 1000000 1e9 = 1000000000
				nfiles++
				nbytes += size
			}
		}
	}
}

/*
递归查询目录，并计算文件大小传输到通道
*/
func walkDir(dir string, wg *sync.WaitGroup, fileSizes chan<- int64) {
	defer wg.Done()
	if cancelled() { // 查看终端通道有没有发送信号
		return
	}
	for _, entry := range recurve(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subDir := filepath.Join(dir, entry.Name())
			go walkDir(subDir, wg, fileSizes) // 递归获取
		} else {
			if info, err := entry.Info(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "du1: %v\n", err)
			} else {
				fileSizes <- info.Size()
			}
		}
	}
}

// 全局通道，上线20，用于处理goroutine数量问题
var ch1 = make(chan struct{}, 20)

/*
读取指定目录下的文件并返回
*/
func recurve(dir string) []os.DirEntry {
	defer func() {
		<-ch1 // 释放凭证
	}()
	ch1 <- struct{}{} // 获取凭证
	if entries, err := os.ReadDir(dir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	} else {
		return entries
	}
}

/*
判断任务有没有结束
*/
func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func test16() {
	if listen, err := net.Listen("tcp", "localhost:8000"); err == nil {
		go broadcaster()
		for {
			if connection, err := listen.Accept(); err == nil {
				go handleConn(connection)
			} else {
				log.Println(err)
				continue
			}
		}
	} else {
		log.Fatal(err)
		return
	}
}

/*
新的goroutine将创建一个新的ch，然后返回信息给client，将创建的ch加入到entering，并发送加入的信息给已经在通道中的client
*/
func handleConn(connection net.Conn) {
	timeout := 10 * time.Second
	ch := make(chan string)
	go clientWriter(connection, ch) // 输出内容到客户端

	_, _ = fmt.Fprintf(connection, "you are: ") // 要求client输入用户名
	scanner := bufio.NewScanner(connection)
	scanner.Scan()
	who := scanner.Text()

	messages <- who + " has arrived" // 先发送client加入的信息，后将client加入通道（顺序很重要）
	entering <- ch
	timer := time.NewTimer(timeout) // 定时器
	go func() {
		<-timer.C
		fmt.Printf("%s has disconnected\n", who)
		_ = connection.Close()
	}()

	input := bufio.NewScanner(connection)
	for input.Scan() { // 读取client发送的数据并推送到其他client
		messages <- who + ": " + input.Text()
		timer.Reset(timeout) // 重置定时器
	}
	leaving <- ch
	messages <- who + " has left"
	_ = connection.Close()

}

func clientWriter(connection net.Conn, ch chan string) {
	for msg := range ch {
		_, _ = fmt.Fprintln(connection, msg)
	}
}

type client chan<- string // 对外发送消息的通道
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]struct{})
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				select {
				case cli <- msg:
				case <-time.After(10 * time.Second):
				}
			}
		case cli := <-entering:
			clients[cli] = struct{}{} // 并不需要value，只是加入map表示client加入
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}

}
