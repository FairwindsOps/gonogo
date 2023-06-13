// Copyright 2021 FairwindsOps, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package helm

import (
	"context"
	"fmt"
	"os"
	"sync"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	// This is required to auth to cloud providers (i.e. GKE)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"k8s.io/apimachinery/pkg/api/meta"
)

// kube wraps a kubernetes client interface
type kube struct {
	Client kubernetes.Interface
}

// GetData fulfills the kubernetes client interface in the fairwinds opa package
func (h *kube) GetData(ctx context.Context, group, kind string) ([]interface{}, error) {
	return nil, nil
}

// getKubeInstance returns a Kubernetes interface
func getKubeInstance() *kube {
	var kubeClient *kube
	once.Do(func() {

		kubeClient = &kube{
			Client: getKubeClient(),
		}

	})
	return kubeClient
}

type dynamicClientInstance struct {
	Client     dynamic.Interface
	RESTMapper meta.RESTMapper
}

var once sync.Once
var clientOnceDynamic sync.Once

// getKubeClient returns a clientset instance
func getKubeClient() *kubernetes.Clientset {
	kubeConf, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting kubeconfig:", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		fmt.Println("Error creating kubernetes client:", err)
		os.Exit(1)
	}
	return clientset
}

// GetDynamicInstance returns a dynamic client instance
func getDynamicInstance() *dynamicClientInstance {
	var dynamicClient *dynamicClientInstance
	clientOnceDynamic.Do(func() {

		dynamicClient = &dynamicClientInstance{
			Client:     getKubeClientDynamic(),
			RESTMapper: getRESTMapper(),
		}
	})
	return dynamicClient
}

func getKubeClientDynamic() dynamic.Interface {
	kubeConf, err := config.GetConfig()
	if err != nil {
		klog.Fatalf("Error getting kubeconfig: %v", err)
	}
	clientset, err := dynamic.NewForConfig(kubeConf)
	if err != nil {
		klog.Fatalf("Error creating dynamic kubernetes client: %v", err)
	}
	return clientset
}

func getRESTMapper() meta.RESTMapper {
	kubeConf, err := config.GetConfig()
	if err != nil {
		klog.Fatalf("Error getting kubeconfig: %v", err)
	}

	httpClient, err := rest.HTTPClientFor(kubeConf)
	if err != nil {
		klog.Fatal("error creating httpClient using kubeconfig: %s", err.Error())
	}

	restmapper, err := apiutil.NewDynamicRESTMapper(kubeConf, httpClient)
	if err != nil {
		klog.Fatalf("Error creating REST Mapper: %v", err)
	}
	return restmapper
}
