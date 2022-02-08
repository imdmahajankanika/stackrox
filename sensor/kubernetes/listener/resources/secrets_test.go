package resources

import (
	"testing"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/registries/types"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/sensor/common/registry"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	openshift311DockerConfigSecret = &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-dockercfg-6167c",
			Namespace: "test-ns",
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": "default",
			},
		},
		Data: map[string][]byte{
			".dockercfg": []byte(`
{
  "docker-registry.default.svc.cluster.local:5000": {
    "username": "serviceaccount",
    "password": "password",
    "email": "serviceaccount@example.org"
  }
}`),
		},
		Type: "kubernetes.io/dockercfg",
	}

	openshift4xDockerConfigSecret = &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-dockercfg-9w5gn",
			Namespace: "test-ns",
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": "default",
			},
		},
		Data: map[string][]byte{
			".dockercfg": []byte(`
{
  "image-registry.openshift-image-registry.svc:5000": {
    "username": "serviceaccount",
    "password": "password",
    "email": "serviceaccount@example.org"
  }
}`),
		},
		Type: "kubernetes.io/dockercfg",
	}
)

// checkTLS is a dummy implementation of registry.CheckTLS
func checkTLS(_ string) (bool, error) {
	return false, nil
}

func TestOpenShiftRegistrySecret_311(t *testing.T) {
	testutils.RunWithFeatureFlagEnabled(t, features.LocalImageScanning, testOpenShiftRegistrySecret311)
}

func testOpenShiftRegistrySecret311(t *testing.T) {
	regStore := registry.NewRegistryStore(checkTLS)
	d := newSecretDispatcher(regStore)

	_ = d.ProcessEvent(openshift311DockerConfigSecret, nil, central.ResourceAction_CREATE_RESOURCE)

	assert.Nil(t, regStore.GetAllInNamespace("random-ns"))

	regs := regStore.GetAllInNamespace(openshift311DockerConfigSecret.GetNamespace())
	assert.NotNil(t, regs)
	assert.Len(t, regs.GetAll(), 1)

	expectedRegConfig := &types.Config{
		Username:         "serviceaccount",
		Password:         "password",
		Insecure:         true,
		URL:              "https://docker-registry.default.svc.cluster.local:5000",
		RegistryHostname: "docker-registry.default.svc.cluster.local:5000",
		Autogenerated:    false,
	}

	assert.Equal(t, expectedRegConfig, regs.GetAll()[0].Config())
}

func TestOpenShiftRegistrySecret_4x(t *testing.T) {
	testutils.RunWithFeatureFlagEnabled(t, features.LocalImageScanning, testOpenShiftRegistrySecret4x)
}

func testOpenShiftRegistrySecret4x(t *testing.T) {
	regStore := registry.NewRegistryStore(checkTLS)
	d := newSecretDispatcher(regStore)

	_ = d.ProcessEvent(openshift4xDockerConfigSecret, nil, central.ResourceAction_CREATE_RESOURCE)

	assert.Nil(t, regStore.GetAllInNamespace("random-ns"))

	regs := regStore.GetAllInNamespace(openshift4xDockerConfigSecret.GetNamespace())
	assert.NotNil(t, regs)
	assert.Len(t, regs.GetAll(), 1)

	expectedRegConfig := &types.Config{
		Username:         "serviceaccount",
		Password:         "password",
		Insecure:         true,
		URL:              "https://image-registry.openshift-image-registry.svc:5000",
		RegistryHostname: "image-registry.openshift-image-registry.svc:5000",
		Autogenerated:    false,
	}

	assert.Equal(t, expectedRegConfig, regs.GetAll()[0].Config())
}
