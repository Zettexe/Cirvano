package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	pb "github.com/zettexe/cirvano/proto/gen/broadcast"
)

var playlistData pb.Update
var playlistDataUpdated = false

type BroadcastServer struct {
	pb.UnimplementedBroadcasterServer
	clients map[string]pb.Broadcaster_BroadcastServer
}

func newServer() *BroadcastServer {
	return &BroadcastServer{
		clients: make(map[string]pb.Broadcaster_BroadcastServer),
	}
}

func (s *BroadcastServer) Broadcast(stream pb.Broadcaster_BroadcastServer) error {
	clientId := "ttest" // TODO: generate a uuid
	s.clients[clientId] = stream

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			delete(s.clients, clientId)
			return nil
		}
		if err != nil {
			delete(s.clients, clientId)
			return err
		}
	}
}

type Players struct {
	control *beep.Ctrl
	volume  *effects.Volume
}

var players Players

func (s *BroadcastServer) PlaySong(ctx context.Context, message *pb.PlayRequest) (*pb.PlayResponse, error) {
	go func() {
		f, err := os.Open(message.Filename)
		if err != nil {
			log.Fatal(err)
		}

		streamer, format, err := mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		players.control = &beep.Ctrl{
			Streamer: beep.Loop(-1, streamer),
			Paused:   false,
		}

		players.volume = &effects.Volume{
			Streamer: players.control,
			Base:     2,
			Volume:   -3,
			Silent:   false,
		}

		fmt.Println("Play Start")

		speaker.Play(players.volume)

		for {
			select {
			case <-time.After(time.Second):
				speaker.Lock()
				playlistData.Duration = format.SampleRate.D(streamer.Position()).String()
				playlistDataUpdated = true
				speaker.Unlock()
			}
		}
	}()

	return &pb.PlayResponse{}, nil
}

func (s *BroadcastServer) SongVolume(ctx context.Context, volumeChange *pb.VolumeChangeRequest) (*pb.VolumeChangeResponse, error) {
	fmt.Printf("Requested volume change by: %f\n", volumeChange.Volume)
	speaker.Lock()
	players.volume.Volume += float64(volumeChange.Volume)
	speaker.Unlock()
	return &pb.VolumeChangeResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bs := newServer()
	pb.RegisterBroadcasterServer(grpcServer, bs)

	go func() {
		for {
			if playlistDataUpdated == true {
				playlistDataUpdated = false
				for _, client := range bs.clients {
					client.Send(&playlistData)
				}
			}
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
