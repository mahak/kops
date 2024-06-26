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

package mockec2

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"k8s.io/kops/util/pkg/awsinterfaces"
)

type MockEC2 struct {
	// Stub out interface
	awsinterfaces.EC2API

	mutex sync.Mutex

	addressNumber int
	Addresses     map[string]*ec2types.Address

	RouteTables map[string]*ec2types.RouteTable

	DhcpOptions map[string]*ec2types.DhcpOptions

	Images []*ec2types.Image

	securityGroupNumber int
	SecurityGroups      map[string]*ec2types.SecurityGroup
	SecurityGroupRules  map[string]*ec2types.SecurityGroupRule

	subnets map[string]*subnetInfo

	Volumes map[string]*ec2types.Volume

	KeyPairs map[string]*ec2types.KeyPairInfo

	Tags []*ec2types.TagDescription

	Vpcs map[string]*vpcInfo

	InternetGateways           map[string]*ec2types.InternetGateway
	EgressOnlyInternetGateways map[string]*ec2types.EgressOnlyInternetGateway

	launchTemplateNumber int
	LaunchTemplates      map[string]*launchTemplateInfo

	NatGateways map[string]*ec2types.NatGateway

	idsMutex sync.Mutex
	ids      map[string]*idAllocator
}

var _ awsinterfaces.EC2API = &MockEC2{}

func (m *MockEC2) All() map[string]interface{} {
	all := make(map[string]interface{})

	for _, o := range m.Addresses {
		all[aws.ToString(o.AllocationId)] = o
	}
	for id, o := range m.RouteTables {
		all[id] = o
	}
	for id, o := range m.DhcpOptions {
		all[id] = o
	}
	for _, o := range m.Images {
		all[aws.ToString(o.ImageId)] = o
	}
	for id, o := range m.SecurityGroups {
		all[id] = o
	}
	for id, o := range m.subnets {
		all[id] = &o.main
	}
	for id, o := range m.Volumes {
		all[id] = o
	}
	for id, o := range m.KeyPairs {
		all[id] = o
	}
	for id, o := range m.Vpcs {
		all[id] = o
	}
	for id, o := range m.InternetGateways {
		all[id] = o
	}
	for id, o := range m.EgressOnlyInternetGateways {
		all[id] = o
	}
	for id, o := range m.LaunchTemplates {
		all[id] = o
	}
	for id, o := range m.NatGateways {
		all[id] = o
	}

	return all
}

type idAllocator struct {
	NextId int
}

func (m *MockEC2) allocateId(prefix string) string {
	m.idsMutex.Lock()
	defer m.idsMutex.Unlock()

	ids := m.ids[prefix]
	if ids == nil {
		if m.ids == nil {
			m.ids = make(map[string]*idAllocator)
		}
		ids = &idAllocator{NextId: 1}
		m.ids[prefix] = ids
	}
	id := ids.NextId
	ids.NextId++
	return fmt.Sprintf("%s-%d", prefix, id)
}
