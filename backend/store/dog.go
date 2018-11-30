package store

import (
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/segmentio/ksuid"
	"log"
	"math/rand"
	"time"
	"github.com/marcsj/grpc-sending-methods/backend/dog"
)
type DogStore interface {
	GetDogStream(tag LocationTag, numDogs int) (chan *dog.Dog, error)
	AddDog(dog *dog.Dog) error
	GetAllDogs(tag LocationTag) []*dog.Dog
	DeleteDog(tag LocationTag, dogID string) error
}

type LocationTag struct {
	LocationID string
	FloorID    string
}

type dogStore struct {
	locations []string
	floors map[string]int
	db map[LocationTag][]*dog.Dog
}

func NewDogStore(numDayCares int, numDogs int) DogStore {
	dogStore := &dogStore{
		locations: make([]string, numDayCares),
		floors: make(map[string]int),
		db: make(map[LocationTag][]*dog.Dog),
	}

	for i := 0; i < numDayCares; i++ {
		location := ksuid.New().String()
		dogStore.locations[i] = location
		dogStore.floors[location] = randomdata.Number(1, 4)
		log.Printf("Location ID: %v\nFloors: %v", location, dogStore.floors[location])
	}

	for i := 0; i < 1000; i++ {
		locationID := dogStore.locations[rand.Intn(numDayCares)]
		newDog := &dog.Dog{
			Id:         ksuid.New().String(),
			Name:       randomdata.SillyName(),
			OwnerId:    ksuid.New().String(),
			LocationId: locationID,
			FloorId:    fmt.Sprint(rand.Intn(dogStore.floors[locationID])+1),
			Location: &dog.Location{
				X: randomdata.Decimal(-20, 20),
				Y: randomdata.Decimal(-20, 20)},
			Status: dog.DogStatus(randomdata.Number(len(dog.DogStatus_value))),
		}
		locationTag := LocationTag{
			LocationID: newDog.LocationId,
			FloorID:    newDog.FloorId,
		}
		dogStore.db[locationTag] = append(dogStore.db[locationTag], newDog)
	}
	log.Println("Finished initializing good puppers(and some not good).")
	// wait to keep logs clean for location IDs
	time.Sleep(500 * time.Millisecond)
	return dogStore
}


func (s dogStore) GetDogStream(tag LocationTag, numDogs int) (chan *dog.Dog, error) {
	if len(s.db[tag]) < 1 {
		return nil, errors.New("no dogs found at that location")
	}

	dogChannel := make(chan *dog.Dog, numDogs)

	for i := 0; i < numDogs; i++ {
		go s.dogCare(dogChannel, tag)
	}
	return dogChannel, nil
}

func (s dogStore) AddDog(dog *dog.Dog) error {
	locationTag := LocationTag{
		LocationID: dog.GetLocationId(),
		FloorID: dog.GetFloorId(),
	}
	if locationTag.FloorID == "" || locationTag.LocationID == "" {
		return errors.New("new dog location invalid")
	}
	s.db[locationTag] = append(s.db[locationTag], dog)
	return nil
}

func (s dogStore) GetAllDogs(tag LocationTag) []*dog.Dog {
	return s.db[tag]
}

func (s dogStore) DeleteDog(tag LocationTag, dogID string) error {
	for i := range s.db[tag] {
		if s.db[tag][i].Id == dogID {
			s.db[tag] = append(s.db[tag][:i], s.db[tag][i+1:]...)
			return nil
		}
	}
	return errors.New("dog not found")
}

func (s dogStore) dogCare(dogChannel chan *dog.Dog, locationTag LocationTag) {
	dogs := s.db[locationTag]
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