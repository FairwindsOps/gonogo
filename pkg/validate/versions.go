package validate

import (
	"k8s.io/client-go/discovery"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func (m* match) checkClusterVersion() error {
	cfg := config.GetConfigOrDie()

	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return err
	}

	version, err := client.ServerVersion()
	if err != nil {
		return  err
	}
	fmt.Println(version)

	
	return nil
}