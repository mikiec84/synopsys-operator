/*
Copyright (C) 2019 Synopsys, Inc.

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

package blackduck

import (
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/api"
	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
)

// Creater interface
type Creater interface {
	Ensure(blackDuck *blackduckapi.Blackduck) error
	Versions() []string
	GetComponents(blackDuck *blackduckapi.Blackduck) (*api.ComponentList, error)
	GetPostgresComponents(blackDuck *blackduckapi.Blackduck) (*api.ComponentList, error)
	GetPVC(blackduck *blackduckapi.Blackduck) []*components.PersistentVolumeClaim
}
