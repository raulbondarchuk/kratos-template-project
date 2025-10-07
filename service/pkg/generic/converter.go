package generic

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertToGoogleTimestamp converts a time.Time to a timestamppb.Timestamp
// t - time.Time
// returns *timestamppb.Timestamp
func ConvertToGoogleTimestamp(t time.Time) *timestamppb.Timestamp {
	if !t.IsZero() {
		t = t.UTC()
	}
	return timestamppb.New(t)
}

// ConvertToTime converts a timestamppb.Timestamp to a time.Time
// t - *timestamppb.Timestamp
// returns time.Time
func ConvertToTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}
