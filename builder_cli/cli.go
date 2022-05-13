package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/builder/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.Register(&utils.DiscoveryClientResolverBuilder{})
}

func main() {
	ctx, cancel := utils.ManualContext("builder-cli", time.Second*10)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "builder")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewBuildClient(conn)

	switch os.Args[1] {
	case "build":
		buildFlags := flag.NewFlagSet("Build", flag.ExitOnError)
		var name = buildFlags.String("name", "", "Id of the record to add")

		if err := buildFlags.Parse(os.Args[2:]); err == nil {
			if len(*name) > 0 {
				_, err := client.Refresh(ctx, &pb.RefreshRequest{Job: *name})
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
			}
		}
	}
}
