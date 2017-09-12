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
