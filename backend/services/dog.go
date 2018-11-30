package services

import (
	"context"
	"github.com/marcsj/grpc-sending-methods/backend/store"

	"github.com/marcsj/grpc-sending-methods/backend/dog"
)

type dogTrackServer struct{
	store store.DogStore
}

type DogTrackServer interface {
	TrackDogs(req *dog.TrackRequest, outStream dog.DogTrack_TrackDogsServer) error
	AddDog(context.Context, *dog.Dog) (*dog.Response, error)
	GetAllDogs(context.Context, *dog.TrackRequest) (*dog.ListDogs, error)
	DeleteDog(context.Context, *dog.DeleteRequest) (*dog.Response, error)
}

func NewDogTrackServer(store store.DogStore) DogTrackServer {
	return dogTrackServer{store: store}
}

func (s dogTrackServer) TrackDogs(req *dog.TrackRequest, outStream dog.DogTrack_TrackDogsServer) error {
	locationTag := store.LocationTag{
		LocationID: req.GetLocationId(),
		FloorID:    req.GetFloorId(),
	}

	dogChannel, err := s.store.GetDogStream(locationTag, 2)
	if err != nil {
		return err
	}

	for {
		select {
		case <-dogChannel:
			err := outStream.Send(<-dogChannel)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s dogTrackServer) AddDog(ctx context.Context, d *dog.Dog) (*dog.Response, error) {
	err := s.store.AddDog(d)
	if err != nil {
		return nil, err
	}
	return &dog.Response{Message: "Dog successfully added."}, nil
}

func (s dogTrackServer) GetAllDogs(ctx context.Context, req *dog.TrackRequest) (*dog.ListDogs, error) {
	locationTag := store.LocationTag{
		LocationID: req.GetLocationId(),
		FloorID: req.GetFloorId(),
	}
	dogs := s.store.GetAllDogs(locationTag)
	return &dog.ListDogs{Dogs: dogs}, nil
}

func (s dogTrackServer) DeleteDog(ctx context.Context, req *dog.DeleteRequest) (*dog.Response, error) {
	locationTag := store.LocationTag{
		LocationID: req.GetLocationId(),
		FloorID: req.GetFloorId(),
	}
	err := s.store.DeleteDog(locationTag, req.GetDogId())
	if err != nil {
		return nil, err
	}
	return &dog.Response{Message: "Dog successfully deleted."}, nil
}
