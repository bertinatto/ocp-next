package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	base := flag.String("base", "", "Upstream branch to start with")
	onlyPatches := flag.Bool("only-patches", false, "Only apply patches, no dependency bump")
	flag.Parse()

	if *base == "" {
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "ERROR: base cannot be empty.\n")
		os.Exit(1)
	}

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

	if err := git.CreateBranch("ocp-next", *base); err != nil {
		panic(err)
	}

	err := git.Merge("ocp/master")
	if err != nil {
		panic(err)
	}

	if err := git.ApplyPatches("patches"); err != nil {
		panic(err)
	}

	if *onlyPatches {
		os.Exit(0)
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

	// Apparently I need to do this twice
	for range 2 {
		err = updateVendor(git.workTree)
		if err != nil {
			panic(err)
		}
	}

	if err := git.Add("."); err != nil {
		panic(err)
	}

	if err := git.Commit("UPSTREAM: <drop>: hack/update-vendor.sh"); err != nil {
		panic(err)
	}
}
