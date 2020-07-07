package imsLb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"imsPb"
	"middleWare"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	LB_PORT       = "27890"
	ETCD_ADDR     = "192.168.0.100:2379"
	SERVER_PREFIX = "/grpcTest/lbTest/serverCenter/"
	SERVER_PORT   = 33333
)

var (
	RRMutex sync.Mutex
)

type RoundRobinLoadBalance interface {
	// round robin scheduler init
	NewRoundRobin()

	// get a service address
	RoundRobinNext(service *ServiceWithThreshold) (serviceAddr string, err error)

	// update scheduler center
	RoundRobinCenterUpdate()

	// set threshold value to
	SetServiceThresholdValue(services ServiceWithThresholdList) (err error)
}

var (
	ServiceMap                     map[string]SubService // {"groupName": SubService}
	LbCliV3                        *clientv3.Client
	LoadBalanceSchedulingAlgorithm = map[string]RoundRobinLoadBalance{
		"BaseRoundRobin":                 new(BaseRoundRobin),
		"RoundRobinWithThresholdLimited": new(RoundRobinWithThresholdLimited),
	}
)

type SubService map[string]map[string]*imsPb.GetServiceResponse // {"serviceName": {"host:port": *imsPb.GetServiceResponse}}

type loadBalance struct {
	schedulingAlgorithm RoundRobinLoadBalance
}

func (lb *loadBalance) GetServiceWithRR(ctx context.Context, req *imsPb.GetServiceRequest) (rsp *imsPb.GetServiceResponse, err error) {

	var (
		serviceHost string
	)

	// TODO: create token if no token, if there is token in ctx, check token

	if serviceHost, err = lb.schedulingAlgorithm.RoundRobinNext(&ServiceWithThreshold{GroupName: req.ServiceGroup, ServiceName: req.ServiceName}); err != nil {
		return
	}
	return ServiceMap[req.ServiceGroup][req.ServiceName][serviceHost], nil
}

func (lb *loadBalance) SetServiceTrafficLimit(ctx context.Context, req *imsPb.ServiceTrafficLimitRequest) (rsp *imsPb.ServiceTrafficLimitResponse, err error) {
	rsp = new(imsPb.ServiceTrafficLimitResponse)
	if err = lb.schedulingAlgorithm.SetServiceThresholdValue(ServiceWithThresholdList{
		&ServiceWithThreshold{
			GroupName:   req.ServiceGroup,
			ServiceName: req.ServiceName,
			ServiceHost: req.ServiceHost,
			ServicePort: req.ServicePort,
			Threshold:   uint8(req.TrafficThreshold),
		}}); err != nil {
		fmt.Println("set threshold value error:", err)
		rsp.Message = "set threshold value error:"
	}

	return
}

func StartLoadBalanceServer(loadBalanceSchedulingAlgorithm string) {
	fmt.Println("LB server start")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("LB server panic:", err)
		}
		fmt.Println("LB server stop")
	}()

	algo, exist := LoadBalanceSchedulingAlgorithm[loadBalanceSchedulingAlgorithm]
	if !exist {
		fmt.Println("LoadBalance scheduling algorithm error")
		panic("LoadBalance scheduling algorithm error")
	}

	var (
		err      error
		listener net.Listener
	)

	// etcd client init
	if LbCliV3, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCD_ADDR},
		DialTimeout: 2 * time.Second,
	}); err != nil {
		fmt.Println("create etcd cli error:", err)
		panic("create etcd cli error")
	}

	// global service map init
	if err = serviceMapInit(SERVER_PREFIX); err != nil {
		panic("serviceMapInit error")
	}

	// load balance init
	algo.NewRoundRobin()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go serviceListWatcher(ctx, algo)

	// rpc server of lb init
	listener, err = net.Listen("tcp", "localhost:"+LB_PORT)
	if err != nil {
		fmt.Printf("lb server listen to port %s error:%v \n", LB_PORT, err)
		panic("lb server start error")
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(middleWare.TokenInterceptor))
	imsPb.RegisterImsLoadBalanceServer(s, &loadBalance{schedulingAlgorithm: algo})
	reflection.Register(s)

	err = s.Serve(listener)
	if err != nil {
		fmt.Println("lb grpc server error:", err)
		panic("lb grpc server error")
	}

	listener.Accept()
}

func serviceMapInit(serviceName string) (err error) {
	var (
		getRsp *clientv3.GetResponse
	)
	if getRsp, err = LbCliV3.Get(context.TODO(), serviceName, clientv3.WithPrefix()); err != nil {
		fmt.Println("get service form etcd error:", err)
		return
	} else {
		ServiceMap = make(map[string]SubService)
		for _, v := range getRsp.Kvs {
			serviceGroup, serviceName, serviceAddr, err1 := serviceKeySplit(string(v.Key))
			if err1 != nil {
				return err1
			}
			serviceGroupAppend(serviceGroup)
			service := new(imsPb.GetServiceResponse)
			if err1 = json.Unmarshal(v.Value, service); err1 != nil {
				return err1
			}
			serviceAppend(serviceName, serviceAddr, service, ServiceMap[serviceGroup])
		}
		return nil
	}
}

