package main

import (
	"fmt"
	"time"
)

/*
概念释义:
0. 4核8线程:8线程指的是指 `intel超线程技术` 操作系统识别的是8个核心,实际只有4个物理CPU能并行执行.
> 0.1 不是真正的物理核心：共享核心的算术逻辑单元(ALU)、缓存等资源
> 0.2 硬件级线程调度：核心内部维护两套线程状态（寄存器组等）
> 0.3 提高资源利用率：当某个线程停顿（如等待数据）时，立即切换执行另一个线程

1. 进程: 应用程序的执行副本, 向操作系统申请资源(硬件资源,文件资源,内存分页,CPU资源等)的最小单位.所有的进程启动后都会有一个`主线程(main thread)`
2. 线程: 是操作系统执执行的最小单位(只依赖CPU和内存).进程和线程没必要将他们存储为一个树状结构,可以分属不同的表,通过`id`进行关联.
> 2.1 应用程序启动后将变成进程,进程启动后拿到 主线程 ,主线程再去到操作系统申请线程,注册到操作系统线程表中,最后由CPU来按照 时间分片 执行
> 2.2 <B>程序语言的线程</B> 与 <B>操作系统的线程</B> 是映射关系,只有操作系统的线程是<B>Real Thread(Kernel Thread)</B>, Java的线程与操作系统的线程是1:1的关系

3. 协程(routine): 又叫做轻量级线程,不在内核态,在用户态空间实现,不依赖操作系统内核.协程的上下文切换开销远小于线程的切换(线程的切换涉及到中断,保存,加载等).
4. goroutine: 是Go语言实现的协程,初始只有2kb(可动态增长), 多个goroutine映射到少量操作系统线程(M:N)
> 4.1 G-P-M 模型
   +---+    +---+    +---+
   | M |    | M |    | M |  (操作系统线程)
   +---+    +---+    +---+
     |        |        |
     v        v        v
   +---+    +---+    +---+
   | P |    | P |    | P |  (逻辑处理器)
   +---+    +---+    +---+
    / \      / \      / \
   v   v    v   v    v   v
 +---+---+ +---+---+ +---+---+
 | G | G | | G | G | | G | G |  (goroutines)
 +---+---+ +---+---+ +---+---+
逻辑处理器的数量由 GOMAXPROCS 指定,默认是 逻辑 CPU 核心数
*/

func multiGoroutine() {
	go func() {
		defer fmt.Println("goroutine A exit ...")
		func() {
			defer fmt.Println("goroutine B exit ...")
			// return or runtime.Goexit()
			// runtime.Goexit() // 退出调用它的协程
			fmt.Println("goroutine B exec")
		}()
		fmt.Println("goroutine A exec")
	}() // 定义并调用
	for { // main退出会导致子协程推出,和java相同
		time.Sleep(1 * time.Second)
	}
}

func main() {
	multiGoroutine()
}
