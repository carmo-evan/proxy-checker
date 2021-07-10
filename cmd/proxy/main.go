package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	enc := csvutil.NewEncoder(writer)

	var wg sync.WaitGroup
	numIPs := 100
	randomIPs := genRandomAddresses(numIPs)
	randomIPs = append(randomIPs, "", "") // working proxies
	resultsChan := make(chan Result)
	for _, ip := range randomIPs {
		wg.Add(1)
		logger := log.With().Str("address", ip).Logger()
		go CheckProxySOCKS(logger, ip+":1080", resultsChan, &wg)
	}
	go func() {
		for result := range resultsChan {
			if err := enc.Encode(result); err != nil {
				log.Fatal().Err(err).Send()
			}
		}
	}()
	wg.Wait()

}

func genRandomAddresses(len int) []string {
	var a []string
	for i := 0; i < len; i++ {
		ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
		a = append(a, ip)
	}
	return a
}
