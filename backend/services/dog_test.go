package services

import (
	"context"
	"github.com/marcsj/streaming-grpc-web-example/backend/dog"
	"google.golang.org/grpc"
	"io"
	"log"
	"testing"
	"time"
)

func TestDogTrackServer_TrackDogs(t *testing.T) {
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
