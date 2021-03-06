/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"fmt"
	"strings"

	synopsysv1 "github.com/blackducksoftware/synopsys-operator/meta-builder/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func patchAlert(alert *synopsysv1.Alert, objects map[string]runtime.Object) map[string]runtime.Object {
	patcher := AlertPatcher{
		alert:   alert,
		objects: objects,
	}
	return patcher.patch()
}

type AlertPatcher struct {
	alert   *synopsysv1.Alert
	objects map[string]runtime.Object
}

func (p *AlertPatcher) patch() map[string]runtime.Object {
	patches := [](func() error){
		p.patchNamespace,
		p.patchEnvirons,
		p.patchSecrets,
		p.patchStandAlone,
		p.patchPersistentStorage,
		p.patchExposeUserInterface,
		p.patchAlertImage,
	}
	for _, f := range patches {
		err := f()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	return p.objects
}

func (p *AlertPatcher) patchNamespace() error {
	accessor := meta.NewAccessor()
	for _, runtimeObject := range p.objects {
		accessor.SetNamespace(runtimeObject, p.alert.Spec.Namespace)
	}
	return nil
}

func (p *AlertPatcher) patchEnvirons() error {
	ConfigMapUniqueID := "ConfigMap.default.demo-alert-blackduck-config"
	configMapRuntimeObject, ok := p.objects[ConfigMapUniqueID]
	if !ok {
		return nil
	}
	configMap := configMapRuntimeObject.(*corev1.ConfigMap)
	for _, e := range p.alert.Spec.Environs {
		vals := strings.Split(e, ":") // TODO - doesn't handle multiple colons
		if len(vals) != 2 {
			fmt.Printf("Could not split environ '%s' on ':'\n", e) // TODO change to log
			continue
		}
		environKey := strings.TrimSpace(vals[0])
		environVal := strings.TrimSpace(vals[1])
		configMap.Data[environKey] = environVal
	}
	return nil
}

func (p *AlertPatcher) patchSecrets() error {
	SecretUniqueID := "Secret.default.demo-alert-secret"
	secretRuntimeObject, ok := p.objects[SecretUniqueID]
	if !ok {
		return nil
	}
	secret := secretRuntimeObject.(*corev1.Secret)
	for _, s := range p.alert.Spec.Environs {
		vals := strings.Split(s, ":") // TODO - doesn't handle multiple colons
		if len(vals) != 2 {
			fmt.Printf("Could not split environ '%s' on ':'\n", s) // TODO change to log
			continue
		}
		secretKey := strings.TrimSpace(vals[0])
		secretVal := strings.TrimSpace(vals[1])
		secret.Data[secretKey] = []byte(secretVal)
	}
	return nil
}

func (p *AlertPatcher) patchDesiredState() error {
	accessor := meta.NewAccessor()
	if strings.EqualFold(p.alert.Spec.DesiredState, "STOP") {
		for uniqueID, runtimeObject := range p.objects {
			if k, _ := accessor.Kind(runtimeObject); k != "PersistentVolumeClaim" {
				delete(p.objects, uniqueID)
			}
		}
	}
	return nil
}

func (p *AlertPatcher) patchPort() error {
	port := *p.alert.Spec.Port
	ReplicationContollerUniqueID := "ReplicationController.default.demo-alert-alert"
	replicationControllerRuntimeObject, ok := p.objects[ReplicationContollerUniqueID]
	if !ok {
		return nil
	}
	replicationController := replicationControllerRuntimeObject.(*corev1.ReplicationController)
	replicationController.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = port
	replicationController.Spec.Template.Spec.Containers[0].Ports[0].Protocol = corev1.ProtocolTCP

	ServiceUniqueID := "Service.default.demo-alert-alert"
	serviceRuntimeObject, ok := p.objects[ServiceUniqueID]
	if !ok {
		return nil
	}
	service := serviceRuntimeObject.(*corev1.Service)
	service.Spec.Ports[0].Name = fmt.Sprintf("port-%d", port)
	service.Spec.Ports[0].Port = port
	service.Spec.Ports[0].TargetPort = intstr.IntOrString{IntVal: port}
	service.Spec.Ports[0].Protocol = corev1.ProtocolTCP

	ServiceExposedUniqueID := "Service.default.demo-alert-exposed"
	serviceExposedRuntimeObject, ok := p.objects[ServiceExposedUniqueID]
	if !ok {
		return nil
	}
	serviceExposed := serviceExposedRuntimeObject.(*corev1.Service)
	serviceExposed.Spec.Ports[0].Name = fmt.Sprintf("port-%d", port)
	service.Spec.Ports[0].Port = port
	service.Spec.Ports[0].TargetPort = intstr.IntOrString{IntVal: port}
	service.Spec.Ports[0].Protocol = corev1.ProtocolTCP

	// TODO: Support Openshift Routes
	// RouteUniqueID := "Route.default.demo-alert-route"
	// routeRuntimeObject := p.objects[RouteUniqueID]

	return nil
}

func (p *AlertPatcher) patchAlertImage() error {
	uniqueID := "ReplicationController.default.demo-alert-alert"
	alertReplicationControllerRuntimeObject, ok := p.objects[uniqueID]
	if !ok {
		return nil
	}
	alertReplicationController := alertReplicationControllerRuntimeObject.(*corev1.ReplicationController)
	alertReplicationController.Spec.Template.Spec.Containers[0].Image = p.alert.Spec.AlertImage
	return nil
}

func (p *AlertPatcher) patchAlertMemory() error {
	uniqueID := "ReplicationController.default.demo-alert-alert"
	alertReplicationControllerRuntimeObject, ok := p.objects[uniqueID]
	if !ok {
		return nil
	}
	alertReplicationController := alertReplicationControllerRuntimeObject.(*corev1.ReplicationController)
	minAndMaxMem, _ := resource.ParseQuantity(p.alert.Spec.AlertMemory)
	alertReplicationController.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceMemory] = minAndMaxMem
	alertReplicationController.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] = minAndMaxMem
	return nil
}

