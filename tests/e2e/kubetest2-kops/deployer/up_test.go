/*
Copyright 2021 The Kubernetes Authors.

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

package deployer

import (
	"reflect"
	"testing"
)

func TestAppendIfUnset(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		arg      string
		val      string
		expected []string
	}{
		{
			"empty",
			[]string{},
			"--foo",
			"bar",
			[]string{"--foo", "bar"},
		},
		{
			"unset",
			[]string{"--baz"},
			"--foo",
			"bar",
			[]string{"--baz", "--foo", "bar"},
		},
		{
			"set without value",
			[]string{"--foo"},
			"--foo",
			"bar",
			[]string{"--foo"},
		},
		{
			"set with different value",
			[]string{"--foo", "123"},
			"--foo",
			"bar",
			[]string{"--foo", "123"},
		},
		{
			"set with same value",
			[]string{"--foo", "bar"},
			"--foo",
			"bar",
			[]string{"--foo", "bar"},
		},
		{
			"set with same value and equals sign",
			[]string{"--foo=bar", "--baz=bar"},
			"--foo",
			"bar",
			[]string{"--foo=bar", "--baz=bar"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := appendIfUnset(tc.args, tc.arg, tc.val)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("arguments didn't match: %v vs %v", actual, tc.expected)
			}
		})
	}
}

func TestProwJobLabel(t *testing.T) {
	cases := []struct {
		name          string
		cloudProvider string
		jobName       string
		expected      string
		expectedOK    bool
	}{
		{
			name:          "empty job name omits label",
			cloudProvider: "aws",
			jobName:       "",
			expectedOK:    false,
		},
		{
			name:          "aws preserves slash and dots",
			cloudProvider: "aws",
			jobName:       "pull-kops-e2e-k8s-aws",
			expected:      "prow.k8s.io/job=pull-kops-e2e-k8s-aws",
			expectedOK:    true,
		},
		{
			name:          "azure replaces slash with underscore",
			cloudProvider: "azure",
			jobName:       "pull-kops-e2e-k8s-azure",
			expected:      "prow.k8s.io_job=pull-kops-e2e-k8s-azure",
			expectedOK:    true,
		},
		{
			name:          "gce normalizes key and value",
			cloudProvider: "gce",
			jobName:       "Pull-Kops-E2E-K8s-GCE",
			expected:      "prow_k8s_io_job=pull-kops-e2e-k8s-gce",
			expectedOK:    true,
		},
		{
			name:          "gce truncates value to 63 chars",
			cloudProvider: "gce",
			jobName:       "pull-kops-e2e-kubernetes-aws-canary-very-long-suffix-that-exceeds-the-gce-label-limit",
			expected:      "prow_k8s_io_job=" + "pull-kops-e2e-kubernetes-aws-canary-very-long-suffix-that-excee",
			expectedOK:    true,
		},
		{
			name:          "gce replaces invalid chars with dash",
			cloudProvider: "gce",
			jobName:       "job.name/with:invalid",
			expected:      "prow_k8s_io_job=job-name-with-invalid",
			expectedOK:    true,
		},
		{
			name:          "digitalocean omits label",
			cloudProvider: "digitalocean",
			jobName:       "pull-kops-e2e-k8s-do",
			expectedOK:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := prowJobLabel(tc.cloudProvider, tc.jobName)
			if ok != tc.expectedOK {
				t.Errorf("ok mismatch: got %v, want %v", ok, tc.expectedOK)
			}
			if actual != tc.expected {
				t.Errorf("label mismatch: got %q, want %q", actual, tc.expected)
			}
		})
	}
}