func serviceKeySplit(key string) (serviceGroup, serviceName, serviceAddr string, err error) {
	ss := strings.Split(key, "/")
	if len(ss) <= 3 {
		fmt.Println("service key error:", key)
		err = errors.New("key error")
		return
	}
	serviceGroup, serviceName, serviceAddr = ss[len(ss)-3], ss[len(ss)-2], ss[len(ss)-1]
	return
}

func serviceGroupAppend(group string) {
	if _, exist := ServiceMap[group]; !exist {
		ServiceMap[group] = make(map[string]map[string]*imsPb.GetServiceResponse)
	}
}

func serviceAppend(name, addr string, rs *imsPb.GetServiceResponse, services SubService) {
	if _, exist := services[name]; !exist {
		services[name] = make(map[string]*imsPb.GetServiceResponse)
	}
	services[name][addr] = rs
}

func isGroup(key []byte) (bool, error) {
	wholeName := strings.Split(string(key), "/")
	switch len(wholeName) - len(strings.Split(SERVER_PREFIX, "/")) {
	case 1:
		return true, nil
	case 2:
		return false, nil
	default:
		return false, errors.New("format of key error")
	}
}

func serviceListWatcher(ctx context.Context, rr RoundRobinLoadBalance) {
	fmt.Println("service watcher start")
	var (
		watchRsp clientv3.WatchResponse
		err      error
	)
	serviceWatcherChan := LbCliV3.Watch(ctx, SERVER_PREFIX, clientv3.WithPrefix())
	for {
		watchRsp = <-serviceWatcherChan
		for _, event := range watchRsp.Events {
			switch {
			case event.IsCreate():
				if err = eventCreateHandler(event.Kv.Key, event.Kv.Value); err != nil {
					fmt.Println(err)
				}
			case event.IsModify():
				if err = eventModifyHandler(event.Kv.Key, event.Kv.Value); err != nil {
					fmt.Println(err)
				}
			default:
				if err = eventDeleteHandler(event.Kv.Key); err != nil {
					fmt.Println(err)
				}
			}

			RRMutex.Lock()
			rr.RoundRobinCenterUpdate()
			RRMutex.Unlock()
		}
		fmt.Println(ServiceMap)
	}
}

func eventCreateHandler(k, v []byte) error {
	groupName, serviceName, serviceAddr, err := serviceKeySplit(string(k))
	if err != nil {
		return err
	}
	serviceGroupAppend(groupName)
	services, _ := ServiceMap[groupName]
	serviceInfo := new(imsPb.GetServiceResponse)
	_ = json.Unmarshal(v, serviceInfo)
	serviceAppend(serviceName, serviceAddr, serviceInfo, services)
	return nil
}

func eventModifyHandler(k, v []byte) error {
	isGroupKey, err := isGroup(k)
	if err != nil {
		return err
	}
	if isGroupKey {
		return errors.New(fmt.Sprintf("key error: %s is group, service group can not be modified", string(k)))
	}

	groupName, serviceName, serviceAddr, err1 := serviceKeySplit(string(k))
	if err1 != nil {
		return err1
	}
	services, exist := ServiceMap[groupName]
	if !exist {
		return errors.New(fmt.Sprintf("service group[%s] dose not exist", groupName))
	}
	_, exist1 := services[serviceName]
	if !exist1 {
		return errors.New(fmt.Sprintf("service[%s] dose not exist", serviceName))
	}

	serviceInfo := new(imsPb.GetServiceResponse)
	_ = json.Unmarshal(v, serviceInfo)
	ServiceMap[groupName][serviceName][serviceAddr] = serviceInfo
	return nil
}

func eventDeleteHandler(k []byte) error {
	isGroupKey, err := isGroup(k)
	if err != nil {
		return err
	}

	if isGroupKey {
		delete(ServiceMap, string(k))
	} else {
		groupName, serviceName, serviceAddr, err1 := serviceKeySplit(string(k))
		if err1 != nil {
			return err1
		}

		services, exist := ServiceMap[groupName]
		if !exist {
			return errors.New(fmt.Sprintf("service group[%s] dose not exist", groupName))
		}

		_, exist = services[serviceName]
		if !exist {
			return errors.New(fmt.Sprintf("service[%s] dose not exist", serviceName))
		}

		delete(services[serviceName], serviceAddr)
	}
	return nil
}
