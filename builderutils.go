package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	WORKING_DIR = "/media/scratch/goscratch"
	OLD_GO      = "1.18.1"
	NEW_GO      = "1.19"
)

func (s *Server) runBuild(ctx context.Context, gha string) error {
	// Only build one thing at a time
	s.lock.Lock()
	defer s.lock.Unlock()

	err := os.MkdirAll(WORKING_DIR, 0700)
	if err != nil {
		return status.Errorf(codes.AlreadyExists, fmt.Sprintf("Cannot make scratch dir: %v", err))
	}
	defer os.RemoveAll(WORKING_DIR)

	os.Chdir(WORKING_DIR)

	out1, err := exec.Command("git", "clone", gha, "./").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "clone (%v) %v -> %v", s.Registry.Identifier, err, string(out1))
	}

	branch := fmt.Sprintf("update-%v", time.Now().Unix())
	out1a, err := exec.Command("git", "checkout", "-b", branch).CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "checkout (%v) %v -> %v", s.Registry.Identifier, err, string(out1a))
	}

	out2, err := exec.Command("go", "get", "-u", "./...").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "go get (%v) %v -> %v", s.Registry.Identifier, err, string(out2))
	}

	out3, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "go mod (%v) %v -> %v", s.Registry.Identifier, err, string(out3))
	}

	// Add the PR closer file
	out5, err := exec.Command("awk", fmt.Sprintf("{gsub(/%v/, \"%v\")}", OLD_GO, NEW_GO), ".github/workflows/*.yml").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "awk (%v) -> %v, %v", s.Registry.Identifier, err, string(out5))
	}

	out6, err := exec.Command("git", "commit", "-am", "DownstreamUpdates").CombinedOutput()
	if err != nil {
		if !strings.Contains(string(out6), "nothing to commit") {
			return status.Errorf(codes.FailedPrecondition, "commit (%v) %v -> %v", s.Registry.Identifier, err, string(out6))
		}
	}

	out4, err1 := exec.Command("git", "push", "origin", branch).CombinedOutput()

	if err1 != nil {
		return status.Errorf(codes.FailedPrecondition, "(%v) Unable to push: %v -> %v", s.Registry.Identifier, err1, string(out4))
	}

	return nil
}
