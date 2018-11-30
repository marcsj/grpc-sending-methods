package services

import (
	"github.com/marcsj/grpc-sending-methods/backend/store"

	"github.com/marcsj/grpc-sending-methods/backend/dog"
)

type dogTrackServer struct{
	store store.DogStore
}

type DogTrackServer interface {
	TrackDogs(req *dog.TrackRequest, outStream dog.DogTrack_TrackDogsServer) error
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
