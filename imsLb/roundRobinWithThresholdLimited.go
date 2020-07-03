package imsLb

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	RRThresholdLimitedNoisyPoint uint32 = 10 // number of noisy point for one service
)

type RRWithThresholdLimitedMap map[string][]string
type ServiceWithThresholdList []*ServiceWithThreshold
type ServiceWithThreshold struct {
	GroupName   string
	ServiceName string
	ServiceHost string
	ServicePort string
	Threshold   uint8
}

type RoundRobinWithThresholdLimited struct{
	RRMap RRWithThresholdLimitedMap
	NextServiceMap map[string]int		// index of next service
}

// RoundRobin with threshold limited, traffic limits are implemented in this method.
// Traffic id limited by threshold value(threshold range is 0-10), default is 10.
// If threshold value of a service is 0, no traffic will be send to this service

// Create a global map used to maintain all services.
// It will return a map,
// the number of every services in this map is RRThresholdLimitedNoisyPoint(default is 10)
func (r *RoundRobinWithThresholdLimited) NewRoundRobin() {
	r.RRMap = make(RRWithThresholdLimitedMap)
	r.NextServiceMap = make(map[string]int)
	for groupName, services := range ServiceMap {
		for serviceName, hosts := range services {
			r.RRMap[groupName+"/"+serviceName] = make([]string, 0)
			r.NextServiceMap[groupName+"/"+serviceName] = 0
			for key, _ := range hosts {
				for i := uint32(0); i < RRThresholdLimitedNoisyPoint; i++ {
					r.RRMap[groupName+"/"+serviceName] = append(r.RRMap[groupName+"/"+serviceName], key)
				}
			}
			serviceListShuffle(r.RRMap[groupName+"/"+serviceName])
		}
	}
	fmt.Println(r.RRMap)
}

func (r *RoundRobinWithThresholdLimited) RoundRobinCenterUpdate() {
	r.NewRoundRobin()
}

func (r *RoundRobinWithThresholdLimited) RoundRobinNext(service *ServiceWithThreshold) (serviceAddr string, err error) {
	key := service.GroupName + "/" + service.ServiceName
	if _, exist := r.RRMap[key]; !exist {
		err = errors.New(fmt.Sprintf("%s dose not exist", key))
		return
	}

	RRMutex.Lock()
	defer RRMutex.Unlock()
	serviceAddr = r.RRMap[key][r.NextServiceMap[key]]
	if r.NextServiceMap[key] + 1 < len(r.RRMap[key]) {
		r.NextServiceMap[key] += 1
	} else {
		r.NextServiceMap[key] = 0
	}

	fmt.Println(r.NextServiceMap[key])
	return
}

func (r *RoundRobinWithThresholdLimited) SetServiceThresholdValue(services ServiceWithThresholdList) (err error) {
	fmt.Println("start set threshold value")
	fmt.Println(services)
	serviceNotExist := true

	for _, service := range services {
		key := service.GroupName + "/" + service.ServiceName
		if hostList, exist := r.RRMap[key]; !exist {
			continue
		} else {
			serviceNotExist = false
			serviceInfo := service.ServiceHost + ":" + service.ServicePort
			sum := serviceListSort(&hostList, serviceInfo)
			serviceNoisyPointNum := int(uint32(service.Threshold) * RRThresholdLimitedNoisyPoint / 10)
			if serviceNoisyPointNum == sum {
				continue
			}
			switch {
			case serviceNoisyPointNum < sum:
				hostList = hostList[sum-serviceNoisyPointNum:]
			default:
				tmp := make([]string, serviceNoisyPointNum-sum)
				for j := 0; j < serviceNoisyPointNum-sum; j++ {
					tmp [j] = serviceInfo
				}
				hostList = append(tmp, hostList...)
			}
			serviceListShuffle(hostList)
			RRMutex.Lock()
			r.RRMap[key] = hostList
			RRMutex.Unlock()
		}
	}

	fmt.Println(r.RRMap)

	if serviceNotExist {
		return errors.New("all service are not in RRThresholdMap")
	}
	return nil
}

// select all service info(host:port), and let them to the top of list
// Return sum of noisy point for this service
func serviceListSort(list *[]string, serviceInfo string) (serviceNum int) {
	next := 0
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == serviceInfo {
			(*list)[next], (*list)[i] = (*list)[i], (*list)[next]
			next += 1
		}
	}
	return next
}

func serviceListShuffle(list []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
}
