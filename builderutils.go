package main

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/net/context"
)

const (
	WORKING_DIR = "/media/scratch/goscratch"
)

func (s *Server) runBuild(ctx context.Context, gha string) error {
	os.MkdirAll(WORKING_DIR, 0700)
	defer os.RemoveAll(WORKING_DIR)

	os.Chdir(WORKING_DIR)

	out1, err := exec.Command("git", "clone", gha, "./").CombinedOutput()
	if err != nil {
		return fmt.Errorf("(%v) %v -> %v", s.Registry.Identifier, err, string(out1))
	}

	out2, err := exec.Command("go", "get", "-u", "./...").CombinedOutput()
	if err != nil {
		return fmt.Errorf("(%v) %v -> %v", s.Registry.Identifier, err, string(out2))
	}

	out3, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return fmt.Errorf("(%v) %v -> %v", s.Registry.Identifier, err, string(out3))
	}

	out4, err1 := exec.Command("git", "push", "origin", "main").CombinedOutput()
	out5, err2 := exec.Command("git", "push", "origin", "master").CombinedOutput()

	if err1 != nil && err2 != nil {
		return fmt.Errorf("(%v) Unable to push: %v or %v -> %v, %v", s.Registry.Identifier, err1, err2, out4, out5)
	}

	return nil
}
