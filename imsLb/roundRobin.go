package imsLb

import (
	"errors"
	"fmt"
)

const (
	SERVICE_CHANNEL_LENGTH = 256
)

type BaseRoundRobin struct{
	RRMap map[string]chan string
}

// All service info in a channel.
// When request coming, it will get a service info from the channel,
// and then put this service info to this channel auto.
func (b *BaseRoundRobin) NewRoundRobin() {
	b.RRMap = make(map[string]chan string)
	for groupName, services := range ServiceMap {
		for serviceName, hosts := range services {
			length := SERVICE_CHANNEL_LENGTH
			for len(hosts) > length {
				length = 2 * length
			}
			b.RRMap[groupName + "/" + serviceName] = make(chan string, length)

			for key, _ := range hosts {
				b.RRMap[groupName + "/" + serviceName] <- key
			}
		}
	}
}

func (b *BaseRoundRobin) RoundRobinNext(service *ServiceWithThreshold) (serviceAddr string, err error) {
	RRMutex.Lock()
	defer RRMutex.Unlock()

	if len(b.RRMap) == 0 {
		err = errors.New(fmt.Sprintf("no service named %s/%s", service.GroupName, service.ServiceName))
		return
	}

	hostChan, exist := b.RRMap[service.GroupName + "/" + service.ServiceName]
	if !exist {
		err = errors.New(fmt.Sprintf("no service named %s/%s", service.GroupName, service.ServiceName))
		return
	}

	serviceAddr = <-hostChan
	hostChan <- serviceAddr
	return serviceAddr, nil
}

func (b *BaseRoundRobin) RoundRobinCenterUpdate() {
	b.NewRoundRobin()
}

func (b *BaseRoundRobin) SetServiceThresholdValue(services ServiceWithThresholdList) (err error) {
	return nil
}