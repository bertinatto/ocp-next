package main

import (
	"os"
	"os/exec"
)

func main() {
	git := New("kubernetes")

	if err := git.Clone("upstream", "https://github.com/kubernetes/kubernetes.git"); err != nil {
		panic(err)
	}

	if err := git.Config("Fabio Bertinatto", "fbertina@redhat.com"); err != nil {
		panic(err)
	}

	git.AddRemote("ocp", "https://github.com/openshift/kubernetes.git")
	git.AddRemote("origin", "git@github.com:bertinatto/kubernetes.git")

	if err := git.Fetch("upstream", ""); err != nil {
		panic(err)
	}

	if err := git.Fetch("ocp", "master"); err != nil {
		panic(err)
	}

	if err := git.Clean(); err != nil {
		panic(err)
	}
	if err := git.CheckoutReset(); err != nil {
		panic(err)
	}

	if err := git.CreateBranch("ocp-next", "upstream/release-1.32"); err != nil {
		// if err := git.CreateBranch("ocp-next", "v1.32.0-rc.0"); err != nil {
		panic(err)
	}

	err := git.Merge("ocp/master")
	if err != nil {
		panic(err)
	}

	if err := git.ApplyPatches("patches"); err != nil {
		panic(err)
	}

	err = updateDependencies(git.workTree)
	if err != nil {
		panic(err)
	}

	err = updateVendor(git.workTree)
	if err != nil {
		panic(err)
	}

	if err := git.Add("."); err != nil {
		panic(err)
	}

	if err := git.Commit("UPSTREAM: <drop>: hack/update-vendor.sh"); err != nil {
		panic(err)
	}

	// TODO: run this inside container
	cmd := exec.Command("kubernetes/hack/install-etcd.sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	err = updateGenerated(git.workTree)
	if err != nil {
		panic(err)
	}

	if err := git.Add("."); err != nil {
		panic(err)
	}

	if err := git.Commit("UPSTREAM: <drop>: make update"); err != nil {
		panic(err)
	}

	err = updateVendor(git.workTree)
	if err != nil {
		panic(err)
	}

	if err := git.Add("."); err != nil {
		panic(err)
	}

	if err := git.Commit("UPSTREAM: <drop>: hack/update-vendor.sh"); err != nil {
		panic(err)
	}
}
