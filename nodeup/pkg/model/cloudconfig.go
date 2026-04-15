/*
Copyright 2019 The Kubernetes Authors.

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

package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/openstack"
	"k8s.io/kops/upup/pkg/fi/nodeup/nodetasks"
)

const (
	CloudConfigFilePath       = "/etc/kubernetes/cloud.config"
	AzureCloudConfigFilePath  = "/etc/kubernetes/azure.json"
	InTreeCloudConfigFilePath = "/etc/kubernetes/in-tree-cloud.config"

	// VM UUID is set by cloud-init
	VM_UUID_FILE_PATH = "/etc/vmware/vm_uuid"
)

// azureCloudConfig is the configuration passed to Cloud Provider Azure.
// The specification is described in https://cloud-provider-azure.sigs.k8s.io/install/configs/.
// Field order follows the upstream documentation: auth configs first, then cluster config.
type azureCloudConfig struct {
	// Auth configs
	TenantID                    string `json:"tenantId,omitempty"`
	SubscriptionID              string `json:"subscriptionId,omitempty"`
	UseManagedIdentityExtension bool   `json:"useManagedIdentityExtension,omitempty"`

	// Cluster config
	ResourceGroup               string `json:"resourceGroup,omitempty"`
	Location                    string `json:"location,omitempty"`
	VnetName                    string `json:"vnetName,omitempty"`
	SubnetName                  string `json:"subnetName,omitempty"`
	SecurityGroupName           string `json:"securityGroupName,omitempty"`
	UseInstanceMetadata         bool   `json:"useInstanceMetadata,omitempty"`
	DisableAvailabilitySetNodes bool   `json:"disableAvailabilitySetNodes,omitempty"`
}

// CloudConfigBuilder creates the cloud configuration file
type CloudConfigBuilder struct {
	*NodeupModelContext
}

var _ fi.NodeupModelBuilder = &CloudConfigBuilder{}

func (b *CloudConfigBuilder) Build(c *fi.NodeupModelBuilderContext) error {
	if !b.HasAPIServer && b.NodeupConfig.KubeletConfig.CloudProvider == "external" {
		return nil
	}

	if err := b.build(c, true); err != nil {
		return err
	}
	if err := b.build(c, false); err != nil {
		return err
	}
	return nil
}

func (b *CloudConfigBuilder) build(c *fi.NodeupModelBuilderContext, inTree bool) error {
	// Add cloud config file if needed
	var lines []string

	cloudProvider := b.CloudProvider()

	var config string
	requireGlobal := true
	switch cloudProvider {
	case kops.CloudProviderGCE:
		if b.NodeupConfig.NodeTags != nil {
			lines = append(lines, "node-tags = "+*b.NodeupConfig.NodeTags)
		}
		if b.NodeupConfig.NodeInstancePrefix != nil {
			lines = append(lines, "node-instance-prefix = "+*b.NodeupConfig.NodeInstancePrefix)
		}
		if b.NodeupConfig.Multizone != nil {
			lines = append(lines, fmt.Sprintf("multizone = %t", *b.NodeupConfig.Multizone))
		}
	case kops.CloudProviderAWS:
		if b.NodeupConfig.DisableSecurityGroupIngress != nil {
			lines = append(lines, fmt.Sprintf("DisableSecurityGroupIngress = %t", *b.NodeupConfig.DisableSecurityGroupIngress))
		}
		if b.NodeupConfig.ElbSecurityGroup != nil {
			lines = append(lines, "ElbSecurityGroup = "+*b.NodeupConfig.ElbSecurityGroup)
		}
		if !inTree {
			for _, family := range b.NodeupConfig.NodeIPFamilies {
				lines = append(lines, "NodeIPFamilies = "+family)
			}
		}
	case kops.CloudProviderOpenstack:
		lines = append(lines, openstack.MakeCloudConfig(b.NodeupConfig.Openstack)...)

	case kops.CloudProviderAzure:
		requireGlobal = false

		vnetName := b.NodeupConfig.Networking.NetworkID
		if vnetName == "" {
			vnetName = b.NodeupConfig.ClusterName
		}

		c := &azureCloudConfig{
			// Auth
			TenantID:                    b.NodeupConfig.AzureTenantID,
			SubscriptionID:              b.NodeupConfig.AzureSubscriptionID,
			UseManagedIdentityExtension: true,
			// Cluster
			ResourceGroup:               b.NodeupConfig.AzureResourceGroup,
			Location:                    b.NodeupConfig.AzureLocation,
			VnetName:                    vnetName,
			SubnetName:                  b.NodeupConfig.AzureSubnetName,
			SecurityGroupName:           b.NodeupConfig.AzureSecurityGroupName,
			UseInstanceMetadata:         true,
			DisableAvailabilitySetNodes: true,
		}
		data, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("error marshalling azure config: %s", err)
		}
		config = string(data)
	}

	if requireGlobal {
		config = "[global]\n" + strings.Join(lines, "\n") + "\n"
	}
	path := CloudConfigFilePath
	if inTree {
		path = InTreeCloudConfigFilePath
	} else if cloudProvider == kops.CloudProviderAzure {
		path = AzureCloudConfigFilePath
	}
	t := &nodetasks.File{
		Path:     path,
		Contents: fi.NewStringResource(config),
		Type:     nodetasks.FileType_File,
	}
	c.AddTask(t)

	return nil
}
