package runtime

import (
	"fmt"
	"log"
)

func Run(config string) error {
	err := ParseConf(config)
	log.Println("Inserting subscribers into database")
	if err != nil {
		return err
	}

	var plmnID string = fmt.Sprintf("%s%s", PopulateConfig.MCC, PopulateConfig.MNC)

	for _, imsi := range PopulateConfig.IMSI {
		smData := generateSubs(imsi, plmnID, PopulateConfig.Slices)
		InsertSubscriber(imsi, plmnID, *smData)
	}
	return nil
}
