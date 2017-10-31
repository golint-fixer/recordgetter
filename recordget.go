package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbc "github.com/brotherlogic/cardserver/card"
	pb "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pbrg "github.com/brotherlogic/recordgetter/proto"
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

func (s *Server) getRelease(ctx context.Context, id int32) (*pbd.Release, error) {
	host, port := s.GetIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)

	return client.GetSingleRelease(ctx, &pbd.Release{Id: id})
}

func (s *Server) saveRelease(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	host, port := s.GetIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)

	return client.UpdateRating(ctx, in)
}

func (s *Server) moveReleaseToListeningBox(ctx context.Context, in *pbd.Release) (*pb.Empty, error) {
	host, port := s.GetIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	return client.MoveToFolder(ctx, &pb.ReleaseMove{Release: in, NewFolderId: 673768})
}

func (s *Server) getReleaseFromPile(folderName string) (*pbd.Release, *pb.ReleaseMetadata) {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := s.GetIP("discogssyncer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, _ := client.GetReleasesInFolder(context.Background(), folderList)

	if len(r.Releases) == 0 {
		return nil, nil
	}

	var newRel *pbd.Release
	newRel = nil
	pDate := int64(math.MaxInt64)
	for _, rel := range r.Releases {
		meta, err2 := client.GetMetadata(context.Background(), rel)
		if err2 == nil {
			if meta.DateAdded > (time.Now().AddDate(0, -3, 0).Unix()) && meta.DateAdded < pDate {
				newRel = rel
				pDate = meta.DateAdded
			}
		}
	}

	if newRel == nil {
		newRel = r.Releases[rand.Intn(len(r.Releases))]
	}
	meta, _ := client.GetMetadata(context.Background(), newRel)
	return newRel, meta
}

func (s *Server) getReleaseFromCollection(allowSeven bool) (*pbd.Release, *pb.ReleaseMetadata) {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := s.GetIP("discogssyncer")
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
	r, _ := client.GetReleasesInFolder(context.Background(), folderList)

	retRel := r.Releases[rand.Intn(len(r.Releases))]
	meta, _ := client.GetMetadata(context.Background(), retRel)

	return retRel, meta
}

func (s *Server) getReleaseWithID(folderName string, id int) *pbd.Release {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := s.GetIP("discogssyncer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())

	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	folderList := &pb.FolderList{}
	folder := &pbd.Folder{Name: folderName}
	folderList.Folders = append(folderList.Folders, folder)
	r, _ := client.GetReleasesInFolder(context.Background(), folderList)

	for _, release := range r.Releases {
		if int(release.Id) == id {
			return release
		}
	}
	return nil
}

func (s *Server) deleteCard(hash string) {
	host, port := s.GetIP("cardserver")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	client.DeleteCards(context.Background(), &pbc.DeleteRequest{Hash: hash})
}

func (s *Server) scoreCard(releaseID int, rating int) bool {
	host, port := s.GetIP("discogssyncer")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	allowSeven := true
	defer conn.Close()
	client := pb.NewDiscogsServiceClient(conn)
	release := s.getReleaseWithID("ListeningPile", releaseID)
	if release == nil {
		release = s.getReleaseWithID("7s", releaseID)
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

func (s *Server) hasCurrentCard() bool {
	//Get the latest card from the cardserver
	cServer, cPort := s.GetIP("cardserver")
	conn, _ := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	cardList, _ := client.GetCards(context.Background(), &pbc.Empty{})

	for _, card := range cardList.Cards {
		if card.Hash == "discogs" {
			return true
		}
	}
	return false
}

func (s *Server) addCards(cardList *pbc.CardList) {
	cServer, cPort := s.GetIP("cardserver")
	conn, _ := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)
	client.AddCards(context.Background(), cardList)
}

func (s Server) processCard() (bool, error) {
	//Get the latest card from the cardserver
	cServer, cPort := s.GetIP("cardserver")
	conn, _ := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
	defer conn.Close()
	client := pbc.NewCardServiceClient(conn)

	allowSeven := true

	cardList, err := client.GetCards(context.Background(), &pbc.Empty{})
	if err != nil {
		return false, err
	}

	for _, card := range cardList.Cards {
		if card.Hash == "discogs-process" {
			releaseID, _ := strconv.Atoi(card.Text)
			if card.ActionMetadata != nil {
				rating, _ := strconv.Atoi(card.ActionMetadata[0])
				if s.delivering {
					allowSeven = s.scoreCard(releaseID, rating)
				}
			} else {
				if s.delivering {
					allowSeven = s.scoreCard(releaseID, -1)
				}
			}
			if s.delivering {
				s.deleteCard(card.Hash)
			}
		}
	}

	return allowSeven, nil
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
		time.Sleep(wait)
		s.runSingle()
	}
}

func (s Server) runSingle() {
	foundCard := s.hasCurrentCard()
	allowSeven, err := s.processCard()

	if err != nil {
		return
	}

	cards := pbc.CardList{}

	if !foundCard {
		rel, meta := s.getReleaseFromPile("ListeningPile")

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
			rel, _ := s.getReleaseFromCollection(allowSeven)
			card := getCard(rel)
			card.Action = pbc.Card_DISMISS
			if rel.Rating <= 0 {
				card.Result = &pbc.Card{Hash: "discogs-process", Priority: -10, Text: strconv.Itoa(int(rel.Id))}
				card.Action = pbc.Card_RATE
			}
			cards.Cards = append(cards.Cards, &card)
		}
	}

	if s.delivering {
		s.addCards(&cards)
	}
}

//Init a record getter
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}, serving: true, delivering: true}
	s.Register = s
	return s
}

// DoRegister does RPC registration
func (s Server) DoRegister(server *grpc.Server) {
	pbrg.RegisterRecordGetterServer(server, &s)
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
	var quiet = flag.Bool("quiet", false, "Show all output")
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
