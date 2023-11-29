package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/zettexe/cirvano/proto/gen/broadcast"
)

func main() {
	for {

		conn, err := grpc.Dial("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewBroadcasterClient(conn)

		stream, err := client.Broadcast(context.Background())
		if err != nil {
			log.Fatalf("Error creating stream: %v", err)
		}

		if err := stream.Send(&pb.RegisterRequest{}); err != nil {
			log.Fatalf("Error sending request: %v", err)
		}

		go func() {
			for {
				update, err := stream.Recv()
				if err != nil {
					log.Fatalf("Error receiving update: %v", err)
				}
				log.Printf("Received update: duration: %s\n", update.GetDuration())
			}
		}()

		time.Sleep(time.Second * 10)

		log.Println("Sending play request")
		if _, err := client.PlaySong(stream.Context(), &pb.PlayRequest{Filename: "TOOL - The Pot.mp3"}); err != nil {
			log.Fatalf("Error sending play request: %v", err)
		}

		time.Sleep(time.Second * 10)

		log.Println("Sending volume change request")
		if _, err := client.SongVolume(stream.Context(), &pb.VolumeChangeRequest{Volume: -0.1}); err != nil {
			log.Fatalf("Error sending volume change request: %v", err)
		}

		time.Sleep(time.Second)

		log.Println("Sending volume change request")
		if _, err := client.SongVolume(stream.Context(), &pb.VolumeChangeRequest{Volume: -0.1}); err != nil {
			log.Fatalf("Error sending volume change request: %v", err)
		}

		time.Sleep(time.Second)

		log.Println("Sending volume change request")
		if _, err := client.SongVolume(stream.Context(), &pb.VolumeChangeRequest{Volume: -0.1}); err != nil {
			log.Fatalf("Error sending volume change request: %v", err)
		}

		time.Sleep(time.Second)

		log.Println("Sending volume change request")
		if _, err := client.SongVolume(stream.Context(), &pb.VolumeChangeRequest{Volume: -0.1}); err != nil {
			log.Fatalf("Error sending volume change request: %v", err)
		}

		time.Sleep(time.Second)

		log.Println("Sending volume change request")
		if _, err := client.SongVolume(stream.Context(), &pb.VolumeChangeRequest{Volume: -0.1}); err != nil {
			log.Fatalf("Error sending volume change request: %v", err)
		}

		for {

		}
	}
}
