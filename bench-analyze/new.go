package main

import (
	"fmt"
	"log"

	"github.com/gyuho/dataframe"
)

func main() {
	fr, err := dataframe.NewFromCSV(nil, "testdata/bench-01-consul-1-monitor.csv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fr)
}

func aggregate(fpaths ...string) (dataframe.Frame, error) {

}
