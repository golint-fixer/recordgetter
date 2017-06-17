package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbc "github.com/brotherlogic/cardserver/card"
	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	serving    bool
	delivering bool
}

const (
	wait = 5 * time.Second
)

func getIP(servername string) (string, int) {
	log.Printf("Looking for %v", servername)
	conn, _ := grpc.Dial("192.168.86.64:50055", grpc.WithInsecure())
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

func getReleaseFromPile(folderName string) (*pbd.Release, *pb.ReleaseMetadata) {
	log.Printf("Getting release from %v", folderName)
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := getIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Error dialling server: %v", err)
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatalf("Problem getting releases from Pile %v", err)
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
	log.Printf("Found %v", newRel)
	return newRel, meta
}

func getReleaseFromCollection(allowSeven bool) (*pbd.Release, *pb.ReleaseMetadata) {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := getIP("discogssyncer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
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
		log.Fatalf("Problem getting releases from collection %v", err)
	}

	log.Printf("Trying to get from %v: %v", len(r.Releases), r.Releases)
	retRel := r.Releases[rand.Intn(len(r.Releases))]
	meta, err := client.GetMetadata(context.Background(), retRel)
	if err != nil {
		log.Fatalf("Problem getting metadata %v", err)
	}
	return retRel, meta
}

func getReleaseWithID(folderName string, id int) *pbd.Release {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := getIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Cannto dial %v", err)
	}

	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, err := client.GetReleasesInFolder(context.Background(), folderList)
	if err != nil {
		log.Fatalf("Problem getting releases with a given ID %v", err)
	}

	for _, release := range r.Releases {
		if int(release.Id) == id {
			return release
		}
	}

	return nil
}

func deleteCard(hash string) {
	host, port := getIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	client.DeleteCards(context.Background(), &pbc.DeleteRequest{Hash: hash})
}

func scoreCard(releaseID int, rating int) bool {
	host, port := getIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	allowSeven := true
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	release := getReleaseWithID("ListeningPile", releaseID)
	if release == nil {
		release = getReleaseWithID("7s", releaseID)
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

func hasCurrentCard() bool {
	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver")
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	log.Printf("GETTING CARDS")
	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	if err != nil {
		log.Fatalf("Whoops2: %v", err)
	}

	for _, card := range cardList.Cards {
		log.Printf("CHECKING %v", card)
		if card.Hash == "discogs" {
			return true
		}
	}
	return false
}

func addCards(cardList *pbc.CardList) {
	cServer, cPort := getIP("cardserver")
	conn, err := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	log.Printf("CALLING ADD CARDS")
	client.AddCards(context.Background(), cardList)
	log.Printf("DONE")
}

func (s Server) processCard() bool {
	//Get the latest card from the cardserver
	cServer, cPort := getIP("cardserver")
	conn, _ := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	allowSeven := true

	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	if err != nil {
		panic(err)
	}

	for _, card := range cardList.Cards {
		if card.Hash == "discogs-process" {
			releaseID, _ := strconv.Atoi(card.Text)
			if card.ActionMetadata != nil {
				rating, _ := strconv.Atoi(card.ActionMetadata[0])
				if s.delivering {
					allowSeven = scoreCard(releaseID, rating)
				}
			} else {
				if s.delivering {
					allowSeven = scoreCard(releaseID, -1)
				}
			}
			if s.delivering {
				deleteCard(card.Hash)
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

// GetRecords runs the get records loop
func (s Server) GetRecords() {
	for s.serving {
		log.Printf("Sleepinging %v", wait)
		time.Sleep(wait)
		log.Printf("Running a single")
		s.runSingle()
	}
}

func (s Server) runSingle() {
	log.Printf("Logging is on!")

	foundCard := hasCurrentCard()
	allowSeven := s.processCard()
	cards := pbc.CardList{}

	log.Printf("CURRENT Card: %v", foundCard)

	if !foundCard {
		rel, meta := getReleaseFromPile("ListeningPile")

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
			rel, _ := getReleaseFromCollection(allowSeven)
			card := getCard(rel)
			card.Action = pbc.Card_DISMISS
			if rel.Rating <= 0 {
				card.Result = &pbc.Card{Hash: "discogs-process", Priority: -10, Text: strconv.Itoa(int(rel.Id))}
				card.Action = pbc.Card_RATE
			}
			cards.Cards = append(cards.Cards, &card)
		}
	}

	log.Printf("RUNNING SINGLE: %v", s.delivering)
	if s.delivering {
		addCards(&cards)
	}
	log.Printf("DONE")
}

//Init a record getter
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}, serving: true, delivering: true}
	s.Register = s
	return s
}

// DoRegister does RPC registration
func (s Server) DoRegister(server *grpc.Server) {
	// Noop
}

// ReportHealth alerts if we're not healthy
func (s Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s Server) Mote(master bool) error {
	s.delivering = master
	return nil
}

func main() {
	var quiet = flag.Bool("quiet", true, "Show all output")
	flag.Parse()

	server := Init()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	server.PrepServer()
	server.RegisterServer("recordgetter", false)
	server.RegisterServingTask(server.GetRecords)
	server.Serve()
}
