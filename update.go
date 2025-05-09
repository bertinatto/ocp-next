package main

import (
	"fmt"
	"os"
	"os/exec"
)

func updateDependencies(workTree string) error {
	deps := [][]string{
		{"github.com/onsi/ginkgo/v2=github.com/openshift/onsi-ginkgo/v2", "v2.21-openshift-4.20"},
		{"github.com/openshift/api=github.com/bertinatto/api", "bump-v1.33"},
		{"github.com/openshift/client-go=github.com/bertinatto/client-go", "bump-v1.33"},
		{"github.com/openshift/library-go=github.com/bertinatto/library-go", "bump-v1.33"},
		{"github.com/openshift/apiserver-library-go=github.com/bertinatto/apiserver-library-go", "bump-v1.33"},
		// FIXME
		// {"github.com/openshift/api", "latest"},
		// {"github.com/openshift/client-go", "latest"},
		// {"github.com/openshift/library-go", "latest"},
		// {"github.com/openshift/apiserver-library-go", "latest"},
		// Kevin
		// {"github.com/openshift/api=github.com/kevinrizza/api-1", "update-to-1.33"},
		// {"github.com/openshift/client-go=github.com/kevinrizza/client-go", "update-kube-1.33"},
		// {"github.com/openshift/library-go=github.com/kevinrizza/library-go", "1.33-kube-update"},
		// {"github.com/openshift/apiserver-library-go=github.com/kevinrizza/apiserver-library-go", "update-kube-1.33"},
	}

	for i := range deps {
		dep := deps[i][0]
		version := deps[i][1]
		cmd := exec.Command("podman",
			"run",
			"-it",
			"--rm",
			"-v",
			".:/go/k8s.io/kubernetes:Z",
			"--workdir=/go/k8s.io/kubernetes",
			"--env", "GOPROXY=direct",
			"--env", "OS_RUN_WITHOUT_DOCKER=yes",
			"--env", "FORCE_HOST_GO=1",
			"registry.ci.openshift.org/openshift/release:rhel-9-release-golang-1.24-nofips-openshift-4.19",
			"hack/pin-dependency.sh",
			dep,
			version,
		)
		cmd.Dir = workTree
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run %q: %w", cmd, err)
		}
	}
	return nil
}

func updateVendor(workTree string) error {
	cmd := exec.Command("podman",
		"run",
		"-it",
		"--rm",
		"-v",
		".:/go/k8s.io/kubernetes:Z",
		"--workdir=/go/k8s.io/kubernetes",
		// "--env", "GOPROXY=direct",
		"--env", "OS_RUN_WITHOUT_DOCKER=yes",
		"--env", "FORCE_HOST_GO=1",
		"registry.ci.openshift.org/openshift/release:rhel-9-release-golang-1.24-nofips-openshift-4.19",
		"hack/update-vendor.sh")
	cmd.Dir = workTree
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run %q: %w", cmd, err)
	}
	return nil
}

func updateGenerated(workTree string) error {
	cmds := []*exec.Cmd{
		// exec.Command("podman",
		// 	"run",
		// 	"-it",
		// 	"--rm",
		// 	"-v",
		// 	".:/go/k8s.io/kubernetes:Z",
		// 	"--workdir=/go/k8s.io/kubernetes",
		// 	"--env", "OS_RUN_WITHOUT_DOCKER=yes",
		// 	"--env", "FORCE_HOST_GO=1",
		// 	"registry.ci.openshift.org/openshift/release:rhel-9-release-golang-1.21-openshift-4.16",
		// 	"hack/install-etcd.sh"),
		exec.Command("podman",
			"run",
			"-it",
			"--rm",
			"-v",
			".:/go/k8s.io/kubernetes:Z",
			"--workdir=/go/k8s.io/kubernetes",
			"--env", "OS_RUN_WITHOUT_DOCKER=yes",
			"--env", "FORCE_HOST_GO=1",
			"--env", "PATH=/opt/google/protobuf/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/go/bin:/go/k8s.io/kubernetes/third_party/etcd:/go/k8s.io/kubernetes/third_party/protoc",
			"registry.ci.openshift.org/openshift/release:rhel-9-release-golang-1.24-nofips-openshift-4.19",
			"make", "update"),
	}
	for _, cmd := range cmds {
		cmd.Dir = workTree
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run %q: %w", cmd, err)
		}
	}
	return nil
}
