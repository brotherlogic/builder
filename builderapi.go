package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/builder/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	// QUEUE - Where we store incoming requests
	QUEUE = "/github.com/brotherlogic/recordadder/queue"
)

var (
	btime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "builder_last_build_time",
		Help: "The size of the print queue",
	}, []string{"job"})
)

// AddRecord adds a record into the system
func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	s.CtxLog(ctx, fmt.Sprintf("Building for %v", req.GetJob()))
	t1 := time.Now()
	err := s.runBuild(ctx, fmt.Sprintf("git@github.com:brotherlogic/%v", req.GetJob()))
	btime.With(prometheus.Labels{"job": req.GetJob()}).Set(float64(time.Since(t1).Seconds()))

	s.CtxLog(ctx, fmt.Sprintf("Build result: %v", err))
	if err != nil {
		s.BounceIssue(ctx, "Refresh Build Error", fmt.Sprintf("%v", err), req.GetJob())
	}
	return &pb.RefreshResponse{}, err
}
