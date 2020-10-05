package grpc

import (
	"time"

	grpcTypes "github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func ParseTimestampProto(ts string) (*timestamp.Timestamp, error) {
	date, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return &timestamp.Timestamp{}, err
	}

	return grpcTypes.TimestampProto(date)
}

func ParseDateProto(d string) (*timestamp.Timestamp, error) {
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		return &timestamp.Timestamp{}, err
	}

	return grpcTypes.TimestampProto(date)
}
