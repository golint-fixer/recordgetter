package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.Empty) (*pbd.Release, error) {
	t := time.Now()
	if s.state.CurrentPick != nil {
		s.LogFunction("GetRecord-cache", t)
		return s.state.CurrentPick, nil
	}

	rel, _ := s.getReleaseFromPile("ListeningPile")
	if rel == nil {
		rel, _ = s.getReleaseFromCollection(true)
	}

	s.state.CurrentPick = rel
	s.saveState()

	s.LogFunction("GetRecord", t)
	return rel, nil
}

//Listened marks a record as Listened
func (s *Server) Listened(ctx context.Context, in *pbd.Release) (*pbd.Release, error) {
	rel, err := s.getRelease(ctx, in.Id)
	s.Log(fmt.Sprintf("Marking as listened: %v", rel))
	if err != nil {
		return nil, err
	}

	if in.Rating != rel.Rating {
		s.saveRelease(ctx, in)
	}

	s.moveReleaseToListeningBox(ctx, in)

	return s.GetRecord(ctx, &pb.Empty{})
}

//Force forces a repick
func (s *Server) Force(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	s.state.CurrentPick = nil
	s.saveState()
	return &pb.Empty{}, nil
}
