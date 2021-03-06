package registry

import (
	"log"
	"time"

	"github.com/reportportal/goRP/commons"
)

const (
	retryTimeout  time.Duration = time.Second * 5
	retryAttempts int           = 3
)

//ServiceDiscovery provides methods to interact with registry (service discovery) service
type ServiceDiscovery interface {
	Register() error
	Deregister() error
	DoWithClient(func(client interface{}) (interface{}, error)) (interface{}, error)
}

//Register registers instance giving several tries
func Register(discovery ServiceDiscovery) {
	err := tryRegister(discovery)
	if nil != err {
		log.Fatal(err)
	}

	commons.ShutdownHook(func() error {
		log.Println("try to deregister")
		return Deregister(discovery)

	})
}

//Deregister de-registers instance giving several tries
func Deregister(discovery ServiceDiscovery) error {
	return tryDeregister(discovery)
}

func tryRegister(discovery ServiceDiscovery) error {
	return commons.Retry(retryAttempts, retryTimeout, func() error {
		e := discovery.Register()
		if nil != e {
			log.Printf("Cannot register service: %s", e)
		} else {
			log.Print("Service Registered!")
		}
		return e
	})
}

func tryDeregister(discovery ServiceDiscovery) error {
	return commons.Retry(retryAttempts, retryTimeout, func() error {
		return discovery.Deregister()
	})
}
