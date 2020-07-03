package imsRegister

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"imsPb"
	"time"
)

const (
	ETCD_ADDR            = "192.168.0.100:2379"
	REGISTER_SERVER_PORT = "27891"
	SERVER_PREFIX        = "/grpcTest/lbTest/serverCenter/"
	LEASE_TIME           = 2
)

var (
	RegisterCliv3 *clientv3.Client
)

func RegisterService(ctx context.Context, req *imsPb.ServiceRegisterRequest) {
	var err error
	if RegisterCliv3, err = clientv3.New(clientv3.Config{Endpoints: []string{"192.168.0.100:2379"}, DialTimeout: 2 * time.Second}); err != nil {
		fmt.Println("create etcd client error:", err)
		panic("create etcd client error")
	}
	defer RegisterCliv3.Close()

	serviceInfo := imsPb.GetServiceResponse{ServiceHost: req.ServiceHost, ServicePort: req.ServicePort}
	key := SERVER_PREFIX + req.ServiceGroup + "/" + req.ServiceName + "/" + req.ServiceHost + ":" + req.ServicePort
	value, _ := json.Marshal(serviceInfo)

	lease := clientv3.NewLease(RegisterCliv3)
	var (
		leaseGrantRsp *clientv3.LeaseGrantResponse
	)
	if leaseGrantRsp, err = lease.Grant(ctx, LEASE_TIME); err != nil {
		fmt.Println("etcd lease grant error:", err)
		panic("etcd lease grant error")
	}
	if _, err = lease.KeepAlive(ctx, leaseGrantRsp.ID); err != nil {
		fmt.Println("lease keep alive auto error:", err)
		panic("lease keep alive auto error")
	}

	if _, err = RegisterCliv3.Put(context.TODO(), key, string(value), clientv3.WithLease(leaseGrantRsp.ID)); err != nil {
		fmt.Println("put service info to etcd error:", err)
		panic("put service info to etcd error")
	}

	select {
	case m := <-ctx.Done():
		fmt.Println("register done", m)
	}
}

//type RegisterServer struct{}
//
//func (rs *RegisterServer) ServiceRegister(ctx context.Context, req *imsPb.ServiceRegisterRequest) (rsp *imsPb.ServiceRegisterResponse, err error) {
//	rsp = new(imsPb.ServiceRegisterResponse)
//
//	return
//}

//func StartRegisterServer() {
//	// rpc server of register init
//	listener, err := net.Listen("tcp", "localhost:"+REGISTER_SERVER_PORT)
//	if err != nil {
//		fmt.Printf("register server listen to port %s error:%v \n", REGISTER_SERVER_PORT, err)
//		panic("register server start error")
//	}
//
//	s := grpc.NewServer()
//	imsPb.RegisterImsRegisterServer(s, &RegisterServer{})
//	reflection.Register(s)
//
//	err = s.Serve(listener)
//	if err != nil {
//		fmt.Println("register grpc server error:", err)
//		panic("register grpc server error")
//	}
//
//	// etcd client init
//	if RegisterCliv3, err = clientv3.New(clientv3.Config{
//		Endpoints:   []string{ETCD_ADDR},
//		DialTimeout: 2 * time.Second,
//	}); err != nil {
//		fmt.Println("create etcd cli error:", err)
//		panic("create etcd cli error")
//	}
//
//	listener.Accept()
//}
