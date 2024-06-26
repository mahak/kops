#!/usr/bin/env bash

# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REPO_ROOT=$(git rev-parse --show-toplevel);
source "${REPO_ROOT}"/tests/e2e/scenarios/lib/common.sh

kops-acquire-latest

OVERRIDES="${OVERRIDES-} --set=cluster.spec.cloudProvider.aws.ebsCSIDriver.enabled=true"
OVERRIDES="$OVERRIDES --set=cluster.spec.snapshotController.enabled=true"
OVERRIDES="$OVERRIDES --set=cluster.spec.certManager.enabled=true"
OVERRIDES="$OVERRIDES --master-size=t3.medium --node-size=c5.large"

kops-up

ZONE=$(${KOPS} get ig -o json | jq -r '[.[] | select(.spec.role=="Node") | .spec.subnets[0]][0]')
REPORT_DIR="${ARTIFACTS:-$(pwd)/_artifacts}/aws-ebs-csi-driver/"

# shellcheck disable=SC2164
cd "$(mktemp -dt kops.XXXXXXXXX)"
go get github.com/onsi/ginkgo/ginkgo

CSI_VERSION=$(kubectl get deployment -n kube-system ebs-csi-controller -o jsonpath='{.spec.template.spec.containers[?(@.name=="ebs-plugin")].image}' | cut -d':' -f2-)
CLONE_ARGS=
if [ -n "$CSI_VERSION" ]; then
    CLONE_ARGS="-b ${CSI_VERSION}"
fi
# shellcheck disable=SC2086
git clone ${CLONE_ARGS} https://github.com/kubernetes-sigs/aws-ebs-csi-driver.git .

# shellcheck disable=SC2164
cd tests/e2e-kubernetes/

ginkgo --nodes=25 ./... -- -cluster-tag="${CLUSTER_NAME}" -ginkgo.skip="\[Disruptive\]" -report-dir="${REPORT_DIR}" -gce-zone="${ZONE}"
