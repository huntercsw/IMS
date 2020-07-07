package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"imsPb"
	"imsRegister"
	"middleWare"
	"os"
	"os/signal"
	"time"
)

var (
	group = "imsTest"
	name  = "yinuo"
	host  = "192.168.1.151"
	port  = "65432"
	res   = make([]string, 0)
)

func main() {
	//fmt.Println("service start")
	ctx, cancel := context.WithCancel(context.TODO())
	defer func() {
		cancel()
		fmt.Println("service stop")
	}()

	go imsRegister.RegisterService(ctx, &imsPb.ServiceRegisterRequest{
		ServiceGroup: group, ServiceName: name, ServiceHost: host, ServicePort: port,
	})

	time.Sleep(1 * time.Second)

	start := time.Now().UnixNano()

	for i := 0; i < 2000; i++ {
		token, err := getToken()
		if err != nil {
			fmt.Println("get token error:", err)
		}

		//fmt.Println("start test")
		test(token)
	}

	fmt.Println(time.Now().UnixNano() - start)

	//n := 0
	//for _, ip := range res {
	//	if ip == "192.168.0.100:65432" {
	//		n += 1
	//	}
	//}
	//fmt.Println("192.168.0.100:65432----", n)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	fmt.Println(s)
}

func test(token *imsPb.ImsToken) {
	conn, err := grpc.Dial(":27890", grpc.WithInsecure(), grpc.WithPerRPCCredentials(&middleWare.AuthenticationInToken{Token: token}))
	if err != nil {
		fmt.Println("connect to server error:", err)
		return
	}

	defer conn.Close()

	c := imsPb.NewImsLoadBalanceClient(conn)

	r, err1 := c.GetServiceWithRR(context.Background(), &imsPb.GetServiceRequest{ServiceGroup: "imsTest", ServiceName: "yinuo"})
	if err1 != nil {
		fmt.Println("call grpc error:", err1)
		return
	}
	res = append(res, r.ServiceHost+":"+r.ServicePort)
}

func getToken() (token *imsPb.ImsToken, err error) {
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":"+middleWare.TICKET_OFFICE_PORT, grpc.WithInsecure())
	if err != nil {
		fmt.Println("connect to server error:", err)
		return
	}

	defer conn.Close()

	c := imsPb.NewImsTicketOfficeClient(conn)

	token, err = c.GetToken(context.Background(), &imsPb.GetTokenRequest{DomainName: "yinuo.com.cn", ServiceName: "yinuo"})
	if err != nil {
		fmt.Println("call grpc error:", err)
		return
	}

	return
}
