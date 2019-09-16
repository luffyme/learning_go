package main

import (
	"fmt"
	"time"
	"runtime"
)

func main() {
	//1.协程的启动
	//Go 语言里创建一个协程非常简单，使用 go 关键词加上一个函数调用就可以了
	//Go 语言会启动一个新的协程，函数调用将成为这个协程的入口。
	//Go里面协程之间不存在那么多层级关系，只有一个主协程，其他的都是子协程，子协程之间是平行关系。
	fmt.Println("run in main goroutine")
	go func() {
		fmt.Println("run in child goroutine")
		go func() {
			fmt.Println("run in grand child goroutine")
			go func() {
				fmt.Println("run in grand grand child goroutine")
			}()
		}()
	}()
	time.Sleep(time.Second)
	fmt.Println("main goroutine will quit")
	
	//2.子协程的异常退出
	//在使用子协程时一定要特别注意保护好每个子协程，确保它们正常安全的运行。
	//因为子协程的异常退出会将异常传播到主协程，直接会导致主协程也跟着挂掉，然后整个程序就崩溃了。
	
	fmt.Println("run in main goroutine")
	go func() {
		fmt.Println("run in child goroutine")
		//panic("wtf")
	}()
	time.Sleep(time.Second)
	fmt.Println("main goroutine will quit")
	
	//3.协程的本质
	//一个进程内部可以运行多个线程，而每个线程又可以运行很多线程。
	//线程要负责对协程进行调度，保证每个协程都有机会得到执行。
	//当一个协程睡眠的时候，它要将线程的运行权让给其他协程来运行，而不能持续霸占这个线程，同一个线程内部最多只会有一个协程正在运行。

	//线程的调度是由操作系统进行调度的，调度算法运行在内核态，而协程的调度是由Go语言的运行时负责，调度算法在用户态。
	
	//协程可以简化为三个状态，运行态，就绪态，休眠态。
	//同一个线程中最多只会存在一个处于运行态的协程。
	//就绪态的协程是指哪些具备了运行能力但是还没有得到运行机会的协程，他们随时都会被调度到运行态。
	//休眠态的协程还不具备运行能力，他们在等待某些条件的发生，比如IO操作的完成， 休眠时间的结束等。

	//操作系统对线程的调度是抢占式的，也就是说，如果单个线程的死循环不会影响其他线程的执行，每个线程的运行收到时间片的限制。
	//Go语言运行时对协程的调度并不是抢占式的，如果单个协程通过死循环霸占了线程的执行权，那么这个线程就没有机会取运行其他协程了，这个线程就假死了。
	//不过一个进程内部往往有多个线程，假死了一个没事，全部假死了才会导致整个进程卡死。

	//每个线程都会包含多个就绪态态的协程形成了一个就绪队列，如果这个线程因为某个协程死循环导致假死，那么队列上所有的就绪态协程是不是就没有机会得到运行了呢？
	//Go 语言运行时调度器采用了 work-stealing 算法，当某个线程空闲时，也就是该线程上所有的协程都在休眠（或者一个协程都没有），它就会去其它线程的就绪队列上去偷一些协程来运行。

	//4.设置线程数
	//Go 运行时会将线程数会被设置为机器 CPU 逻辑核心数。
	//runtime 包提供了 GOMAXPROCS(n int) 函数允许我们动态调整线程数，如果参数 n<=0，会返回修改前的线程数
	
	// 读取默认的线程数
	fmt.Println(runtime.GOMAXPROCS(0))
	// 设置线程数为 10
	runtime.GOMAXPROCS(10)
	// 读取当前的线程数
	fmt.Println(runtime.GOMAXPROCS(0))

	//获取当前的协程数量可以使用 runtime 包提供的 NumGoroutine() 方法
	fmt.Println(runtime.NumGoroutine())
	for i:=0;i<10;i++ {
		go func(){
			for {
				time.Sleep(time.Second)
			}
		}()
	}
	fmt.Println(runtime.NumGoroutine())
	
	//5.协程的应用
	//在 HTTP API 应用中，每一个 HTTP 请求，服务器都会单独开辟一个协程来处理。
	//在这个请求处理过程中，要进行很多 IO 调用，比如访问数据库、访问缓存、调用外部系统等，协程会休眠，IO 处理完成后协程又会再次被调度运行。待请求的响应回复完毕后，链接断开，这个协程的寿命也就到此结束。

	//在消息推送系统中，客户端的链接寿命很长，大部分时间这个链接都是空闲状态，客户端会每隔几十秒周期性使用心跳来告知服务器你不要断开我。
	//在服务器端，每一个来自客户端链接的维持都需要单独一个协程。
	//因为消息推送系统维持的链接普遍很闲，单台服务器往往可以轻松撑起百万链接，这些维持链接的协程只有在推送消息或者心跳消息到来时才会变成就绪态被调度运行。

	//聊天系统也是长链接系统，它内部来往的消息要比消息推送系统频繁很多，限于 CPU 和 网卡的压力，它能撑住的连接数要比推送系统少很多。
	//不过原理是类似的，都是一个链接由一个协程长期维持，连接断开协程也就消亡。
}