package main

import (
	"errors"
	"time"

	"golang.org/x/net/context"

	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.Empty) (*pbrc.Record, error) {
	t := time.Now()
	if s.state.CurrentPick != nil {
		s.LogFunction("GetRecord-cache", t)
		return s.state.CurrentPick, nil
	}

	rec, err := s.getReleaseFromPile()
	if err != nil {
		return nil, err
	}

	s.state.CurrentPick = rec
	s.saveState()

	s.LogFunction("GetRecord", t)
	return rec, nil
}

//Listened marks a record as Listened
func (s *Server) Listened(ctx context.Context, in *pbrc.Record) (*pbrc.Record, error) {
	return nil, errors.New("UNIMPLEMENTED")
}

//Force forces a repick
func (s *Server) Force(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	s.state.CurrentPick = nil
	s.saveState()
	return &pb.Empty{}, nil
}
