package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/brotherlogic/goserver/utils"

	pb "github.com/brotherlogic/builder/proto"
	dpb "github.com/brotherlogic/discovery/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.ManualContext("builder-cli", time.Minute*10)
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "builder")
	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewBuildClient(conn)

	switch os.Args[1] {
	case "fullbuild":
		ctx, cancel := utils.ManualContext("builder-cli", time.Hour*2)
		defer cancel()
		conn2, err2 := utils.LFDial(utils.Discover)
		if err2 != nil {
			log.Fatalf("Unable to dial: %v", err2)
		}
		dclient := dpb.NewDiscoveryServiceV2Client(conn2)
		alljobs, err := dclient.Get(ctx, &dpb.GetRequest{})
		if err != nil {
			log.Fatalf("All jobs request failed: %v", err)
		}

		jobm := make(map[string]bool)
		for _, j := range alljobs.GetServices() {
			jobm[j.GetName()] = true
		}

		wg := &sync.WaitGroup{}
		for j := range jobm {
			fmt.Printf("Building %v\n", j)
			wg.Add(1)
			go func(job string) {
				_, err := client.Refresh(ctx, &pb.RefreshRequest{Job: job})
				fmt.Printf("Built %v 0> %v\n", job, err)
				wg.Done()
			}(j)
		}
		wg.Wait()
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
