/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	piftester "github.com/blackducksoftware/synopsys-operator/contrib/hydra/pkg/apps/piftester"
	"github.com/blackducksoftware/synopsys-operator/contrib/hydra/pkg/kubebuilder"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	configPath := os.Args[1]
	auxConfigPath := os.Args[2]
	config := piftester.ReadConfig(configPath)
	if config == nil {
		panic("didn't find config")
	}
	auxConfig := piftester.ReadAuxiliaryConfig(auxConfigPath)
	if auxConfig == nil {
		panic("didn't find auxconfig")
	}
	config.AuxConfig = auxConfig
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %s\n", string(jsonBytes))
	runPiftester(config)
}

func runPiftester(config *piftester.Config) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags(config.MasterURL, config.KubeConfigPath)
	//		kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err)
	}

	var resources kubebuilder.Resources
	if config.AuxConfig.IsOpenshift {
		resources = piftester.NewOpenshift(config)
	} else {
		resources = piftester.NewKube(config)
	}

	namespace := config.AuxConfig.Namespace
	builder := kubebuilder.NewBuilder(namespace, resources, clientset)
	builder.CreateResources()
}
