package consul

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// ServiceData defines wrapped service data from consul
type ServiceData struct {
	ServiceName string
	ServiceID   string
	Port        string
	Addr        string
}

var gServiceDataMap = map[string]*ServiceData{}                   // serviceID : serviceData
var gServiceDataNameIndex = map[string](map[string]interface{}){} // serviceName: serviceIDs
var gServiceDataMutex = &sync.RWMutex{}

// get all know services by Catalog.Services.
func collectAllServices() {
	catalog := gConsulClient.Catalog()
	services, _, err := catalog.Services(nil)
	if err != nil {
		logrus.Error("query catalog.services failed.", err)
		return
	}
	for name := range services {
		collectServiceDatas(name)
	}
	log := logrus.WithFields(logrus.Fields{})
	for serviceID, data := range gServiceDataMap {
		log = log.WithField(serviceID, data)
	}
	log.Debug("collect all services")
}

func collectHealthServiceIDs(serviceName string) ([]string, error) {
	health := gConsulClient.Health()
	entries, _, err := health.Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("query health.service failed:%v", err)
	}
	passingIDs := []string{}
	for _, entry := range entries {
		passingIDs = append(passingIDs, entry.Service.ID)
	}
	sort.Strings(passingIDs)
	return passingIDs, nil
}

func removeUnhealthServices(serviceName string, passingIDs []string) {
	gServiceDataMutex.Lock()
	defer gServiceDataMutex.Unlock()

	services, exists := gServiceDataNameIndex[serviceName]
	if !exists {
		return
	}
	for serviceID := range services {
		if x := sort.SearchStrings(passingIDs, serviceID); x >= 0 {
			continue
		}
		log := logrus.WithFields(logrus.Fields{
			"service_name": serviceName,
			"service_ID":   serviceID,
		})
		log.Warn("service not in health status")
		delete(gServiceDataMap, serviceID)
		delete(services, serviceID)
	}
}

func collectServiceDatas(serviceName string) {
	log := logrus.WithField("service", serviceName)
	passingIDs, err := collectHealthServiceIDs(serviceName)
	if err != nil {
		log.Error(err)
		return
	}
	removeUnhealthServices(serviceName, passingIDs)
	catalog := gConsulClient.Catalog()
	ss, _, err := catalog.Service(serviceName, "", nil)
	if err != nil {
		log.Error("call catalog.Service failed", err)
		return
	}

	gServiceDataMutex.Lock()
	defer gServiceDataMutex.Unlock()
	for _, s := range ss {
		gServiceDataMap[s.ServiceID] = &ServiceData{
			ServiceName: serviceName,
			ServiceID:   s.ServiceID,
			Port:        strconv.Itoa(s.ServicePort),
			Addr:        s.ServiceAddress,
		}
		if _, found := gServiceDataNameIndex[serviceName]; !found {
			gServiceDataNameIndex[serviceName] = map[string]interface{}{}
		}
		gServiceDataNameIndex[serviceName][s.ServiceID] = nil
	}
}

func GetServiceDataByID(serviceID string) *ServiceData {
	gServiceDataMutex.RLock()
	defer gServiceDataMutex.RUnlock()
	if sd, found := gServiceDataMap[serviceID]; found {
		return sd
	}
	return nil
}

func GetServiceDatasByName(serviceName string) []ServiceData {
	gServiceDataMutex.RLock()
	defer gServiceDataMutex.RUnlock()
	serviceIDs := gServiceDataNameIndex[serviceName]
	result := []ServiceData{}
	for sid := range serviceIDs {
		sd := gServiceDataMap[sid]
		result = append(result, *sd)
	}
	return result
}

func setupCollector() {
	collectAllServices()
	go func() {
		for {
			collectAllServices()
			time.Sleep(time.Second * 30)
		}
	}()
}
