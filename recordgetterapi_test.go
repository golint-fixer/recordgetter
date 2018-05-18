package main

import (
	"context"
	"testing"
	"time"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordgetter/proto"
)

type testGetter struct {
	records []*pbrc.Record
}

func (tg *testGetter) getRecords() (*pbrc.GetRecordsResponse, error) {
	return &pbrc.GetRecordsResponse{Records: tg.records}, nil
}

func TestRecordGetDiskReturn(t *testing.T) {
	s := InitTestServer()
	s.rGetter = &testGetter{records: []*pbrc.Record{
		&pbrc.Record{Release: &pbgd.Release{InstanceId: 1234, FormatQuantity: 2}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PRE_PROFESSOR, DateAdded: 1234}},
	}}

	resp, err := s.GetRecord(context.Background(), &pb.GetRecordRequest{})

	if err != nil {
		t.Fatalf("Error getting record: %v", err)
	}

	if resp.Disk != 1 {
		t.Errorf("Disk was not reported: %v", resp)
	}
}

func TestRecordGetDiskSkipOnDate(t *testing.T) {
	s := InitTestServer()
	s.rGetter = &testGetter{records: []*pbrc.Record{
		&pbrc.Record{Release: &pbgd.Release{InstanceId: 12, FormatQuantity: 2}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PRE_FRESHMAN, DateAdded: 12}},
		&pbrc.Record{Release: &pbgd.Release{InstanceId: 1234, FormatQuantity: 2}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PRE_PROFESSOR, DateAdded: 1234}},
	}}
	s.state.Scores = append(s.state.Scores, &pb.DiskScore{InstanceId: 1234, DiskNumber: 1, ScoreDate: time.Now().Unix(), Score: 5})

	resp, err := s.GetRecord(context.Background(), &pb.GetRecordRequest{})
	if err != nil {
		t.Fatalf("Error getting record: %v", err)
	}

	if resp.Record.GetRelease().InstanceId != 12 {
		t.Errorf("Wrong record returned: %v", resp)
	}
}

func TestRecordGetNextDisk(t *testing.T) {
	s := InitTestServer()
	s.rGetter = &testGetter{records: []*pbrc.Record{
		&pbrc.Record{Release: &pbgd.Release{InstanceId: 1234, FormatQuantity: 2}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PRE_FRESHMAN, DateAdded: 12}},
		&pbrc.Record{Release: &pbgd.Release{InstanceId: 1234, FormatQuantity: 2}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PRE_PROFESSOR, DateAdded: 1234}},
	}}

	s.state.Scores = append(s.state.Scores, &pb.DiskScore{InstanceId: 1234, DiskNumber: 1, ScoreDate: time.Now().AddDate(0, -1, 0).Unix(), Score: 5})

	resp, err := s.GetRecord(context.Background(), &pb.GetRecordRequest{})
	if err != nil {
		t.Fatalf("Error getting record: %v", err)
	}

	if resp.Record.GetRelease().InstanceId != 1234 {
		t.Errorf("Wrong record returned: %v", resp)
	}

	if resp.Disk != 2 {
		t.Errorf("Wrong disk number returned %v", resp)
	}
}
