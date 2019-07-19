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

package actions

import (
	"encoding/json"
	"fmt"
	v1Configmap "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/configmap/global/v1"
	"sort"
	"strings"

	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps"
	blackduckclientset "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/client/clientset/versioned"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Blackduck)
// DB Table: Plural (Blackducks)
// Resource: Plural (Blackducks)
// Path: Plural (/blackducks)
// View Template Folder: Plural (/templates/blackducks/)

// BlackducksResource is the resource for the Blackduck model
type BlackducksResource struct {
	buffalo.Resource
	config          *protoform.Config
	kubeConfig      *rest.Config
	kubeClient      *kubernetes.Clientset
	blackduckClient *blackduckclientset.Clientset
}

// NewBlackduckResource will instantiate the Black Duck Resource
func NewBlackduckResource(config *protoform.Config, kubeConfig *rest.Config) (*BlackducksResource, error) {
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create kube client due to %+v", err)
	}
	hubClient, err := blackduckclientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create Black Duck client due to %+v", err)
	}
	return &BlackducksResource{config: config, kubeConfig: kubeConfig, kubeClient: kubeClient, blackduckClient: hubClient}, nil
}

// List gets all Hubs. This function is mapped to the path
// GET /blackducks
func (v BlackducksResource) List(c buffalo.Context) error {
	SetVersion(c)
	blackducks, err := util.ListHubs(v.blackduckClient, v.config.Namespace)
	if err != nil {
		return c.Error(500, err)
	}
	// Make blackducks available inside the html template
	c.Set("blackducks", blackducks.Items)
	return c.Render(200, r.HTML("blackducks/index.html", "old_application.html"))
}

// parseParams will split the input namespace/name string and return the namespace and name field separately
func parseParams(name string) (string, string) {
	values := strings.SplitN(name, ":", 2)
	return values[0], values[1]
}

// Show gets the data for one Blackduck. This function is mapped to
// the path GET /blackducks/{blackduck_id}
func (v BlackducksResource) Show(c buffalo.Context) error {
	SetVersion(c)
	namespace, name := parseParams(c.Param("blackduck_id"))
	blackduck, err := util.GetHub(v.blackduckClient, namespace, name)
	if err != nil {
		return c.Error(500, err)
	}
	// Make blackduck available inside the html template
	c.Set("blackduck", blackduck)
	return c.Render(200, r.HTML("blackducks/show.html", "old_application.html"))
}

// New renders the form for creating a new Blackduck.
// This function is mapped to the path GET /blackducks/new
func (v BlackducksResource) New(c buffalo.Context) error {
	SetVersion(c)
	blackduckSpec := util.GetBlackDuckDefaultPersistentStorageLatest()
	blackduck := &blackduckapi.Blackduck{}
	blackduck.Spec = *blackduckSpec
	if v.config.IsClusterScoped {
		blackduck.Spec.Namespace = ""
	} else {
		blackduck.Spec.Namespace = v.config.Namespace
	}
	blackduck.Spec.PersistentStorage = true
	blackduck.Spec.PVCStorageClass = ""
	blackduck.Spec.ScanType = "Artifacts"

	err := v.common(c, blackduck, false)
	if err != nil {
		return err
	}
	// Make blackduck available inside the html template
	c.Set("blackduck", blackduck)

	return c.Render(200, r.HTML("blackducks/new.html", "old_application.html"))
}

