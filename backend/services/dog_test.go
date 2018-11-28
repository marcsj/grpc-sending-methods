package services

import (
	"context"
	"github.com/marcsj/streaming-grpc-web-example/backend/dog"
	"golang.org/x/net/websocket"
	"google.golang.org/grpc"
	"io"
	"log"
	"testing"
	"time"
)

func TestGRPCServer_TrackDogs(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	ctx := context.Background()
	clientDeadline := time.Now().Add(time.Duration(1) * time.Minute)
	ctx, _ = context.WithDeadline(ctx, clientDeadline)

	client := dog.NewDogTrackClient(conn)
	request := &dog.TrackRequest{
		LocationId: "<location_id>",
		FloorId: "1",
	}
	stream, err := client.TrackDogs(ctx, request)
	if err != nil {
		t.Error(err)
	}

	for {
		dog, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}
		log.Println(
			dog.Id,
			": Name: ", dog.Name,
			"Owned by: ", dog.OwnerId,
			"Status: ", dog.Status.String(),
			dog.Location.X, dog.Location.Y)
	}
}

func TestWSServer_TrackDogs(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var err error
	conn, err := websocket.Dial("ws://localhost:8081/v1/dogs/track", "", "http://localhost/")
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte("test message"))
	connectionBytes := make([]byte, 256)
	_, err = conn.Read(connectionBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(connectionBytes))
	err = conn.Close()
	if err != nil {
		t.Error(err)
	}
}