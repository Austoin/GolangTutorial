package main

// 本文件演示 Go 语言的并发编程：Goroutine 和 Channel

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ========== Goroutine ==========

	// 1. Goroutine 基础
	// 使用 go 关键字创建 Goroutine（轻量级线程）
	// Goroutine 由 Go 运行时管理，非常轻量（初始栈约 2KB）

	fmt.Println("主程序开始")

	// 启动 Goroutine
	// go 关键字后跟函数调用
	go say("Hello") // 并发执行
	go say("Concurrent")
	go say("World")

	// 主程序等待一段时间，让 Goroutine 有机会执行
	// 实际开发中应使用 sync.WaitGroup
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n主程序结束")

	// 2. 匿名 Goroutine
	// 可以直接启动匿名函数作为 Goroutine
	go func() {
		fmt.Println("匿名 Goroutine")
	}()

	time.Sleep(50 * time.Millisecond)

	// ========== Channel ==========

	// 3. Channel 基础
	// Channel 用于 Goroutine 之间传递数据
	// 声明：var ch chan 类型
	// 使用 make 创建：ch := make(chan 类型)

	// 创建无缓冲 Channel
	// 无缓冲 Channel：发送和接收必须同时进行
	ch1 := make(chan string)

	// 启动 Goroutine 发送数据
	go func() {
		ch1 <- "来自 Goroutine 的消息" // 发送数据到 Channel
	}()

	// 从 Channel 接收数据
	// 接收会阻塞直到有数据可用
	msg := <-ch1
	fmt.Printf("收到消息: %s\n", msg)

	// 4. 有缓冲 Channel
	// 创建有缓冲的 Channel，容量为 3
	ch2 := make(chan int, 3)

	// 发送数据到有缓冲 Channel
	ch2 <- 1 // 不会阻塞，因为还有空间
	ch2 <- 2
	ch2 <- 3
	// ch2 <- 4  // 会阻塞，因为缓冲区已满

	// 接收数据
	fmt.Printf("ch2 接收: %d\n", <-ch2)
	fmt.Printf("ch2 接收: %d\n", <-ch2)
	fmt.Printf("ch2 接收: %d\n", <-ch2)

	// 5. Channel 方向
	// 可以限制 Channel 的方向（只发送或只接收）
	// 接收端：<-chan Type
	// 发送端：chan<- Type

	// 创建发送端 Channel
	sendCh := make(chan int, 2)

	// 启动发送协程
	go sender(sendCh)

	// 接收数据
	for i := 0; i < 2; i++ {
		fmt.Printf("主程序接收: %d\n", <-sendCh)
	}

	// ========== 关闭 Channel ==========

	// 6. 关闭 Channel
	// 发送方可以关闭 Channel
	// 接收方可以通过第二个返回值判断 Channel 是否已关闭
	ch3 := make(chan int, 5)

	// 启动协程发送数据
	go func() {
		for i := 1; i <= 5; i++ {
			ch3 <- i
		}
		close(ch3) // 关闭 Channel
	}()

	// 接收数据
	for {
		// v 是接收到的值
		// ok 是布尔值，true 表示 Channel 未关闭，false 表示已关闭
		v, ok := <-ch3
		if !ok {
			fmt.Println("Channel 已关闭")
			break
		}
		fmt.Printf("接收: %d\n", v)
	}

	// 7. range 遍历 Channel
	// 更加简洁的方式遍历 Channel
	ch4 := make(chan int, 3)
	go func() {
		ch4 <- 10
		ch4 <- 20
		ch4 <- 30
		close(ch4)
	}()

	// 使用 range 遍历，收到 close 信号后自动退出
	for v := range ch4 {
		fmt.Printf("range 接收: %d\n", v)
	}

	// ========== sync 包 ==========

	// 8. WaitGroup
	// 等待一组 Goroutine 完成
	var wg sync.WaitGroup

	// 启动 3 个 Goroutine
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 增加等待计数
		go func(id int) {
			defer wg.Done() // 完成后减少计数
			fmt.Printf("Goroutine %d 开始\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Goroutine %d 结束\n", id)
		}(i)
	}

	wg.Wait() // 等待所有 Goroutine 完成
	fmt.Println("所有 Goroutine 完成")

	// 9. Mutex（互斥锁）
	// 保护共享资源的访问
	var (
		counter int
		mutex   sync.Mutex
	)

	// 启动多个 Goroutine 修改 counter
	for i := 0; i < 1000; i++ {
		go func() {
			mutex.Lock()   // 加锁
			counter++      // 访问共享资源
			mutex.Unlock() // 解锁
		}()
	}

	// 等待一段时间让所有 Goroutine 完成
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Counter: %d\n", counter) // 应该是 1000

	// 10. RWMutex（读写锁）
	// 适合读多写少的场景
	var rwMutex sync.RWMutex
	data := 0

	// 读操作
	for i := 0; i < 5; i++ {
		go func() {
			rwMutex.RLock() // 读锁
			fmt.Printf("读操作: %d\n", data)
			time.Sleep(10 * time.Millisecond)
			rwMutex.RUnlock()
		}()
	}

	// 写操作
	for i := 0; i < 2; i++ {
		go func(id int) {
			rwMutex.Lock() // 写锁
			data = id
			fmt.Printf("写操作: %d\n", data)
			time.Sleep(10 * time.Millisecond)
			rwMutex.Unlock()
		}(i * 10)
	}

	time.Sleep(100 * time.Millisecond)

	// ========== Select ==========

	// 11. Select 语句
	// 多路复用，可以同时等待多个 Channel
	ch5 := make(chan int)
	ch6 := make(chan int)

	// 启动两个协程，分别向两个 Channel 发送数据
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch5 <- 100
	}()

	go func() {
		time.Sleep(30 * time.Millisecond)
		ch6 <- 200
	}()

	// 使用 select 等待任意一个 Channel 就绪
	select {
	case v1 := <-ch5:
		fmt.Printf("从 ch5 收到: %d\n", v1)
	case v2 := <-ch6:
		fmt.Printf("从 ch6 收到: %d\n", v2)
		// default:
		//     fmt.Println("都没有数据")
	}

	// 12. 超时处理
	// 使用 select 实现超时
	timeout := time.After(100 * time.Millisecond)
	done := make(chan bool)

	go func() {
		time.Sleep(200 * time.Millisecond) // 模拟耗时操作
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("操作完成")
	case <-timeout:
		fmt.Println("操作超时")
	}

	// 13. 死锁检测
	// Go 会在运行时检测死锁
	// chDeadlock := make(chan int)
	// <-chDeadlock  // 没有发送方，会死锁
}

// ========== 辅助函数 ==========

// say 打印消息
func say(s string) {
	fmt.Println(s)
}

// sender 发送数据到 Channel
// 接收只读 Channel
func sender(ch chan<- int) {
	ch <- 10
	ch <- 20
}

// ========== 总结 ==========
// 1. Goroutine 使用 go 关键字创建
// 2. Channel 用于 Goroutine 通信
// 3. 无缓冲 Channel 需要发送和接收同时进行
// 4. 有缓冲 Channel 可以暂存数据
// 5. 关闭 Channel 表示不再发送数据
// 6. WaitGroup 等待一组 Goroutine
// 7. Mutex 互斥锁保护共享资源
// 8. RWMutex 读写锁适合读多写少
// 9. Select 多路复用等待多个 Channel
// 10. Channel 是 Go 并发的核心
