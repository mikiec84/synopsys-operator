// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Blackduck) DeepCopyInto(out *Blackduck) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.View.DeepCopyInto(&out.View)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Blackduck.
func (in *Blackduck) DeepCopy() *Blackduck {
	if in == nil {
		return nil
	}
	out := new(Blackduck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Blackduck) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlackduckList) DeepCopyInto(out *BlackduckList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Blackduck, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlackduckList.
func (in *BlackduckList) DeepCopy() *BlackduckList {
	if in == nil {
		return nil
	}
	out := new(BlackduckList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BlackduckList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlackduckSpec) DeepCopyInto(out *BlackduckSpec) {
	*out = *in
	if in.ExternalPostgres != nil {
		in, out := &in.ExternalPostgres, &out.ExternalPostgres
		*out = new(PostgresExternalDBConfig)
		**out = **in
	}
	if in.PVC != nil {
		in, out := &in.PVC, &out.PVC
		*out = make([]PVC, len(*in))
		copy(*out, *in)
	}
	if in.NodeAffinities != nil {
		in, out := &in.NodeAffinities, &out.NodeAffinities
		*out = make(map[string][]NodeAffinity, len(*in))
		for key, val := range *in {
			var outVal []NodeAffinity
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]NodeAffinity, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
	if in.Environs != nil {
		in, out := &in.Environs, &out.Environs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ImageRegistries != nil {
		in, out := &in.ImageRegistries, &out.ImageRegistries
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.RegistryConfiguration.DeepCopyInto(&out.RegistryConfiguration)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlackduckSpec.
func (in *BlackduckSpec) DeepCopy() *BlackduckSpec {
	if in == nil {
		return nil
	}
	out := new(BlackduckSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlackduckStatus) DeepCopyInto(out *BlackduckStatus) {
	*out = *in
	if in.PVCVolumeName != nil {
		in, out := &in.PVCVolumeName, &out.PVCVolumeName
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlackduckStatus.
func (in *BlackduckStatus) DeepCopy() *BlackduckStatus {
	if in == nil {
		return nil
	}
	out := new(BlackduckStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlackduckView) DeepCopyInto(out *BlackduckView) {
	*out = *in
	if in.Clones != nil {
		in, out := &in.Clones, &out.Clones
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.StorageClasses != nil {
		in, out := &in.StorageClasses, &out.StorageClasses
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.CertificateNames != nil {
		in, out := &in.CertificateNames, &out.CertificateNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Environs != nil {
		in, out := &in.Environs, &out.Environs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ContainerTags != nil {
		in, out := &in.ContainerTags, &out.ContainerTags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SupportedVersions != nil {
		in, out := &in.SupportedVersions, &out.SupportedVersions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlackduckView.
func (in *BlackduckView) DeepCopy() *BlackduckView {
	if in == nil {
		return nil
	}
	out := new(BlackduckView)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Environs) DeepCopyInto(out *Environs) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Environs.
func (in *Environs) DeepCopy() *Environs {
	if in == nil {
		return nil
	}
	out := new(Environs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAffinity) DeepCopyInto(out *NodeAffinity) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAffinity.
func (in *NodeAffinity) DeepCopy() *NodeAffinity {
	if in == nil {
		return nil
	}
	out := new(NodeAffinity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PVC) DeepCopyInto(out *PVC) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PVC.
func (in *PVC) DeepCopy() *PVC {
	if in == nil {
		return nil
	}
	out := new(PVC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PostgresExternalDBConfig) DeepCopyInto(out *PostgresExternalDBConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PostgresExternalDBConfig.
func (in *PostgresExternalDBConfig) DeepCopy() *PostgresExternalDBConfig {
	if in == nil {
		return nil
	}
	out := new(PostgresExternalDBConfig)
	in.DeepCopyInto(out)
	return out
}
