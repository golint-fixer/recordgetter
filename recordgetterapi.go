package main

import (
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.Empty) (*pbd.Release, error) {
	rel, _ := s.getReleaseFromPile("ListeningPile")
	if rel == nil {
		rel, _ = s.getReleaseFromCollection(true)
	}
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
