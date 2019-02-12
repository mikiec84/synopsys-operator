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

package ctl

import (
	"fmt"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	horizoncomponents "github.com/blackducksoftware/horizon/pkg/components"
)

func GetOperatorReplicationController(namespace, synopsysOperatorImage, blackduckRegistrationKey string) (*horizoncomponents.ReplicationController, error) {
	// Add the Replication Controller to the Deployer
	var synopsysOperatorRCReplicas int32 = 1
	synopsysOperatorRC := horizoncomponents.NewReplicationController(horizonapi.ReplicationControllerConfig{
		APIVersion: "v1",
		//ClusterName:  "string",
		Name:      "synopsys-operator",
		Namespace: namespace,
		Replicas:  &synopsysOperatorRCReplicas,
		//ReadySeconds: "int32",
	})

	synopsysOperatorRC.AddLabelSelectors(map[string]string{"name": "synopsys-operator"})

	synopsysOperatorPod := horizoncomponents.NewPod(horizonapi.PodConfig{
		APIVersion: "v1",
		//ClusterName:            "string",
		Name:           "synopsys-operator",
		Namespace:      namespace,
		ServiceAccount: "synopsys-operator",
		//RestartPolicy:          "RestartPolicyType",
		//TerminationGracePeriod: "*int64",
		//ActiveDeadline:         "*int64",
		//Node:                   "string",
		//FSGID:                  "*int64",
		//Hostname:               "string",
		//SchedulerName:          "string",
		//DNSPolicy:              "DNSPolicType",
		//PriorityValue:          "*int32",
		//PriorityClass:          "string",
		//SELinux:                "*SELinuxType",
		//RunAsUser:              "*int64",
		//RunAsGroup:             "*int64",
		//ForceNonRoot:           "*bool",
	})

	synopsysOperatorPodLabels := map[string]string{"name": "synopsys-operator"}

	synopsysOperatorContainer := horizoncomponents.NewContainer(horizonapi.ContainerConfig{
		Name:       "synopsys-operator",
		Args:       []string{"/etc/synopsys-operator/config.json"},
		Command:    []string{"./operator"},
		Image:      synopsysOperatorImage,
		PullPolicy: horizonapi.PullAlways,
		//MinCPU:                   "string",
		//MaxCPU:                   "string",
		//MinMem:                   "string",
		//MaxMem:                   "string",
		//Privileged:               "*bool",
		//AllowPrivilegeEscalation: "*bool",
		//ReadOnlyFS:               "*bool",
		//ForceNonRoot:             "*bool",
		//SELinux:                  "*SELinuxType",
		//UID:                      "*int64",
		//AllocateStdin:            "bool",
		//StdinOnce:                "bool",
		//AllocateTTY:              "bool",
		//WorkingDirectory:         "string",
		//TerminationMsgPath:       "string",
		//TerminationMsgPolicy:     "TerminationMessagePolicyType",
	})
	synopsysOperatorContainer.AddPort(horizonapi.PortConfig{
		//Name:          "string",
		//Protocol:      "ProtocolType",
		//IP:            "string",
		//HostPort:      "string",
		ContainerPort: "8080",
	})
	synopsysOperatorContainer.AddEnv(horizonapi.EnvConfig{
		NameOrPrefix: "REGISTRATION_KEY",
		Type:         horizonapi.EnvVal,
		KeyOrVal:     blackduckRegistrationKey,
		//FromName:     "string",
	})
	synopsysOperatorContainer.AddVolumeMount(horizonapi.VolumeMountConfig{
		MountPath: "/etc/synopsys-operator",
		//Propagation: "*MountPropagationType",
		Name: "synopsys-operator",
		//SubPath:     "string",
		//ReadOnly:    "*bool",
	})

	synopsysOperatorContainerUI := horizoncomponents.NewContainer(horizonapi.ContainerConfig{
		Name: "synopsys-operator-ui",
		//Args:                     "[]string",
		Command:    []string{"./app"},
		Image:      synopsysOperatorImage,
		PullPolicy: horizonapi.PullAlways,
		//MinCPU:                   "string",
		//MaxCPU:                   "string",
		//MinMem:                   "string",
		//MaxMem:                   "string",
		//Privileged:               "*bool",
		//AllowPrivilegeEscalation: "*bool",
		//ReadOnlyFS:               "*bool",
		//ForceNonRoot:             "*bool",
		//SELinux:                  "*SELinuxType",
		//UID:                      "*int64",
		//AllocateStdin:            "bool",
		//StdinOnce:                "bool",
		//AllocateTTY:              "bool",
		//WorkingDirectory:         "string",
		//TerminationMsgPath:       "string",
		//TerminationMsgPolicy:     "TerminationMessagePolicyType",
	})
	synopsysOperatorContainerUI.AddPort(horizonapi.PortConfig{
		//Name:          "string",
		//Protocol:      "ProtocolType",
		//IP:            "string",
		//HostPort:      "string",
		ContainerPort: "3000",
	})
	synopsysOperatorContainerUI.AddEnv(horizonapi.EnvConfig{
		NameOrPrefix: "ADDR",
		Type:         horizonapi.EnvVal,
		KeyOrVal:     "0.0.0.0",
		//FromName:     "string",
	})
	synopsysOperatorContainerUI.AddEnv(horizonapi.EnvConfig{
		NameOrPrefix: "PORT",
		Type:         horizonapi.EnvVal,
		KeyOrVal:     "3000",
		//FromName:     "string",
	})
	synopsysOperatorContainerUI.AddEnv(horizonapi.EnvConfig{
		NameOrPrefix: "GO_ENV",
		Type:         horizonapi.EnvVal,
		KeyOrVal:     "development",
		//FromName:     "string",
	})

	// Create config map volume
	var synopsysOperatorVolumeDefaultMode int32 = 420
	synopsysOperatorVolume := horizoncomponents.NewConfigMapVolume(horizonapi.ConfigMapOrSecretVolumeConfig{
		VolumeName:      "synopsys-operator",
		MapOrSecretName: "synopsys-operator",
		//Items:           "map[string]KeyAndMode",
		DefaultMode: &synopsysOperatorVolumeDefaultMode,
		//Required:        "*bool",
	})

	synopsysOperatorPod.AddLabels(synopsysOperatorPodLabels)
	synopsysOperatorPod.AddContainer(synopsysOperatorContainer)
	synopsysOperatorPod.AddContainer(synopsysOperatorContainerUI)
	synopsysOperatorPod.AddVolume(synopsysOperatorVolume)
	synopsysOperatorRC.AddPod(synopsysOperatorPod)

	return synopsysOperatorRC, nil
}

