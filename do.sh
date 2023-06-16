#!/bin/bash

set -eo pipefail

main() {
    local branch_name="$1"

    clone_repos
    pushd kubernetes || exit 1
    create_branch "$branch_name"
    merge_changes
    apply_patches
    update_dependencies
    update_vendor
    # FIXME
    # update_and_build
    # commit_and_push
    popd
}

clone_repos() {
    git clone --origin upstream "https://github.com/kubernetes/kubernetes.git" || true
    git remote add openshift --fetch "https://github.com/openshift/kubernetes.git" || true
}

create_branch() {
    local branch_name="$1"
    git clean -fd
    git checkout -- .
    git checkout -B "$branch_name" upstream/master
}

merge_changes() {
    git merge -s ours --no-edit openshift/master
}

apply_patches() {
    local patch
    echo "Starting to apply patches"
    for patch in ../patches/*.patch; do
        if ! git am "$patch"; then
            echo "Failed to apply $patch"
	    exit 1
        else
            echo "Applied $patch"
        fi
    done
}

update_dependencies() {
    GOPROXY=direct hack/pin-dependency.sh github.com/onsi/ginkgo/v2=github.com/openshift/onsi-ginkgo/v2 v2.9-openshift-4.14
    GOPROXY=direct hack/pin-dependency.sh github.com/openshift/api=github.com/bertinatto/api ocp-next
    GOPROXY=direct hack/pin-dependency.sh github.com/openshift/client-go=github.com/bertinatto/client-go ocp-next
    GOPROXY=direct hack/pin-dependency.sh github.com/openshift/library-go=github.com/bertinatto/library-go ocp-next
    GOPROXY=direct hack/pin-dependency.sh github.com/openshift/apiserver-library-go=github.com/bertinatto/apiserver-library-go ocp-next
}

update_vendor() {
    hack/update-vendor.sh
}

update_and_build() {
    eval "$(hack/install-etcd.sh | grep "export PATH")"
    # FIXME: this is for local testing
    sudo PATH="$PATH" make clean
    sudo PATH="$PATH" make update
}

commit_and_push() {
    git add .
    git commit -m "UPSTREAM: <drop>: update dependencies and generated files"
    git remote add origin --fetch "https://github.com/bertinatto/kubernetes.git" || true
    git push origin HEAD -f
}

if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
    main "$@"
fi
