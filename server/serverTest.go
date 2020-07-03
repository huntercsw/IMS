package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"imsPb"
	"imsRegister"
	"os"
	"os/signal"
	"time"
)

var (
	group = "imsTest"
	name = "yinuo"
	host = "192.168.1.151"
	port = "65432"
	res = make([]string, 0)
)

func main() {
	fmt.Println("service start")
	ctx, cancel := context.WithCancel(context.TODO())
	defer func() {
		cancel()
		fmt.Println("service stop")
	}()

	go imsRegister.RegisterService(ctx, &imsPb.ServiceRegisterRequest{
		ServiceGroup:group, ServiceName:name, ServiceHost:host, ServicePort:port,
	})


	time.Sleep(60*time.Second)
	for i := 0; i < 100; i++ {
		test()
	}

	n := 0
	for _, ip := range res {
		if ip == "192.168.0.100:65432" {
			n += 1
		}
	}
	fmt.Println("192.168.0.100:65432----", n)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println(s)
}

func test() {
	conn, err := grpc.Dial(":27890", grpc.WithInsecure())
	if err != nil {
		fmt.Println("connect to server error:", err)
		return
	}

	defer conn.Close()

	c := imsPb.NewImsLoadBalanceClient(conn)

	r, err1 := c.GetServiceWithRR(context.Background(), &imsPb.GetServiceRequest{ServiceGroup:"imsTest", ServiceName:"yinuo"})
	if err1 != nil {
		fmt.Println("call grpc error:", err)
		return
	}
	fmt.Println(r.ServiceHost, r.ServicePort)
	res = append(res, r.ServiceHost + ":" + r.ServicePort)
}
