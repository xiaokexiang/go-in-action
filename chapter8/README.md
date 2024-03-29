## 7. 并发和共享变量

> 并发：每个goroutine中的步骤是顺序执行的，但是多个goroutine不能保证先后执行顺序。
>
> 竞态：在多个goroutine按某些交错的顺序执行时程序无法给出正确的结果。
>
> 数据竞态：两个goroutine并发读写`同一个变量`且至少一个是写入。
>
> 互斥：允许多个goroutine访问同一个变量，但`同一时间只有一个goroutine可以访问`。

### 7.1 sync.Mutex

在多个goroutine获取通过sync.Mutex互斥锁获取共享变量时，没有获取到锁的goroutine会阻塞到已获取锁的goroutine释放锁，在一个goroutine中的加锁与释放锁的中间区域成为`临界区`，临界区域内可自由读取和修改共享变量。（`go的互斥锁不支持重入`）。

```go
var (
	mute    sync.Mutex // 声明互斥锁
	balance int // 共享变量需要紧接着Mutex声明之后
)
func test2(){
    go func() {
        mute.Lock() // 尝试获取锁，若无法获取会阻塞到锁被其他goroutine释放
        defer mute.Unlock() // 配合defer使用
		// 临界区域开始，临界区域内可自由读取和修改共享变量
		balance += 200
		fmt.Printf("当前余额: %d\n", balance)
		// 临界区域结束
    }()
    
    mute.Lock()
	balance += 300
	fmt.Printf("当前余额: %d\n", balance)
	mute.Unlock()

	time.Sleep(time.Second)
}
```

### 7.2 sync.RWMutex

多读单写锁：允许`只读操作并发`执行，但写操作需要获得`完全独享`的访问权限。

```go
var (
	m sync.RWMutex // 读写锁
	b int // 共享变量
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
```

> 读写锁只适用于获取读锁并且锁竞争比较激烈的场景，竞争不激烈时比普通的互斥锁慢。

### 7.3 内存同步

在单个goroutine中，执行顺序是`串行一致`的

```go
var x,y int

go func(){
    x = 1
    fmt.Printf("y: %d\n", y)
}()

go func(){
    y = 1
    fmt.Printf("x: %d\n", x)
}()
// 有概率出现如下结果
// x:0 y:0 ?
// y:0 x:0 ? 
```

> 1. 在单个goroutine中，语句的执行顺序是`串行一致`的。缺少同步操作的前提下，多个goroutine之间的执行顺序无法保证。
> 2. `内存可见性`：多个处理器中，每个处理器都有自己的内存的本地缓存，在必要时才会将数据刷回内存。会导致一个goroutine的写入操作对另一个goroutine是不可见的。
> 3. `编译器和CPU重排序`：编译器和处理器可能会对代码进行重新排序，以优化执行效率。因为上文中赋值的操作和print对应不同的变量，编译器可能会交换两个语句的执行顺序。

### 7.4 sync.Once

```go
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
	once.Do(func() {
		fmt.Println("init")
		p = Person{"jack"}
	})
}
```

> sync.Once中的Do方法每次调用时都会`锁定互斥量并检查里面的bool值`，为false就执行传入的函数，为true就不执行，对所有goroutine可见。实现`禁止重排序 + 互斥锁`的作用（类似java中dcl + volatile的效果）。

### 7.5 goroutine与线程

#### 7.5.1 栈

> 每个操作系统都有一个固定大小的栈内存，主要用于保存函数调用期间那些`正在执行或临时暂停`的函数中的`局部变量`。

goroutine在生命周期开始的时栈大小为`2KB`，但是它的大小不是固定的，是可以按需增大和缩小，最大可达`1GB`。

#### 7.5.2 调度

> CPU通过调用`调度器`的内核函数，这个函数会暂停当前正在运行的线程，将它寄存器的信息保存到内存，查看线程列表并决定接下来运行哪一个线程，再从内存恢复线程的注册表信息，最后执行选中的线程。

go运行时包含一个自己的调度器，这个调度器使用一个`m:n`调度技术（`复用/调度m个goroutine到n个OS线程`），与内核调度器工作类似，但是go调度器只需要关心单个go程序的goroutine调度问题。

go调度器不是由硬件时钟来定期触发的，而是由特定的go语言结构来触发的，当一个goroutine调用`time.Sleep()或被通道阻塞或对互斥量操作时`，调度器就会将这个goroutine设置为休眠模式，并运行其他goroutine直到前一个可重新唤醒为止，相比内核调度器调度一个线程的成本要低得多。

> 1. go程序的`主线程`负责执行goroutine的调度工作，调度器会决定将新的goroutine放到哪个线程（processor）去执行。
> 2. 调度器会将goroutine添加到`每个线程的本地队列`中。当有线程空闲时，它会从本地队列中获取goroutine并执行它。
> 3. 如果线程的本地队列为空，processor会从`全局队列`中获取goroutine，全局队列存储所有未分配的goroutine。
> 4. 调度器会根据`抢占调度、工作窃取`等方式，在某个goroutine执行时间过长或发生阻塞时中断该goroutine的执行，也可以在某个processor队列为空时，从其他processor的队列中窃取任务执行，实现负载均衡。

#### 7.5.3 GOMAXPROCS

> GOMAXPROCS设置需要多少个OS的线程来同时执行Go代码。默认是`cpu核心数量`。
>
> 正在休眠或者被通道通信阻塞的goroutine不占用线程。

```go
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
```

#### 7.5.4 goroutine标识

goroutine和java中的线程不同，后者会有一个独特的标识（例如线程id），go不引入唯一标识的原因：主要是为了保持简洁和易用性，避免额外的开销。其次go推荐使用通道和同步安全的传递数据，也就无须关注goroutine的标识符。