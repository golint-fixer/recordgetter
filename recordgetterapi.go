package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordgetter/proto"
)

//GetRecord gets a record
func (s *Server) GetRecord(ctx context.Context, in *pb.GetRecordRequest) (*pb.GetRecordResponse, error) {
	t := time.Now()
	if s.state.CurrentPick != nil {
		if in.GetRefresh() {
			rec, err := s.rGetter.getRelease(ctx, s.state.CurrentPick.Release.InstanceId)
			if err == nil && len(rec.GetRecords()) == 1 {
				s.state.CurrentPick = rec.GetRecords()[0]
			}
			s.LogMilestone("GetRecord", "Refresh", t)
		}
		s.LogFunction("GetRecord", t)
		disk := int32(1)
		for _, score := range s.state.Scores {
			if score.InstanceId == s.state.CurrentPick.GetRelease().InstanceId {
				if score.DiskNumber >= disk {
					disk = score.DiskNumber + 1
				}
			}
		}

		return &pb.GetRecordResponse{Record: s.state.CurrentPick, NumListens: getNumListens(s.state.CurrentPick), Disk: disk}, nil
	}

	rec, err := s.getReleaseFromPile(t)
	s.LogMilestone("GetRecord", "GetRelease", t)
	if err != nil {
		return nil, err
	}

	disk := int32(1)
	if s.state.Scores != nil {
		for _, score := range s.state.Scores {
			s.Log(fmt.Sprintf("%v and %v", score, rec))
			time.Sleep(time.Second * 10)
			if score.InstanceId == rec.GetRelease().InstanceId {
				if score.DiskNumber >= disk {
					disk = score.DiskNumber + 1
				}
			}
		}
	}

	s.state.CurrentPick = rec
	s.saveState()

	s.LogFunction("GetRecord", t)
	return &pb.GetRecordResponse{Record: rec, NumListens: getNumListens(rec), Disk: disk}, nil
}

//Listened marks a record as Listened
func (s *Server) Listened(ctx context.Context, in *pbrc.Record) (*pb.Empty, error) {
	score := s.getScore(in)
	if score >= 0 {
		in.Release.Rating = score
		s.updater.update(ctx, in)
	}
	s.state.CurrentPick = nil
	s.saveState()
	return &pb.Empty{}, nil
}

//Force forces a repick
func (s *Server) Force(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	s.state.CurrentPick = nil
	s.saveState()
	return &pb.Empty{}, nil
}
