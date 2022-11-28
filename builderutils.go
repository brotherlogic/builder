package main

import (
	"os"
	"os/exec"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	WORKING_DIR = "/media/scratch/goscratch"
)

func (s *Server) runBuild(ctx context.Context, gha string) error {
	// Only build one thing at a time
	s.lock.Lock()
	defer s.lock.Unlock()

	os.MkdirAll(WORKING_DIR, 0700)
	defer os.RemoveAll(WORKING_DIR)

	os.Chdir(WORKING_DIR)

	out1, err := exec.Command("git", "clone", gha, "./").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "clone (%v) %v -> %v", s.Registry.Identifier, err, string(out1))
	}

	out2, err := exec.Command("go", "get", "-u", "./...").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "go get (%v) %v -> %v", s.Registry.Identifier, err, string(out2))
	}

	out3, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "go mod (%v) %v -> %v", s.Registry.Identifier, err, string(out3))
	}

	out6, err := exec.Command("git", "commit", "-am", "DownstreamUpdates").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "commit (%v) %v -> %v", s.Registry.Identifier, err, string(out6))
	}

	out4, err1 := exec.Command("git", "push", "origin", "main").CombinedOutput()
	out5, err2 := exec.Command("git", "push", "origin", "master").CombinedOutput()

	if err1 != nil && err2 != nil {
		return status.Errorf(codes.FailedPrecondition, "(%v) Unable to push: %v or %v -> %v, %v", s.Registry.Identifier, err1, err2, string(out4), string(out5))
	}

	return nil
}
