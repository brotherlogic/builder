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

	err := exec.Command("git", "clone", "./", gha).Run()
	if err != nil {
		return err
	}

	err = exec.Command("go", "get", "-u", "./...").Run()
	if err != nil {
		return err
	}

	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}

	err1 := exec.Command("git", "push", "origin", "main").Run()
	err2 := exec.Command("git", "push", "origin", "master").Run()

	if err1 != nil && err2 != nil {
		return fmt.Errorf("Unable to push: %v or %v", err1, err2)
	}

	return nil
}
