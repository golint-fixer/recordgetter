package main

import (
	"time"

	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func getNumListens(rc *pbrc.Record) int32 {
	if rc.GetMetadata().GetDateAdded() > (time.Now().AddDate(0, -3, 0).Unix()) {
		return 1
	}
	return 3
}
