package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbc "github.com/brotherlogic/cardserver/card"
	pb "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pbrg "github.com/brotherlogic/recordgetter/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	serving    bool
	delivering bool
	state      *pbrg.State
}

const (
	wait = 5 * time.Second

	//KEY under which we store the collection
	KEY = "/github.com/brotherlogic/recordgetter/state"
)

func (s *Server) getRelease(ctx context.Context, instance int32) (*pbrc.GetRecordsResponse, error) {
	host, port := s.GetIP("recordcollection")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pbrc.NewRecordCollectionServiceClient(conn)

	return client.GetRecords(ctx, &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbd.Release{InstanceId: instance}}})
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

func (s *Server) update(rec *pbrc.Record) error {
	host, port := s.GetIP("recordcollection")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pbrc.NewRecordCollectionServiceClient(conn)
	_, err = client.UpdateRecord(context.Background(), &pbrc.UpdateRecordRequest{Update: rec})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) getReleaseFromPile(t time.Time) (*pbrc.Record, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	host, port := s.GetIP("recordcollection")
	conn, err := grpc.Dial(host+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := pbrc.NewRecordCollectionServiceClient(conn)

	//Only get clean records
	s.LogMilestone("GetRecord", "PreGetRecords", t)
	r, err := client.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Metadata: &pbrc.ReleaseMetadata{Dirty: false}, Release: &pbd.Release{FolderId: 812802}}}, grpc.MaxCallRecvMsgSize(1024*1024*1024))
	if err != nil {
		return nil, err
	}
	s.LogMilestone("GetRecord", "PostGetRecords", t)

	if len(r.GetRecords()) == 0 {
		return nil, nil
	}

	var newRec *pbrc.Record
	newRec = nil
	pDate := int64(0)
	for _, rc := range r.GetRecords() {
		if rc.GetMetadata().DateAdded > pDate && rc.GetRelease().Rating == 0 && !rc.GetMetadata().GetDirty() {
			pDate = rc.GetMetadata().DateAdded
			newRec = rc
		}
	}
	s.LogMilestone("GetRecord", "RecordSelected", t)

	return newRec, nil
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

	retRel := r.Records[rand.Intn(len(r.Records))].GetRelease()
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

	for _, release := range r.Records {
		if int(release.GetRelease().Id) == id {
			return release.GetRelease()
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
	if cPort > 0 {
		conn, _ := grpc.Dial(cServer+":"+strconv.Itoa(cPort), grpc.WithInsecure())
		defer conn.Close()
		client := pbc.NewCardServiceClient(conn)

		cardList, err := client.GetCards(context.Background(), &pbc.Empty{})

		if err == nil {
			for _, card := range cardList.Cards {
				if card.Hash == "discogs" {
					return true
				}
			}
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

	card := pbc.Card{Text: pbd.GetReleaseArtist(rel) + " - " + rel.Title, Hash: "discogs", Image: imageURL, Priority: 100}
	return card
}

//Init a record getter
func Init() *Server {
	s := &Server{GoServer: &goserver.GoServer{}, serving: true, delivering: true, state: &pbrg.State{}}
	s.Register = s
	s.PrepServer()
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pbrg.RegisterRecordGetterServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(master bool) error {
	s.delivering = master

	if master {
		return s.readState()
	}

	return nil
}

// GetState gets the state of the server
func (s Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

// This is the only method that interacts with disk
func (s *Server) readState() error {
	state := &pbrg.State{}
	data, _, err := s.KSclient.Read(KEY, state)

	if err != nil {
		return err
	}

	if data != nil {
		s.state = data.(*pbrg.State)
	}

	return nil
}

func (s *Server) saveState() {
	s.KSclient.Save(KEY, s.state)
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

	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	server.RegisterServer("recordgetter", false)
	//server.RegisterServingTask(server.GetRecords)
	server.Log("Starting!")
	server.Serve()
}
