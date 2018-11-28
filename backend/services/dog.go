package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/marcsj/grpc-sending-methods/backend/dog"
	"github.com/segmentio/ksuid"
)

type DogTrackServer struct{}

func (s DogTrackServer) TrackDogs(req *dog.TrackRequest, outStream dog.DogTrack_TrackDogsServer) error {
	locationTag := LocationTag{
		locationID: req.GetLocationId(),
		floorID:    req.GetFloorId(),
	}
	if len(daycareDogs[locationTag]) < 1 {
		return errors.New("no dogs found at that location")
	}

	for i := range daycareDogs[locationTag] {
		err := outStream.Send(daycareDogs[locationTag][i])
		if err != nil {
			return err
		}
	}

	dogChannel := make(chan *dog.Dog)
	for i := 0; i < 2; i++ {
		go DogCare(dogChannel, req.GetLocationId(), req.GetFloorId())
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

/*
	Ordinarily, this file would not contain any code apart from the server.
	In this case since we are showing our data without a real store(apart from in memory),
	we are using the same file. This example is just to show how a streaming setup might work.
*/
func InitializeDogs() {
	for i := 0; i < NUM_DAYCARES; i++ {
		location := ksuid.New().String()
		Locations[i] = location
		Floors[location] = randomdata.Number(1, 4)
		log.Printf("Location ID: %v\nFloors: %v", location, Floors[location])
	}
	for i := 0; i < 1000; i++ {
		locationID := Locations[randomdata.Number(0, NUM_DAYCARES)]
		Dogs[i] = &dog.Dog{
			Id:         ksuid.New().String(),
			Name:       randomdata.SillyName(),
			OwnerId:    ksuid.New().String(),
			LocationId: locationID,
			FloorId:    fmt.Sprint(randomdata.Number(Floors[locationID])),
			Location: &dog.Location{
				X: randomdata.Decimal(-20, 20),
				Y: randomdata.Decimal(-20, 20)},
			Status: dog.DogStatus(randomdata.Number(len(dog.DogStatus_value))),
		}
		locationTag := LocationTag{
			locationID: Dogs[i].LocationId,
			floorID:    Dogs[i].FloorId,
		}
		daycareDogs[locationTag] = append(daycareDogs[locationTag], Dogs[i])
	}
	log.Println("Finished initializing good puppers(and some not good).")
}

const NUM_DAYCARES = 5

var Dogs = make([]*dog.Dog, 1000)
var Locations = make([]string, NUM_DAYCARES)
var Floors = make(map[string]int, 0)

var daycareDogs = make(map[LocationTag][]*dog.Dog, 0)

type LocationTag struct {
	locationID string
	floorID    string
}

func DogCare(dogChannel chan *dog.Dog, locationID string, floorID string) {
	locationTag := LocationTag{
		locationID: locationID,
		floorID:    floorID,
	}
	dogs := daycareDogs[locationTag]
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			for _, d := range dogs {
				change := randomdata.Boolean()
				if change {
					d.Location = &dog.Location{
						X: d.Location.X + randomdata.Decimal(-1, 1),
						Y: d.Location.Y + randomdata.Decimal(-1, 1),
					}
					statusChange := randomdata.Boolean()
					if statusChange {
						d.Status = dog.DogStatus(randomdata.Number(len(dog.DogStatus_value)))
					}
					dogChannel <- d
				}
			}
		}
	}
}
