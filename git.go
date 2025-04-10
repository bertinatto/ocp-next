package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GitWrapper struct {
	workTree string
	gitDir   string
}

func New(dir string) *GitWrapper {
	return &GitWrapper{
		workTree: dir,
		gitDir:   filepath.Join(dir, ".git"),
	}
}

func (g *GitWrapper) runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", "GIT_WORK_TREE", g.workTree))
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", "GIT_DIR", g.gitDir))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git failed to run %q: %w", cmd, err)
	}
	return nil
}

func (g *GitWrapper) Clone(name, url string) error {
	_, err := os.Stat(g.workTree)
	if os.IsNotExist(err) {
		// TODO: switch to a custom runGit with an option to skip setting the env vars
		cmd := exec.Command("git", "clone", "--origin", name, url, g.workTree)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git failed to run %q: %w", cmd, err)
		}
	}
	return nil
}

func (g *GitWrapper) AddRemote(name, url string) error {
	return g.runGit("remote", "add", name, url)
}

func (g *GitWrapper) CreateBranch(name, origin string) error {
	return g.runGit("checkout", "-B", name, origin)
}

func (g *GitWrapper) Merge(branch string) error {
	return g.runGit("merge", "-s", "ours", "--no-edit", branch)
}

func (g *GitWrapper) ApplyPatches(dir string) error {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %q: %w", dir, err)

	}

	fi, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to check directory %q: %w", absPath, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory: %w", absPath, err)
	}

	matches, err := filepath.Glob(filepath.Join(absPath, "*.patch"))
	if err != nil {
		return fmt.Errorf("failed to list patch files: %w", err)
	}

	for _, file := range matches {
		if err := g.Am(file); err != nil {
			return err
		}
		fmt.Printf("Successfully applied %q\n", file)
	}
	return nil
}

func (g *GitWrapper) Add(dir string) error {
	return g.runGit("add", dir)
}

func (g *GitWrapper) Commit(msg string) error {
	return g.runGit("commit", "-m", msg)
}

func (g *GitWrapper) Am(file string) error {
	return g.runGit("am", "--3way", file)
}

func (g *GitWrapper) CheckoutReset() error {
	return g.runGit("checkout", "--", ".")
}

func (g *GitWrapper) Checkout(branch string) error {
	return g.runGit("checkout", branch)
}

func (g *GitWrapper) Clean() error {
	if err := g.runGit("reset", "--hard"); err != nil {
		return err
	}
	return g.runGit("clean", "-fd")
}

func (g *GitWrapper) Fetch(remote, branch string) error {
	return g.runGit("fetch", remote, branch)
}

func (g *GitWrapper) FetchAll() error {
	return g.runGit("fetch")
}

func (g *GitWrapper) Config(name, email string) error {
	if err := g.runGit("config", "--local", "user.name", name); err != nil {
		return err
	}
	return g.runGit("config", "--local", "user.email", email)
}
