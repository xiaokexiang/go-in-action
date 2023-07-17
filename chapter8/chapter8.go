package chapter8

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func Main() {
	//test6()
	//test7([]string{"https://cn.bing.com", "https://baidu.com", "https://leejay.top", "https://cn.bing.com"})
	//test8()
	//test9()
	test10()
}

var ch1 = make(chan int) // 发送存款额
var ch2 = make(chan int) // 接收余额
func Deposit(amount int) {
	ch1 <- amount
}

func Balance() int {
	return <-ch2
}

/*
1. 不要通过共享内存来通信，而应该通过通信来共享内存。
2. 使用通道请求来代理一个受限变量的所有访问的goroutine称为`该变量的监控goroutine`
*/
func teller() {
	var balance int
	for {
		select {
		case amount := <-ch1: // 捕获通道发送的数据
			balance += amount
			fmt.Printf("当前余额: %d\n", balance)
		case ch2 <- balance: // balance初始化是nil，所以默认阻塞直到不为nil
		}
	}
}

/*
init函数：执行先于main，导入包时自动执行
*/
func init() {
	go teller() // 启动监控goroutine
}

func test1() {
	go func() {
		Deposit(200)
	}()
	go Deposit(300)
	time.Sleep(2 * time.Second)
}

var (
	mute    sync.Mutex
	balance int // 共享变量需要紧接着Mutex声明之后
)

func test2() {
	go func() {
		mute.Lock() // 尝试获取锁，若无法获取会阻塞到锁被其他goroutine释放
		// 临界区域开始，临界区域内可自由读取和修改共享变量
		balance += 200
		fmt.Printf("当前余额: %d\n", balance)
		// 临界区域结束
		mute.Unlock() // 建议配合defer使用
	}()

	mute.Lock()
	balance += 300
	fmt.Printf("当前余额: %d\n", balance)
	mute.Unlock()

	time.Sleep(time.Second)
}

var (
	mutex    sync.RWMutex // 读写锁，允许读操作并发执行，写操作则需要完全独享的访问权限
	wg       sync.WaitGroup
	channel  = make(chan struct{})
	channel2 = make(chan struct{})
)

func test3() {
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-channel:
				fmt.Println("任务执行完毕")
				return
			case <-channel2:
				fmt.Printf("balance: %d\n", Balance2())
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			Deposit2(i * 100)
			channel2 <- struct{}{}
		}
		channel <- struct{}{}
	}()

	wg.Wait()
	close(channel)
	close(channel2)
}

func Deposit2(num int) {
	mutex.Lock()
	defer mutex.Unlock()
	balance += num
}

func Balance2() int {
	mutex.RLock()         // 读 加锁
	defer mutex.RUnlock() // 读 解锁
	return balance        // 读锁的临界区内不要有更新共享变量的操作
}

var (
	m sync.RWMutex // 读写锁
	b int          // 共享变量
)

/*
模拟读多写少的场景
*/
func test4() {
	for i := 0; i < 100; i++ {
		go func() {
			balance3()
		}()
	}

	for i := 0; i < 10; i++ {
		deposit3(i * 100)
	}
}

func balance3() {
	defer m.RUnlock()
	m.RLock()
	fmt.Printf("balance: %d\n", b)
}

func deposit3(num int) {
	m.Lock()
	defer m.Unlock()
	balance += num
	fmt.Printf("deposit: %d\n", num)
}

func test5() {
	var x, y int

	go func() {
		x = 1
		fmt.Printf("y: %d\n", y)
	}()

	go func() {
		y = 1
		fmt.Printf("x: %d\n", x)
	}()
	time.Sleep(time.Second)
}

var (
	once sync.Once // 包含bool和Mutex
	p    Person
)

type Person struct {
	Name string
}

func test6() {
	for i := 0; i < 3; i++ {
		go func() {
			initPerson()
			fmt.Printf("person: %#v\n", p)
		}()
	}
	time.Sleep(1 * time.Second)
}

func initPerson() {
	once.Do(func() { // 每次调用do的时候都会锁定互斥量并检查里面的bool值，为false就执行函数，为true就不执行，对所有goroutine可见。
		fmt.Println("init")
		p = Person{"jack"}
	})
}

type Func func(key string) (any, error) // 定义函数类型

func httpGetBody(url string) (any, error) { // 此函数入参与返回与Func类型一致
	fmt.Printf("url: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type result struct {
	value any
	ready chan struct{}
	err   error
}

type Memo struct {
	f     Func               // 函数类型
	cache map[string]*result // key是url地址，value是请求的结果+error
}

func New(f Func) *Memo {
	return &Memo{f, make(map[string]*result)}
}

/*
Get 判断url是否在缓存中，不在就缓存，在就直接返回缓存
1. 通过互斥锁来保证多个goroutine判断缓存的并发安全，如果缓存不存在就初始化
2. 另一个goroutine进来后，会等待直到通道被关闭。
*/
func (memo *Memo) Get(key string) (any, error) {
	rw.Lock()
	r, ok := memo.cache[key]
	if !ok {
		r = &result{ready: make(chan struct{})}
		memo.cache[key] = r
		rw.Unlock()
		(*r).value, (*r).err = memo.f(key) // 进行缓存，调用Func方法执行并缓存
		close((*r).ready)
	} else {
		rw.Unlock()
		<-(*r).ready // 等待数据准备完毕
	}
	return (*r).value, (*r).err
}

var (
	wg1 sync.WaitGroup
	rw  sync.Mutex
)

func test7(urls []string) {
	m := New(httpGetBody) // 初始化
	for _, url := range urls {
		wg1.Add(1)
		go func(u string) {
			defer wg1.Done()
			start := time.Now()
			value, err := m.Get(u)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s, %s, %d bytes\n", u, time.Since(start), len(value.([]byte)))
		}(url)
	}
	wg1.Wait()
}

func test8() {
	start := time.Now()
	ch1 := make(chan any)
	go func() {
		for {
			ch1 <- struct{}{}
		}
	}()

	for range ch1 {

	}
	fmt.Printf("time: %s", time.Since(start))
}

func test9() {
	max := runtime.GOMAXPROCS(-1) // 输入<=0的值就是返回上次一次设置的参数，默认和CPU核数相同
	fmt.Printf("GOMAXPROCS: %d\n", max)
	runtime.GOMAXPROCS(4)
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
}

func test10() {
	// GOMAXPROCS=1和!=1时输出的不同体现goroutine的调度
	runtime.GOMAXPROCS(2)
	for {
		go fmt.Print(0)
		fmt.Print(1)
	}
}
