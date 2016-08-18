package main

import "flag"
import "golang.org/x/net/context"
import "google.golang.org/grpc"
import "log"
import "math/rand"
import "strconv"

import "time"

import pb "github.com/brotherlogic/discogssyncer/server"
import pbd "github.com/brotherlogic/godiscogs"
import pbdi "github.com/brotherlogic/discovery/proto"
import pbc "github.com/brotherlogic/cardserver/card"

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, _ := registry.Discover(context.Background(), &entry)
	return r.Ip, int(r.Port)
}

func getRelease(folderName string, host string, port string) *pbd.Release {
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

	return r.Releases[rand.Intn(len(r.Releases))]
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
		log.Printf("RELEASE = %v -> %v", release.InstanceId, release)
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

func scoreCard(releaseID int, rating int, host string, port string) {
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	release := getReleaseWithID("ListeningPile", host, port, releaseID)
	release.Rating = int32(rating)
	// Update the rating and move to the listening box
	log.Printf("UPDATING: %v", release)
	client.UpdateRating(context.Background(), release)
	log.Printf("MOVING: %v to box", release)
	client.MoveToFolder(context.Background(), &pb.ReleaseMove{Release: release, NewFolderId: 673768})
}

func main() {
	var host = flag.String("host", "10.0.1.17", "Hostname of server.")
	var port = flag.String("port", "50055", "Port number of server")
	var dryRun = flag.Bool("dry_run", false, "If true, takes no action")

	flag.Parse()
	portVal, _ := strconv.Atoi(*port)

	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver", *host, portVal)
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	log.Printf("Read cards: %v", cardList)
	if err != nil {
		panic(err)
	}

	foundCard := false
	for _, card := range cardList.Cards {
		if card.Hash == "discogs" {
			foundCard = true
		}

		if card.Hash == "discogs-process" {
			log.Printf("Processing %v", card)

			//delete the card
			server, port := getIP("cardserver", *host, portVal)
			dServer, dPort := getIP("discogssyncer", *host, portVal)

			log.Printf("Scoring card %v", card)
			releaseID, _ := strconv.Atoi(card.Text)
			rating, _ := strconv.Atoi(card.ActionMetadata[0])
			log.Printf("Scoring %v as %v", releaseID, rating)
			if !*dryRun {
				scoreCard(releaseID, rating, dServer, strconv.Itoa(dPort))
			}

			log.Printf("Deleting %v", card.Hash)
			if !*dryRun {
				deleteCard(card.Hash, server, strconv.Itoa(port))
			}

		}

		if card.Hash == "" {
			log.Printf("OOPS %v", card)
			//delete the card
			server, port := getIP("cardserver", *host, portVal)
			deleteCard(card.Hash, server, strconv.Itoa(port))
		}
	}

	if !foundCard {
		dServer, dPort := getIP("discogssyncer", *host, portVal)
		rel := getRelease("ListeningPile", dServer, strconv.Itoa(dPort))

		cards := pbc.CardList{}

		imageURL := ""
		backupURL := ""
		for _, image := range rel.Images {
			if image.Type == "primary" {
				imageURL = image.Uri
			}
			backupURL = image.Uri
		}
		if imageURL == "" {
			imageURL = backupURL
		}

		cardResponse := &pbc.Card{Hash: "discogs-process", Priority: -10, Text: strconv.Itoa(int(rel.Id))}
		card := pbc.Card{Text: pbd.GetReleaseArtist(*rel) + " - " + rel.Title, Hash: "discogs", Image: imageURL, Action: pbc.Card_RATE, Priority: 100, Result: cardResponse}
		cards.Cards = append(cards.Cards, &card)
		log.Printf("Writing: %v", cards)
		if !*dryRun {
			log.Printf("Writing the card")
			_, err = client.AddCards(context.Background(), &cards)
			if err != nil {
				log.Printf("Problem adding cards %v", err)
			}
		}
	}
}
