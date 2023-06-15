#!/bin/bash

LOG="$(realpath output.log)"

# First, import repos
git clone --origin upstream https://github.com/kubernetes/kubernetes.git
pushd kubernetes || exit
git remote add openshift --fetch https://github.com/openshift/kubernetes.git

# Create a "ocp-next" branch and make sure it's clean
git checkout -B ocp-next upstream/master
git clean -fd
git checkout -- .

# Create a merge commit to bring the openshift changes to our branch
git merge -s ours --no-edit openshift/master

# Apply each patch individually. If one of them fail, abort immediately
echo "Starting: $(date -u +%Y-%m-%dT%H:%M:%S%Z)" > "$LOG"
for patch in ../patches/*.patch; do
    if ! git am "$patch" >> "$LOG" 2>&1; then
        echo "Failed to apply $patch. Check the log at $LOG"
        exit 1
    else
        echo "Applied $patch"
    fi
done

# No patches failed, so update the dependencies

GOPROXY=direct hack/pin-dependency.sh github.com/onsi/ginkgo/v2=github.com/openshift/onsi-ginkgo/v2 v2.9-openshift-4.14
GOPROXY=direct hack/pin-dependency.sh github.com/openshift/api=github.com/bertinatto/api ocp-next
GOPROXY=direct hack/pin-dependency.sh github.com/openshift/client-go=github.com/bertinatto/client-go ocp-next
GOPROXY=direct hack/pin-dependency.sh github.com/openshift/library-go=github.com/bertinatto/library-go ocp-next
GOPROXY=direct hack/pin-dependency.sh github.com/openshift/apiserver-library-go=github.com/bertinatto/apiserver-library-go ocp-next

hack/update-vendor.sh

export PATH="/home/fbertina/src/k8s.io/kubernetes/third_party/etcd:${PATH}"
sudo PATH="$PATH" make update
