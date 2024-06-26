/*
Copyright 2020 The Kubernetes Authors.

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

package azure

import (
	"context"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	network "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

type mockClient struct {
	vmss   []*compute.VirtualMachineScaleSet
	ifaces map[string][]*network.Interface
}

var _ client = &mockClient{}

func (c *mockClient) ListVMScaleSets(ctx context.Context) ([]*compute.VirtualMachineScaleSet, error) {
	return c.vmss, nil
}

func (c *mockClient) ListVMSSNetworkInterfaces(ctx context.Context, vmScaleSetName string) ([]*network.Interface, error) {
	return c.ifaces[vmScaleSetName], nil
}

func newTestInterfaces(ip string) []*network.Interface {
	return []*network.Interface{
		{
			Properties: &network.InterfacePropertiesFormat{
				IPConfigurations: []*network.InterfaceIPConfiguration{
					{
						Properties: &network.InterfaceIPConfigurationPropertiesFormat{
							PrivateIPAddress: to.Ptr(ip),
						},
					},
				},
			},
		},
	}
}

func TestGetSeeds(t *testing.T) {
	const (
		clusterTag  = "KubernetesCluster"
		clusterName = "test-cluster"
	)

	vmssNames := []string{"vmss0", "vmss1", "vmss"}
	ips := []string{"ip0", "ip1", "ip2"}
	client := &mockClient{
		vmss: []*compute.VirtualMachineScaleSet{
			{
				Name: to.Ptr(vmssNames[0]),
				Tags: map[string]*string{
					clusterTag: to.Ptr(clusterName),
				},
			},
			{
				Name: to.Ptr(vmssNames[1]),
				Tags: map[string]*string{
					clusterTag:             to.Ptr(clusterName),
					"not-relevant-tag-key": to.Ptr("val"),
				},
			},
			{
				// Irrelevalent VM that has no matching tag.
				Name: to.Ptr(vmssNames[2]),
				Tags: map[string]*string{
					"not-relevant-tag-key": to.Ptr("val"),
				},
			},
		},
		ifaces: map[string][]*network.Interface{
			vmssNames[0]: newTestInterfaces(ips[0]),
			vmssNames[1]: newTestInterfaces(ips[1]),
			vmssNames[2]: newTestInterfaces(ips[2]),
		},
	}
	provider := SeedProvider{
		client: client,
		tags: map[string]string{
			clusterTag: clusterName,
		},
	}
	actual, err := provider.GetSeeds()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := []string{ips[0], ips[1]}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected seeds %+v, but got %+v", expected, actual)
	}
}
