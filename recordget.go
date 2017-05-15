package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/discogssyncer/server"

	pbdi "github.com/brotherlogic/discovery/proto"
)

import "google.golang.org/grpc"
import "google.golang.org/grpc/grpclog"

import "math/rand"

import "time"

import pbd "github.com/brotherlogic/godiscogs"

import pbc "github.com/brotherlogic/cardserver/card"

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, err := registry.Discover(context.Background(), &entry)
	if err != nil {
		log.Printf("Error discovering %v -> %v", servername, err)
		return "", -1
	}

	log.Printf("Found %v -> %v:%v", servername, r.Ip, r.Port)
	return r.Ip, int(r.Port)
}

func getReleaseFromPile(folderName string, host string, port string) (*pbd.Release, *pb.ReleaseMetadata) {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatalf("Problem getting releases %v", err)
	}

	if len(r.Releases) == 0 {
		log.Printf("No releases in folder: %v", folderList)
		return nil, nil
	}

	var newRel *pbd.Release
	newRel = nil
	for _, rel := range r.Releases {
		meta, err2 := client.GetMetadata(context.Background(), rel)
		if err2 == nil {
			if meta.DateAdded > (time.Now().AddDate(0, -3, 0).Unix()) {
				newRel = rel
			}
		}
	}

	if newRel == nil {
		newRel = r.Releases[rand.Intn(len(r.Releases))]
	}
	meta, err := client.GetMetadata(context.Background(), newRel)
	if err != nil {
		log.Fatalf("Problem getting metadata %v", err)
	}
	return newRel, meta
}

func getReleaseFromCollection(host string, port string, allowSeven bool) (*pbd.Release, *pb.ReleaseMetadata) {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)

	folderList := &pb.FolderList{}
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "12s"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "10s"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "April Orchestra"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Death Waltz"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "IM"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Music Mosaic"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "MusiquePourLImage"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "NumeroLPs"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Outside"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Robbie Basho"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Timing"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "TVMusic"})
	folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "Vinyl Boxsets"})
	if allowSeven {
		folderList.Folders = append(folderList.Folders, &pbd.Folder{Name: "7s"})
	}
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatalf("Problem getting releases %v", err)
	}

	log.Printf("Trying to get from %v: %v", len(r.Releases), r.Releases)
	retRel := r.Releases[rand.Intn(len(r.Releases))]
	meta, err := client.GetMetadata(context.Background(), retRel)
	if err != nil {
		log.Fatalf("Problem getting metadata %v", err)
	}
	return retRel, meta
}

func getReleaseWithID(folderName string, host string, port string, id int) *pbd.Release {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatalf("Problem getting releases %v", err)
	}

	for _, release := range r.Releases {
		if int(release.Id) == id {
			return release
		}
	}

	return nil
}

func deleteCard(hash string, host string, port string) {
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	client.DeleteCards(context.Background(), &pbc.DeleteRequest{Hash: hash})
}

func scoreCard(releaseID int, rating int, host string, port string) bool {
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	allowSeven := true
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	release := getReleaseWithID("ListeningPile", host, port, releaseID)
	if release == nil {
		release = getReleaseWithID("7s", host, port, releaseID)
		allowSeven = false
	}
	if release != nil {
		release.Rating = int32(rating)
		// Update the rating and move to the listening box
		if rating > 0 {
			client.UpdateRating(context.Background(), release)
		}
		client.MoveToFolder(context.Background(), &pb.ReleaseMove{Release: release, NewFolderId: 673768})
	}
	return allowSeven
}

func hasCurrentCard(host string, portVal int) bool {
	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver", host, portVal)
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	if err != nil {
		log.Fatalf("Whoops2: %v", err)
	}

	for _, card := range cardList.Cards {
		if card.Hash == "discogs" {
			return true
		}
	}
	return false
}

func addCards(cardList *pbc.CardList, host string, portVal int) {
	cServer, cPort := getIP("cardserver", host, portVal)
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	client.AddCards(context.Background(), cardList)
}

func processCard(host string, portVal int, dryRun bool) bool {
	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver", host, portVal)
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	allowSeven := true

	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	if err != nil {
		panic(err)
	}

	for _, card := range cardList.Cards {
		if card.Hash == "discogs-process" {

			//delete the card
			server, port := getIP("cardserver", host, portVal)
			dServer, dPort := getIP("discogssyncer", host, portVal)

			releaseID, _ := strconv.Atoi(card.Text)
			if card.ActionMetadata != nil {
				rating, _ := strconv.Atoi(card.ActionMetadata[0])
				if !dryRun {
					allowSeven = scoreCard(releaseID, rating, dServer, strconv.Itoa(dPort))
				}
			} else {
				if !dryRun {
					allowSeven = scoreCard(releaseID, -1, dServer, strconv.Itoa(dPort))
				}
			}
			if !dryRun {
				deleteCard(card.Hash, server, strconv.Itoa(port))
			}
		}
	}

	return allowSeven
}

func getCard(rel *pbd.Release) pbc.Card {
	var imageURL string
	var backupURL string
	for _, image := range rel.Images {
		if image.Type == "primary" {
			imageURL = image.Uri
		}
		backupURL = image.Uri
	}
	if imageURL == "" {
		imageURL = backupURL
	}

	card := pbc.Card{Text: pbd.GetReleaseArtist(*rel) + " - " + rel.Title, Hash: "discogs", Image: imageURL, Priority: 100}
	return card
}

func main() {
	var host = flag.String("host", "192.168.86.34", "Hostname of server.")
	var port = flag.Int("port", 50055, "Port number of server")
	var dryRun = flag.Bool("dry_run", false, "If true, takes no action")
	var quiet = flag.Bool("quiet", false, "Don't log anything.")
	flag.Parse()

	if *quiet {
		log.SetOutput(ioutil.Discard)
		grpclog.SetLogger(log.New(ioutil.Discard, "", -1))
	}

	foundCard := hasCurrentCard(*host, *port)
	allowSeven := processCard(*host, *port, *dryRun)
	cards := pbc.CardList{}

	if !foundCard {
		dServer, dPort := getIP("discogssyncer", *host, *port)
		rel, meta := getReleaseFromPile("ListeningPile", dServer, strconv.Itoa(dPort))

		if rel != nil {
			card := getCard(rel)
			card.Result = &pbc.Card{Hash: "discogs-process", Priority: -10, Text: strconv.Itoa(int(rel.Id))}
			card.Action = pbc.Card_RATE
			addTime := time.Unix(meta.DateAdded, 0)
			if time.Now().Sub(addTime).Hours() < 24*30*3 {
				card.Action = pbc.Card_DISMISS
			}
			cards.Cards = append(cards.Cards, &card)
		} else {
			rel, _ := getReleaseFromCollection(dServer, strconv.Itoa(dPort), allowSeven)
			card := getCard(rel)
			card.Action = pbc.Card_DISMISS
			if rel.Rating <= 0 {
				card.Result = &pbc.Card{Hash: "discogs-process", Priority: -10, Text: strconv.Itoa(int(rel.Id))}
				card.Action = pbc.Card_RATE
			}
			cards.Cards = append(cards.Cards, &card)
		}
	}
	if !*dryRun {
		addCards(&cards, *host, *port)
	}
}