func (v BlackducksResource) common(c buffalo.Context, bd *blackduckapi.Blackduck, isEdit bool) error {
	// External postgres
	if bd.Spec.ExternalPostgres == nil {
		bd.Spec.ExternalPostgres = &blackduckapi.PostgresExternalDBConfig{}
	}

	// PVCs
	if bd.Spec.PVC == nil {
		bd.Spec.PVC = []blackduckapi.PVC{
			{
				Name: "blackduck-postgres",
				Size: "150Gi",
			},
		}
	}

	// PVC storage classes
	if bd.View.StorageClasses == nil {
		var storageList map[string]string
		storageList = make(map[string]string)
		storageClasses, err := util.ListStorageClasses(v.kubeClient)
		if err != nil {
			c.Error(404, fmt.Errorf("\"message\": \"Failed to List the storage class due to %+v\"", err))
		}
		for _, storageClass := range storageClasses.Items {
			storageList[fmt.Sprintf("%s (%s)", storageClass.GetName(), storageClass.Provisioner)] = storageClass.GetName()
		}
		storageList[fmt.Sprintf("%s (%s)", "None", "Disable dynamic provisioner")] = ""
		if isEdit && bd.Spec.PersistentStorage {
			for key, value := range storageList {
				if strings.EqualFold(value, bd.Spec.PVCStorageClass) {
					bd.View.StorageClasses = map[string]string{key: value}
					break
				}
			}
		} else {
			bd.View.StorageClasses = storageList
		}
	}

	blackducks, _ := util.ListHubs(v.blackduckClient, v.config.Namespace)
	// Clone Black Ducks
	if bd.View.Clones == nil {
		keys := make(map[string]string)
		for _, v := range blackducks.Items {
			if strings.EqualFold(v.Status.State, "running") {
				keys[v.Name] = v.Name
			}
		}
		keys["None"] = ""
		if isEdit {
			for key, value := range keys {
				if strings.EqualFold(value, bd.Spec.DbPrototype) {
					bd.View.Clones = map[string]string{key: value}
					break
				}
			}
		} else {
			bd.View.Clones = keys
		}
	}

	// certificates
	if bd.View.CertificateNames == nil {
		certificateNames := []string{"default", "manual"}
		for _, hub := range blackducks.Items {
			if strings.EqualFold(hub.Spec.CertificateName, "manual") {
				certificateNames = append(certificateNames, hub.Spec.Namespace)
			}
		}
		bd.View.CertificateNames = certificateNames
	}

	// environment variables
	if bd.View.Environs == nil {
		env := v1Configmap.GetHubKnobs()
		environs := []string{}
		for key, value := range env {
			if !strings.EqualFold(value, "") {
				environs = append(environs, fmt.Sprintf("%s:%s", key, value))
			}
		}

		if len(bd.Spec.Environs) > 0 {
			bd.View.Environs = bd.Spec.Environs
		} else {
			bd.View.Environs = environs
		}
	}

	// supported versions
	if bd.View.SupportedVersions == nil {
		bd.View.SupportedVersions = apps.NewApp(&protoform.Config{}, v.kubeConfig).Blackduck().Versions()
		sort.Sort(sort.Reverse(sort.StringSlice(bd.View.SupportedVersions)))
	}
	return nil
}

func (v BlackducksResource) redirect(c buffalo.Context, blackduck *blackduckapi.Blackduck, err error) error {
	if err != nil {
		c.Flash().Add("warning", err.Error())
		// Make blackduck available inside the html template
		err = v.common(c, blackduck, false)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debugf("edit hub in create: %+v", blackduck)

		c.Set("blackduck", blackduck)
		return c.Render(422, r.HTML("blackducks/new.html", "old_application.html"))
	}
	return nil
}

