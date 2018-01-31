package main

import (
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func getNumListens(rc *pbrc.Record) int32 {
	if rc.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_PRE_FRESHMAN {
		return 3
	}
	return 1
}