func GetOperatorService(namespace string) (*horizoncomponents.Service, error) {

	// Add the Service to the Deployer
	synopsysOperatorService := horizoncomponents.NewService(horizonapi.ServiceConfig{
		APIVersion: "v1",
		//ClusterName:              "string",
		Name:      "synopsys-operator",
		Namespace: namespace,
		//ExternalName:             "string",
		//IPServiceType:            "ClusterIPServiceType",
		//ClusterIP:                "string",
		//PublishNotReadyAddresses: "bool",
		//TrafficPolicy:            "TrafficPolicyType",
		//Affinity:                 "string",
	})

	synopsysOperatorService.AddSelectors(map[string]string{"name": "synopsys-operator"})
	synopsysOperatorService.AddPort(horizonapi.ServicePortConfig{
		Name:       "synopsys-operator-ui",
		Port:       3000,
		TargetPort: "3000",
		//NodePort:   "int32",
		Protocol: horizonapi.ProtocolTCP,
	})
	synopsysOperatorService.AddPort(horizonapi.ServicePortConfig{
		Name:       "synopsys-operator-ui-standard-port",
		Port:       80,
		TargetPort: "3000",
		//NodePort:   "int32",
		Protocol: horizonapi.ProtocolTCP,
	})
	synopsysOperatorService.AddPort(horizonapi.ServicePortConfig{
		Name:       "synopsys-operator",
		Port:       8080,
		TargetPort: "8080",
		//NodePort:   "int32",
		Protocol: horizonapi.ProtocolTCP,
	})

	return synopsysOperatorService, nil
}

func GetOperatorConfigMap(namespace string) (*horizoncomponents.ConfigMap, error) {
	// Config Map
	synopsysOperatorConfigMap := horizoncomponents.NewConfigMap(horizonapi.ConfigMapConfig{
		APIVersion: "v1",
		//ClusterName: "string",
		Name:      "synopsys-operator",
		Namespace: namespace,
	})

	synopsysOperatorConfigMap.AddData(map[string]string{"config.json": fmt.Sprintf("{\"OperatorTimeBombInSeconds\":\"315576000\", \"DryRun\": false, \"LogLevel\": \"debug\", \"Namespace\": \"%s\", \"Threadiness\": 5, \"PostgresRestartInMins\": 10, \"NFSPath\" : \"/kubenfs\"}", namespace)})

	return synopsysOperatorConfigMap, nil
}

func GetOperatorServiceAccount(namespace string) (*horizoncomponents.ServiceAccount, error) {
	// Service Account
	synopsysOperatorServiceAccount := horizoncomponents.NewServiceAccount(horizonapi.ServiceAccountConfig{
		APIVersion: "v1",
		//ClusterName:    "string",
		Name:      "synopsys-operator",
		Namespace: namespace,
		//AutomountToken: "*bool",
	})

	return synopsysOperatorServiceAccount, nil
}

func GetOperatorClusterRoleBinding(namespace string) (*horizoncomponents.ClusterRoleBinding, error) {
	// Cluster Role Binding
	synopsysOperatorClusterRoleBinding := horizoncomponents.NewClusterRoleBinding(horizonapi.ClusterRoleBindingConfig{
		APIVersion: "rbac.authorization.k8s.io/v1beta1",
		//ClusterName: "string",
		Name:      "synopsys-operator-admin",
		Namespace: namespace,
	})
	synopsysOperatorClusterRoleBinding.AddSubject(horizonapi.SubjectConfig{
		Kind: "ServiceAccount",
		//APIGroup:  "string",
		Name:      "synopsys-operator",
		Namespace: namespace,
	})
	synopsysOperatorClusterRoleBinding.AddRoleRef(horizonapi.RoleRefConfig{
		APIGroup: "",
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
	})

	return synopsysOperatorClusterRoleBinding, nil
}