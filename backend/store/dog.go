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
}

type LocationTag struct {
	LocationID string
	FloorID    string
}

type dogStore struct {
	dogs []*dog.Dog
	locations []string
	floors map[string]int
	db map[LocationTag][]*dog.Dog
}

func NewDogStore(numDayCares int, numDogs int) DogStore {
	dogStore := &dogStore{
		dogs: make([]*dog.Dog, numDogs),
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
		dogStore.dogs[i] = &dog.Dog{
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
			LocationID: dogStore.dogs[i].LocationId,
			FloorID:    dogStore.dogs[i].FloorId,
		}
		dogStore.db[locationTag] = append(dogStore.db[locationTag], dogStore.dogs[i])
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