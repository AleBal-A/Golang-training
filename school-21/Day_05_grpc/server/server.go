package main

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "transmitter/transmitter"
)

type server struct {
	pb.UnimplementedTransmitterServiceServer
}

func (s *server) Transmit(req *pb.TransmitRequest, stream pb.TransmitterService_TransmitServer) error {
	sessionID := uuid.New().String()
	mean := rand.Float64()*20 - 10
	stddev := rand.Float64()*1.2 + 0.3
	log.Printf("New session: %s, Mean: %f, STD: %f, ID_From_Req: %s", sessionID, mean, stddev, req.ClientId)

	for {
		frequency := rand.NormFloat64()*stddev + mean
		timestamp := timestamppb.Now()

		res := &pb.TransmitResponse{
			SessionId:    sessionID,
			Frequency:    frequency,
			TimestampUTC: timestamp,
		}

		if err := stream.Send(res); err != nil {
			return err
		}

		time.Sleep(1000 * time.Millisecond)
	}

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTransmitterServiceServer(s, &server{})

	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve :%v", err)
	}
}
