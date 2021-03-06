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

package util

import (
	"fmt"
	"strings"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
)

// Container defines the configuration for a container
type Container struct {
	ContainerConfig       *horizonapi.ContainerConfig
	EnvConfigs            []*horizonapi.EnvConfig
	VolumeMounts          []*horizonapi.VolumeMountConfig
	PortConfig            []*horizonapi.PortConfig
	ActionConfig          *horizonapi.ActionConfig
	ReadinessProbeConfigs []*horizonapi.ProbeConfig
	LivenessProbeConfigs  []*horizonapi.ProbeConfig
	PreStopConfig         *horizonapi.ActionConfig
}

// PodConfig used for configuring the pod
type PodConfig struct {
	Name                   string
	Labels                 map[string]string
	ServiceAccount         string
	Containers             []*Container
	Volumes                []*components.Volume
	InitContainers         []*Container
	PodAffinityConfigs     map[horizonapi.AffinityType][]*horizonapi.PodAffinityConfig
	PodAntiAffinityConfigs map[horizonapi.AffinityType][]*horizonapi.PodAffinityConfig
	NodeAffinityConfigs    map[horizonapi.AffinityType][]*horizonapi.NodeAffinityConfig
	ImagePullSecrets       []string
	FSGID                  *int64
}

// MergeEnvMaps will merge the source and destination environs. If the same value exist in both, source environ will given more preference
func MergeEnvMaps(source, destination map[string]string) map[string]string {
	// if the source key present in the destination map, it will overrides the destination value
	// if the source value is empty, then delete it from the destination
	for key, value := range source {
		if len(value) == 0 {
			delete(destination, key)
		} else {
			destination[key] = value
		}
	}
	return destination
}

// MergeEnvSlices will merge the source and destination environs. If the same value exist in both, source environ will given more preference
func MergeEnvSlices(source, destination []string) []string {
	// create a destination map
	destinationMap := make(map[string]string)
	for _, value := range destination {
		values := strings.SplitN(value, ":", 2)
		if len(values) == 2 {
			mapKey := strings.TrimSpace(values[0])
			mapValue := strings.TrimSpace(values[1])
			if len(mapKey) > 0 && len(mapValue) > 0 {
				destinationMap[mapKey] = mapValue
			}
		}
	}

	// if the source key present in the destination map, it will overrides the destination value
	// if the source value is empty, then delete it from the destination
	for _, value := range source {
		values := strings.SplitN(value, ":", 2)
		if len(values) == 2 {
			mapKey := strings.TrimSpace(values[0])
			mapValue := strings.TrimSpace(values[1])
			if len(mapValue) == 0 {
				delete(destinationMap, mapKey)
			} else {
				destinationMap[mapKey] = mapValue
			}
		}
	}

	// convert destination map to string array
	mergedValues := []string{}
	for key, value := range destinationMap {
		mergedValues = append(mergedValues, fmt.Sprintf("%s:%s", key, value))
	}
	return mergedValues
}

// UniqueStringSlice returns a unique subset of the string slice provided.
func UniqueStringSlice(input []string) []string {
	output := []string{}
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			output = append(output, val)
		}
	}

	return output
}

// RemoveFromStringSlice will remove the string from the slice and it will maintain the order
func RemoveFromStringSlice(slice []string, str string) []string {
	for index, value := range slice {
		if value == str {
			slice = append(slice[:index], slice[index+1:]...)
		}
	}
	return slice
}

// IsExposeServiceValid validates the expose service type
func IsExposeServiceValid(serviceType string) bool {
	switch strings.ToUpper(serviceType) {
	case NONE, NODEPORT, LOADBALANCER, OPENSHIFT:
		return true
	}
	return false
}