// Create adds a Blackduck to the DB. This function is mapped to the
// path POST /blackducks
func (v BlackducksResource) Create(c buffalo.Context) error {
	SetVersion(c)
	// Allocate an empty Blackduck
	blackduck := &blackduckapi.Blackduck{}

	err := v.postSubmit(c, blackduck)
	if err != nil {
		log.Error(err)
		return v.redirect(c, blackduck, err)
	}

	log.Infof("create Black Duck: %+v", blackduck)

	blackducks, err := util.ListHubs(v.blackduckClient, blackduck.Spec.Namespace)
	if err != nil {
		return fmt.Errorf("unable to list Black Duck instances in namespace '%s' due to %+v", blackduck.Spec.Namespace, err)
	}

	// When running in cluster scope mode, custom resources do not have a namespace so the above command returns everything and we need to check Spec.Namespace.
	for _, bd := range blackducks.Items {
		if strings.EqualFold(bd.Spec.Namespace, blackduck.Spec.Namespace) {
			return fmt.Errorf("a Black Duck instance already exists in namespace '%s', only one instance per namespace is allowed", blackduck.Spec.Namespace)
		}
	}

	_, err = util.GetHub(v.blackduckClient, blackduck.Spec.Namespace, blackduck.Name)
	if err == nil {
		return v.redirect(c, blackduck, fmt.Errorf("already '%s' Black Duck instance exist in '%s' namespace", blackduck.Name, blackduck.Spec.Namespace))
	}

	// Get the Synopsys Operator namespace by CRD scope
	scope := apiextensions.NamespaceScoped
	if v.config.IsClusterScoped {
		scope = apiextensions.ClusterScoped
	}

	operatorNamespace, err := util.GetOperatorNamespaceByCRDScope(v.kubeClient, util.BlackDuckCRDName, scope, blackduck.Spec.Namespace)
	if err != nil {
		return v.redirect(c, blackduck, err)
	}

	_, err = util.CreateHub(v.blackduckClient, operatorNamespace, &blackduckapi.Blackduck{ObjectMeta: metav1.ObjectMeta{Name: blackduck.Name}, Spec: blackduck.Spec})
	if err != nil {
		return v.redirect(c, blackduck, err)
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Black Duck was created successfully")

	c.Set("blackducks", blackducks.Items)
	// and redirect to the blackducks index page
	return c.Redirect(302, "/blackducks/%s", fmt.Sprintf("%s:%s", blackduck.Spec.Namespace, blackduck.Name))
}

// Edit renders a edit form for a Blackduck. This function is
// mapped to the path GET /blackducks/{blackduck_id}/edit
func (v BlackducksResource) Edit(c buffalo.Context) error {
	SetVersion(c)
	namespace, name := parseParams(c.Param("blackduck_id"))
	blackduck, err := util.GetHub(v.blackduckClient, namespace, name)
	if err != nil {
		return c.Error(404, err)
	}

	// decode the password back during edit action
	if blackduck.Spec.ExternalPostgres == nil {
		blackduck.Spec.PostgresPassword, _ = util.Base64Decode(blackduck.Spec.PostgresPassword)
		blackduck.Spec.AdminPassword, _ = util.Base64Decode(blackduck.Spec.AdminPassword)
		blackduck.Spec.UserPassword, _ = util.Base64Decode(blackduck.Spec.UserPassword)
	} else {
		blackduck.Spec.ExternalPostgres.PostgresAdminPassword, _ = util.Base64Decode(blackduck.Spec.ExternalPostgres.PostgresAdminPassword)
		blackduck.Spec.ExternalPostgres.PostgresUserPassword, _ = util.Base64Decode(blackduck.Spec.ExternalPostgres.PostgresUserPassword)
	}

	err = v.common(c, blackduck, true)
	if err != nil {
		return c.Error(500, err)
	}
	return c.Render(200, r.Auto(c, blackduck))
}

func (v BlackducksResource) postSubmit(c buffalo.Context, blackduck *blackduckapi.Blackduck) error {
	// Bind blackduck to the html form elements
	if err := c.Bind(blackduck); err != nil {
		log.Errorf("error binding Black Duck '%+v' due to %+v", c, err)
		return errors.WithStack(err)
	}

	if len(blackduck.View.NodeAffinities) > 0 {
		nodeAffinities := map[string][]blackduckapi.NodeAffinity{}
		err := json.Unmarshal([]byte(blackduck.View.NodeAffinities), &nodeAffinities)
		if err != nil {
			return err
		}
		blackduck.Spec.NodeAffinities = nodeAffinities
	}

	if !blackduck.Spec.PersistentStorage {
		blackduck.Spec.PVC = nil
	} else {
		// Remove postgres volume if we use an external db
		if *blackduck.Spec.ExternalPostgres != (blackduckapi.PostgresExternalDBConfig{}) {
			blackduck.Spec.PVC = nil
		}
	}

	// Change back to nil if the configuration is empty
	if *blackduck.Spec.ExternalPostgres == (blackduckapi.PostgresExternalDBConfig{}) {
		blackduck.Spec.ExternalPostgres = nil
		blackduck.Spec.PostgresPassword = util.Base64Encode([]byte(blackduck.Spec.PostgresPassword))
		blackduck.Spec.AdminPassword = util.Base64Encode([]byte(blackduck.Spec.AdminPassword))
		blackduck.Spec.UserPassword = util.Base64Encode([]byte(blackduck.Spec.UserPassword))
	} else {
		blackduck.Spec.ExternalPostgres.PostgresAdminPassword = util.Base64Encode([]byte(blackduck.Spec.ExternalPostgres.PostgresAdminPassword))
		blackduck.Spec.ExternalPostgres.PostgresUserPassword = util.Base64Encode([]byte(blackduck.Spec.ExternalPostgres.PostgresUserPassword))
	}
	return nil
}

// Update changes a Blackduck in the DB. This function is mapped to
// the path PUT /blackducks/{blackduck_id}
func (v BlackducksResource) Update(c buffalo.Context) error {
	SetVersion(c)
	// Allocate an empty Blackduck
	blackduck := &blackduckapi.Blackduck{}

	err := v.postSubmit(c, blackduck)
	if err != nil {
		log.Error(err)
		return v.redirect(c, blackduck, err)
	}

	latestBlackduck, err := util.GetHub(v.blackduckClient, blackduck.Spec.Namespace, blackduck.Name)
	if err != nil {
		log.Errorf("error getting Black Duck '%s' in namespace '%s' due to %+v", blackduck.Name, blackduck.Spec.Namespace, err)
		return v.redirect(c, blackduck, err)
	}

	latestBlackduck.Spec = blackduck.Spec
	_, err = util.UpdateBlackduck(v.blackduckClient, blackduck.Spec.Namespace, latestBlackduck)

	if err != nil {
		log.Errorf("error updating Black Duck '%s' in namespace '%s' due to %+v", blackduck.Name, blackduck.Spec.Namespace, err)
		return v.redirect(c, blackduck, err)
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Black Duck was updated successfully")

	blackducks, _ := util.ListHubs(v.blackduckClient, v.config.Namespace)
	c.Set("blackducks", blackducks.Items)
	// and redirect to the blackducks index page
	return c.Redirect(302, "/blackducks/%s", fmt.Sprintf("%s:%s", blackduck.Spec.Namespace, blackduck.Name))
}

// Destroy deletes a Blackduck from the DB. This function is mapped
// to the path DELETE /blackducks/{blackduck_id}
func (v BlackducksResource) Destroy(c buffalo.Context) error {
	SetVersion(c)
	log.Infof("delete Black Duck request %v", c.Param("blackduck"))
	namespace, name := parseParams(c.Param("blackduck_id"))
	_, err := util.GetHub(v.blackduckClient, namespace, name)
	// To find the Blackduck the parameter blackduck_id is used.
	if err != nil {
		return c.Error(404, err)
	}

	// This is on the event loop.
	err = v.blackduckClient.SynopsysV1().Blackducks(namespace).Delete(name, &metav1.DeleteOptions{})
	// To find the Blackduck the parameter blackduck_id is used.
	if err != nil {
		log.Errorf("error deleting Black Duck '%s' in namespace '%s' due to %+v", name, namespace, err)
		return c.Error(404, err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Black Duck was deleted successfully")

	// blackducks, _ := util.ListHubs(v.blackduckClient, "")
	// c.Set("hubs", blackducks.Items)

	// Redirect to the blackducks index page
	return c.Redirect(302, "/blackducks")
}

// ChangeState Used to change state of a Blackduck instance
// POST  /blackducks/{blackduck_id}/state
func (v BlackducksResource) ChangeState(c buffalo.Context) error {
	SetVersion(c)
	if c.Param("state") == "" {
		return c.Redirect(400, "/blackducks")
	}

	namespace, name := parseParams(c.Param("blackduck_id"))

	blackduck, err := util.GetHub(v.blackduckClient, namespace, name)
	if err != nil {
		log.Errorf("error getting Black Duck '%s' in namespace '%s' due to %+v", name, namespace, err)
		return v.redirect(c, blackduck, err)
	}

	blackduck.Spec.DesiredState = c.Param("state")

	_, err = v.blackduckClient.SynopsysV1().Blackducks(blackduck.Spec.Namespace).Update(blackduck)
	if err != nil {
		log.Errorf("error updating Black Duck '%s' in namespace '%s' due to %+v", blackduck.Name, blackduck.Spec.Namespace, err)
		return v.redirect(c, blackduck, err)
	}

	// Redirect to the blackducks index page
	return c.Redirect(302, "/blackducks")
}
