package runtime

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Run(config string) error {
	if err := ParseConf(config); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(PopulateConfig.Mongo.URL))
	defer cancel()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("Inserting subscribers into database")

	var plmnID string = fmt.Sprintf("%s%s", PopulateConfig.MCC, PopulateConfig.MNC)

	for _, imsi := range PopulateConfig.IMSI {
		smData := generateSubs(imsi, plmnID, PopulateConfig.Slices)
		if err := InsertSubscriber(client, PopulateConfig.Mongo.Name, imsi, plmnID, *smData); err != nil {
			return err
		}
	}
	return nil
}
