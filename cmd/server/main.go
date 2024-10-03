package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	"codeberg.org/shinyzero0/ostrovok2024-client/proto"
	"codeberg.org/shinyzero0/ostrovok2024-client/utils"
	"google.golang.org/grpc"
)

func main() {
	if err := f(); err != nil {
		log.Fatal(err)
	}
}

type server struct {
	proto.UnimplementedProcessorServer
}

func (s *server) ProcessBatch(_ context.Context, in *proto.BatchedInput) (*proto.BatchedOutput, error) {
	outs, _ := utils.Map(
		in.Records,
		func(i *proto.InputRecord) (o *proto.OutputRecord, err error) {
			return &proto.OutputRecord{
				Bathroom:       0,
				Class:          0,
				BedroomsAmount: 0,
				BeddingType:    0,
				HasBalcony:     false,
				IsClub:         false,
				View:           0,
				Floor:          0,
				Capacity:       0,
				Quality:        0,
			}, nil
		},
	)
	return &proto.BatchedOutput{
		Records: outs,
	}, nil
}

func f() error {
	srvu, ok := os.LookupEnv("SERVER_URI")
	if !ok {
		return errors.New("fuck")
	}
	lis, err := net.Listen("tcp", srvu)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterProcessorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}
