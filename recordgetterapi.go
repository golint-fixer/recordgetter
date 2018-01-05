package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.GetRecordRequest) (*pbrc.Record, error) {
	t := time.Now()
	if s.state.CurrentPick != nil {

		s.Log(fmt.Sprintf("Doing refresh: %v", in.GetRefresh()))
		if in.GetRefresh() {
			rec, err := s.getRelease(ctx, s.state.CurrentPick.Release.InstanceId)
			s.Log(fmt.Sprintf("GOT %v and %v", rec, err))
			if err != nil && len(rec.GetRecords()) == 1 {
				s.state.CurrentPick = rec.GetRecords()[0]
			}
		}

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
	s.update(in)
	s.state.CurrentPick = nil
	return s.GetRecord(ctx, &pb.GetRecordRequest{})
}

//Force forces a repick
func (s *Server) Force(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	s.state.CurrentPick = nil
	s.saveState()
	return &pb.Empty{}, nil
}
