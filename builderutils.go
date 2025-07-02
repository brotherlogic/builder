package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	WORKING_DIR = "/media/scratch/goscratch"
	NEW_GO      = "1.22"
)

var (
	OLD_GO = []string{"1.19.2", "1.19", "1.20"}
)

func (s *Server) runBuild(ctx context.Context, gha string) error {
	// Only build one thing at a time
	s.lock.Lock()
	defer s.lock.Unlock()

	// Track the actual build time
	t1 := time.Now()
	defer btime.With(prometheus.Labels{"job": gha}).Set(float64(time.Since(t1).Seconds()))

	err := os.MkdirAll(WORKING_DIR, 0700)
	if err != nil {
		return status.Errorf(codes.AlreadyExists, "Cannot make scratch dir: %v", err)
	}
	defer os.RemoveAll(WORKING_DIR)

	os.Chdir(WORKING_DIR)

	out1, err := exec.Command("git", "clone", gha, "./").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "clone (%v) %v -> %v", s.Registry.Identifier, err, string(out1))
	}

	// Add the PR closer file
	_, err = exec.Command("curl", "https://raw.githubusercontent.com/brotherlogic/discovery/main/clean_branches.sh", "-o", "./clean_branches.sh").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Failed download: %v", err)
	}
	_, err = exec.Command("chmod", "u+x", "./clean_branches.sh").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Failed chmod: %v", err)
	}
	_, err = exec.Command("./clean_branches.sh").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Failed clean: %v", err)
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

	for _, og := range OLD_GO {
		out5, err := exec.Command("sed", "-i", "-e", fmt.Sprintf("s/%v/%v/g", og, NEW_GO), ".github/workflows/basicrun.yml").CombinedOutput()
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "awk (%v) -> %v, %v", s.Registry.Identifier, err, string(out5))
		}
	}

	// Add the PR closer file
	out7, err := exec.Command("curl", "https://raw.githubusercontent.com/brotherlogic/discovery/main/.github/workflows/close.yml", "-o", ".github/workflows/close.yml").CombinedOutput()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Failed download: %v -> %v", err, string(out7))
	}

	out9, err := exec.Command("git", "add", ".github/workflows/close.yml").CombinedOutput()
	s.CtxLog(ctx, fmt.Sprintf("Here -> %v => %v", string(out9), err))

	// Ensure we add clean branches if not present
	out10, err := exec.Command("git", "add", "clean_branches.sh").CombinedOutput()
	s.CtxLog(ctx, fmt.Sprintf("Here -> %v => %v", string(out10), err))

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
