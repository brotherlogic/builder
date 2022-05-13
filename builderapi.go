package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/builder/proto"
)

const (
	// QUEUE - Where we store incoming requests
	QUEUE = "/github.com/brotherlogic/recordadder/queue"
)

//AddRecord adds a record into the system
func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return &pb.RefreshResponse{}, s.runBuild(ctx, fmt.Sprintf("git@github.com:brotherlogic/%v", req.GetJob()))
}
