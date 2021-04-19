package runtime

import "fmt"

func Run(config string) error {
	err := ParseConf(config)
	if err != nil {
		return err
	}

	for _, i := range PopulateConfig.Slices {
		fmt.Println(i)
	}
	for _, i := range PopulateConfig.IMSI {
		fmt.Println(i)
	}
	return nil
}
