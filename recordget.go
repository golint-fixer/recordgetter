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

func getRelease(folderName string, host string, port string) (*pbd.Release, *pb.ReleaseMetadata) {
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
	if rating > 0 {
		client.UpdateRating(context.Background(), release)
	}
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
	if err != nil {
		panic(err)
	}

	foundCard := false
	for _, card := range cardList.Cards {
		if card.Hash == "discogs" {
			foundCard = true
		}

		if card.Hash == "discogs-process" {

			//delete the card
			server, port := getIP("cardserver", *host, portVal)
			dServer, dPort := getIP("discogssyncer", *host, portVal)

			log.Printf("READ CARD %v", card)

			releaseID, _ := strconv.Atoi(card.Text)
			if card.ActionMetadata != nil {
				rating, _ := strconv.Atoi(card.ActionMetadata[0])
				if !*dryRun {

					log.Printf("RATING CARD %v", card)
					scoreCard(releaseID, rating, dServer, strconv.Itoa(dPort))
				}
			} else {
				if !*dryRun {
					log.Printf("SCORING CARD %v", card)
					scoreCard(releaseID, -1, dServer, strconv.Itoa(dPort))
				}
			}

			if !*dryRun {
				log.Printf("DELETING CARD %v", card)
				deleteCard(card.Hash, server, strconv.Itoa(port))
			}

		}

		if card.Hash == "" {
			//delete the card
			server, port := getIP("cardserver", *host, portVal)
			deleteCard(card.Hash, server, strconv.Itoa(port))
		}
	}

	if !foundCard {
		dServer, dPort := getIP("discogssyncer", *host, portVal)
		rel, meta := getRelease("ListeningPile", dServer, strconv.Itoa(dPort))

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
		addTime := time.Unix(meta.DateAdded, 0)
		if time.Now().Sub(addTime).Hours() < 24*30*3 {
			card.Action = pbc.Card_DISMISS
		}
		cards.Cards = append(cards.Cards, &card)
		if !*dryRun {
			_, err = client.AddCards(context.Background(), &cards)
			if err != nil {
				log.Printf("Problem adding cards %v", err)
			}
		}
	}
}
