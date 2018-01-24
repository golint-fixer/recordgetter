package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
	"github.com/brotherlogic/goserver/utils"
	pbrg "github.com/brotherlogic/recordgetter/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func findServer(name string) (string, int) {
	conn, err := grpc.Dial(utils.Discover, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot reach discover server: %v (trying to discover %v)", err, name)
	}
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	re := &pbdi.RegistryEntry{Name: name}
	r, err := registry.Discover(ctx, re)

	e, ok := status.FromError(err)
	if ok && e.Code() == codes.Unavailable {
		log.Printf("RETRY")
		r, err = registry.Discover(ctx, re)
	}

	if err != nil {
		return "", -1
	}
	return r.Ip, int(r.Port)
}

func clear() {
	host, port := findServer("recordgetter")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbrg.NewRecordGetterClient(conn)
	r, err := client.Force(context.Background(), &pbrg.Empty{})
	fmt.Printf("%v and %v", r, err)
}

func listened(score int32) {
	host, port := findServer("recordgetter")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbrg.NewRecordGetterClient(conn)
	r, err := client.GetRecord(context.Background(), &pbrg.GetRecordRequest{})
	if err != nil {
		log.Fatalf("%v", err)
	}
	r.GetRecord().GetRelease().Rating = score
	_, err = client.Listened(context.Background(), r.GetRecord())
	fmt.Printf("%v", err)
}

func get() {
	host, port := findServer("recordgetter")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbrg.NewRecordGetterClient(conn)
	r, err := client.GetRecord(context.Background(), &pbrg.GetRecordRequest{Refresh: true})
	fmt.Printf("%v and %v", r, err)
}

func score(value int32) {
	host, port := findServer("recordgetter")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbrg.NewRecordGetterClient(conn)
	r, err := client.GetRecord(context.Background(), &pbrg.GetRecordRequest{})
	if err != nil {
		log.Fatalf("Error in scoring: %v", err)
	}
	r.GetRecord().GetMetadata().SetRating = value
	re, err := client.Listened(context.Background(), r.GetRecord())
	fmt.Printf("%v and %v", re, err)
}

func run() (int, error) {
	host, port := findServer("discogssyncer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: "ListeningPile"}
	folderList.Folders = append(folderList.Folders, folder)
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		return 0, err
	}

	return len(r.GetRecords()), nil
}

func main() {
	get()
	val, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing num: %v", err)
	}
	score(int32(val))
}