func (p *AlertPatcher) patchPersistentStorage() error {
	if (p.alert.Spec.PersistentStorage == synopsysv1.PersistentStorage{}) {
		PVCUniqueID := "PersistentVolumeClaim.default.demo-alert-pvc"
		delete(p.objects, PVCUniqueID)
		return nil
	}
	// Patch PVC Name
	PVCUniqueID := "PersistentVolumeClaim.default.demo-alert-pvc"
	PVCRuntimeObject, ok := p.objects[PVCUniqueID]
	if !ok {
		return nil
	}
	pvc := PVCRuntimeObject.(*corev1.PersistentVolumeClaim)

	name := fmt.Sprintf("%s-%s-%s", p.alert.Name, "alert", p.alert.Spec.PersistentStorage.PVCName)
	if p.alert.Annotations["synopsys.com/created.by"] == "pre-2019.6.0" {
		name = p.alert.Spec.PersistentStorage.PVCName
	}
	pvc.Name = name

	return nil
}

func (p *AlertPatcher) patchStandAlone() error {
	if (p.alert.Spec.StandAlone == synopsysv1.StandAlone{}) {
		// Remove Cfssl Resources
		uniqueID := "ReplicationController.default.demo-alert-cfssl"
		delete(p.objects, uniqueID)
		uniqueID = "Service.default.demo-alert-cfssl"
		delete(p.objects, uniqueID)

		// Add Environ to use BlackDuck Cfssl
		ConfigMapUniqueID := "ConfigMap.default.demo-alert-blackduck-config"
		configMapRuntimeObject, ok := p.objects[ConfigMapUniqueID]
		if !ok {
			return nil
		}
		configMap := configMapRuntimeObject.(*corev1.ConfigMap)
		configMap.Data["HUB_CFSSL_HOST"] = fmt.Sprintf("%s-%s-%s", p.alert.Name, "alert", "cfssl")
	} else {
		uniqueID := "ReplicationController.default.demo-alert-cfssl"
		alertCfsslReplicationControllerRuntimeObject, ok := p.objects[uniqueID]
		if !ok {
			return nil
		}
		// patch Cfssl Image
		alertCfsslReplicationController := alertCfsslReplicationControllerRuntimeObject.(*corev1.ReplicationController)
		alertCfsslReplicationController.Spec.Template.Spec.Containers[0].Image = p.alert.Spec.StandAlone.CfsslImage
		// patch Cfssl Memory
		minAndMaxMem, _ := resource.ParseQuantity(p.alert.Spec.StandAlone.CfsslMemory)
		alertCfsslReplicationController.Spec.Template.Spec.Containers[0].Resources.Requests[corev1.ResourceMemory] = minAndMaxMem
		alertCfsslReplicationController.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] = minAndMaxMem
	}
	return nil
}

func (p *AlertPatcher) patchExposeUserInterface() error {
	nodePortUniqueID := "Service.default.demo-alert-exposed"
	loadbalancerUniqueID := "Service.default.demo-alert-exposed"
	routeUniqueID := "Service.default.demo-alert-exposed"
	switch p.alert.Spec.ExposeService {
	case "NODEPORT":
		delete(p.objects, loadbalancerUniqueID)
		delete(p.objects, routeUniqueID)
	case "LOADBALANCER":
		delete(p.objects, nodePortUniqueID)
		delete(p.objects, routeUniqueID)
	case "OPENSHIFT":
		delete(p.objects, nodePortUniqueID)
		delete(p.objects, loadbalancerUniqueID)
	default:
		delete(p.objects, nodePortUniqueID)
		delete(p.objects, loadbalancerUniqueID)
		delete(p.objects, routeUniqueID)
	}
	return nil
}
