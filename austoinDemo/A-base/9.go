package mybase

import (
	"fmt"
	"sync"
	"time"
)

func Lesson9() {
	goroutineDemo() // 启动并发任务
	channelDemo()   //理解协议间如何通信
	waitGroupDemo()    //等待协程完成
	mutexDemo()        //处理共享资源
	selectDemo()       //综合运用，超时控制
}

func goroutineDemo() {
	fmt.Println("主程序开始")

	// 启动三个协程
	go say("Hello")
	go say("Hello, Austoin")
	go say("Hello, Austoin, GoRuntine")

	// 主程序需要等待，否者协程没开始执行主程序就退出了
	time.Sleep(100 * time.Millisecond)

	// 匿名协程：不需要先定义函数
	go func() {
		fmt.Println("匿名协程执行")
	}()

	time.Sleep(50 * time.Millisecond)
}

// 普通函数作为协程
func say(s string) {
	fmt.Println("协程输出：", s)
}

func channelDemo() {
	// 1. 无缓冲 Channel
	ch1 := make(chan string)

	// 启动协程发送数据
	go func() {
		// 这里会阻塞，直到主程序准备接收
		ch1 <- "来自 GoRontine 的消息 (无缓冲 Channel)"
		fmt.Println("发送完成")
	}()

	// 接收数据（这里会阻塞，直到协程发送）
	msg := <-ch1
	fmt.Println("收到消息:", msg, "\n")

	// 2. 有缓冲 Channel
	// 特点：可以暂存数据，缓冲区满时才阻塞
	ch2 := make(chan int, 3) // 容量为3

	// 发送数据，数据未满不阻塞
	ch2 <- 1
	ch2 <- 2
	ch2 <- 3
	// ch2 <- 4 // 阻塞！因为缓冲满了

	// 接收数据
	fmt.Println("有缓冲 Channel 接收:", <- ch2)
	fmt.Println("有缓冲 Channel 接收:", <- ch2)
	fmt.Println("有缓冲 Channel 接收:", <- ch2, "\n")

	// 3. 关闭 Channel + range 遍历
	// 发送方关闭Channel后，接收方仍然可以读取剩余数据
	ch3 := make(chan int, 5)
	go func ()  {
		for i := 1; i <= 5; i++ {
			ch3 <- i
		}	
		close(ch3) // 关闭Channel
// 关闭通道（close(ch)）的本质是告诉接收方不会再有新数据发送到这个通道了
// 而不是 禁止读取这个通道。
	}()

	// 使用 range 遍历Channel（会自动等待数据，直到Channel关闭）
	fmt.Println("range 遍历")
	for v := range ch3{
		fmt.Println("ch3 的数据:",v)
	}
	fmt.Println()

// for v := range ch3 是 Go 为通道接收数据设计的便捷循环语法；
// 它完全等价于 “用<-显式接收数据 + 判读通道状态 + 循环退出” 的组合逻辑。
}

func waitGroupDemo(){
	var wg sync.WaitGroup
	
	// 启动三个协程
	for i := 1; i <= 3; i++{
		wg.Add(1) // 等待计数 +1

		// 启动协程，传入参数 i
		go func (id int)  {
			defer wg.Done() 	// 协程结束时，计数 -1

// defer作用是将其后的语句延迟到当前函数
// 或匿名函数执行完毕时（无论正常结束还是异常终止）再执行。

			fmt.Printf("Goroutine %d 开始\n", id)
			time.Sleep(50 * time.Millisecond) // 模拟耗时
			fmt.Printf("Goroutine %d 结束\n", id)
		}(i)
	}

	wg.Wait() //等待所有协程完成
	fmt.Println("所有 Goroutine 完成", "\n")
}

func mutexDemo(){
	var (
	  counter int   // 共享变量
	  mutex   sync.Mutex // 互斥锁
	)

	// 启动 1000 个协程
	for i := 0; i < 1000; i++ {
		go func ()  {
			mutex.Lock() // 加锁（其他协程会阻塞等待）
			counter ++ // 安全地访问共享变量
			mutex.Unlock() // 解锁
		}()
	}

	// 等待所有协程完成
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("共享变量 Counter: %d (应为1000)\n\n", counter)
}

func selectDemo(){
	ch5 := make(chan int)
	ch6 := make(chan int)

	// 两个协程，分别向两个Channel(通道)发送数据
	go func ()  {
		time.Sleep(50 * time.Millisecond)
		ch5 <- 100
	}()

	go func ()  {
		time.Sleep(30 * time.Millisecond)
		ch6 <- 200
	}()

	// select 等待任意一个Channel就绪
	// 多通道监听，谁先有数据先执行谁
	select {
	case v1 := <- ch5:
		fmt.Printf("从 ch5 收到: %d\n", v1)
	case v2 := <- ch6:
		fmt.Printf("从 ch6 收到: %d\n", v2)
	}

	// 超时控制：time.After 返回超时Channel
// time.After：返回一个通道，100ms后会向该通道发送当前时间（time.Time类型）
	timeout := time.After(100 * time.Millisecond)
	done := make(chan bool)

	go func ()  {
		// 模拟耗时操作（200ms）
		time.Sleep(200 * time.Millisecond)
		done <- true
	}()

	select {
	case <- done:
		fmt.Println("操作完成")
	case <-timeout:
		fmt.Println("操作超时（预期）")
	}
}