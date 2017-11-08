package main

import (
	"log"
	"time"

	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.Empty) (*pbd.Release, error) {
	t := time.Now()
	log.Printf("HERE %v and %v", s, s.state)
	if s.state.CurrentPick != nil {
		s.LogFunction("GetRecord-cache", t)
		log.Printf("HUH")
		return s.state.CurrentPick, nil
	}

	log.Printf("GETTING RELEASE")
	rel, _ := s.getReleaseFromPile("ListeningPile")
	log.Printf("GOT1 %v", rel)
	if rel == nil {
		rel, _ = s.getReleaseFromCollection(true)
		log.Printf("NOW %v", rel)
	}

	log.Printf("GOT %v", rel)

	s.state.CurrentPick = rel
	s.saveState()

	log.Printf("NOW %v and %v", s, s.state)
	s.LogFunction("GetRecord", t)
	return rel, nil
}

//Listened marks a record as Listened
func (s *Server) Listened(ctx context.Context, in *pbd.Release) (*pbd.Release, error) {
	rel, err := s.getRelease(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if in.Rating != rel.Rating {
		s.saveRelease(ctx, in)
	}

	s.moveReleaseToListeningBox(ctx, in)

	return s.GetRecord(ctx, &pb.Empty{})
}
