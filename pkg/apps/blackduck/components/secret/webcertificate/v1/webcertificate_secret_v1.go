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

package v1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/store"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/types"
	apputils "github.com/blackducksoftware/synopsys-operator/pkg/apps/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"k8s.io/client-go/kubernetes"
)

// BdSecret holds the Black Duck secret configuration
type BdSecret struct {
	config     *protoform.Config
	kubeClient *kubernetes.Clientset
	blackDuck  *blackduckapi.Blackduck
}

func init() {
	store.Register(types.BlackDuckWebCertificateSecretV1, NewBdSecret)
}

// NewBdSecret returns the Black Duck secret configuration
func NewBdSecret(config *protoform.Config, kubeClient *kubernetes.Clientset, cr interface{}) (types.SecretInterface, error) {
	blackDuck, ok := cr.(*blackduckapi.Blackduck)
	if !ok {
		return nil, fmt.Errorf("unable to cast the interface to Black Duck object")
	}
	return &BdSecret{config: config, kubeClient: kubeClient, blackDuck: blackDuck}, nil
}

// GetSecret returns the secret
func (b *BdSecret) GetSecret() (*components.Secret, error) {
	certificateSecret := components.NewSecret(horizonapi.SecretConfig{Namespace: b.blackDuck.Spec.Namespace, Name: apputils.GetResourceName(b.blackDuck.Name, util.BlackDuckName, "webserver-certificate"), Type: horizonapi.SecretTypeOpaque})

	cert, key, err := b.getTLSCertKeyOrCreate()
	if err != nil {
		return nil, fmt.Errorf("unable to create the self signed certificate due to %+v", err)
	}
	certificateSecret.AddData(map[string][]byte{"WEBSERVER_CUSTOM_CERT_FILE": []byte(cert), "WEBSERVER_CUSTOM_KEY_FILE": []byte(key)})
	certificateSecret.AddLabels(apputils.GetVersionLabel("secret", b.blackDuck.Name, b.blackDuck.Spec.Version))

	return certificateSecret, nil
}

func (b *BdSecret) getTLSCertKeyOrCreate() (string, string, error) {
	if len(b.blackDuck.Spec.Certificate) > 0 && len(b.blackDuck.Spec.CertificateKey) > 0 {
		return b.blackDuck.Spec.Certificate, b.blackDuck.Spec.CertificateKey, nil
	}

	// Cert copy
	if len(b.blackDuck.Spec.CertificateName) > 0 && !strings.EqualFold(b.blackDuck.Spec.CertificateName, "default") {
		secret, err := util.GetSecret(b.kubeClient, b.blackDuck.Spec.CertificateName, apputils.GetResourceName(b.blackDuck.Name, util.BlackDuckName, "webserver-certificate"))
		if err == nil {
			cert, certok := secret.Data["WEBSERVER_CUSTOM_CERT_FILE"]
			key, keyok := secret.Data["WEBSERVER_CUSTOM_KEY_FILE"]
			if certok && keyok {
				return string(cert), string(key), nil
			}
		}
	}

	// default cert
	secret, err := util.GetSecret(b.kubeClient, b.config.Namespace, "blackduck-certificate")
	if err == nil {
		data := secret.Data
		if len(data) >= 2 {
			cert, certok := secret.Data["WEBSERVER_CUSTOM_CERT_FILE"]
			key, keyok := secret.Data["WEBSERVER_CUSTOM_KEY_FILE"]
			if !certok || !keyok {
				util.DeleteSecret(b.kubeClient, b.blackDuck.Spec.Namespace, apputils.GetResourceName(b.blackDuck.Name, util.BlackDuckName, "webserver-certificate"))
			} else {
				return string(cert), string(key), nil
			}
		}
	}

	// Default
	return CreateSelfSignedCert()
}

// CreateSelfSignedCert will create a random self signed certificate
func CreateSelfSignedCert() (string, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	template := x509.Certificate{
		SerialNumber: max,
		Subject: pkix.Name{
			CommonName:         "Black Duck",
			OrganizationalUnit: []string{"Cloud Native"},
			Organization:       []string{"Black Duck By Synopsys"},
			Locality:           []string{"Burlington"},
			StreetAddress:      []string{"800 District Avenue"},
			Province:           []string{"Massachusetts"},
			Country:            []string{"US"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365 * 3),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return "", "", err
	}
	certificate := &bytes.Buffer{}
	key := &bytes.Buffer{}
	pem.Encode(certificate, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	pemBlock, err := pemBlockForKey(priv)
	if err != nil {
		return "", "", err
	}

	pem.Encode(key, pemBlock)
	return certificate.String(), key.String(), nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal ECDSA private key: %v", err)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, nil
	}
}
