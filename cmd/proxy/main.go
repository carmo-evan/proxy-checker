package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/jszwec/csvutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	doControlRequest()

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	enc := csvutil.NewEncoder(writer)
	seenIPs := make(map[string]bool)
	tries := make([]struct{}, 100)

	for range tries {
		var wg sync.WaitGroup
		numIPs := 10
		randomIPs := genRandomAddresses(numIPs, seenIPs)
		randomIPs = append(randomIPs, "119.8.44.87") // working proxies
		resultsChan := make(chan Result)
		go func() {
			for result := range resultsChan {
				if err := enc.Encode(result); err != nil {
					log.Fatal().Err(err).Send()
				}
				writer.Flush()
			}
		}()
		for _, ip := range randomIPs {
			wg.Add(1)
			seenIPs[ip] = true
			logger := log.With().Str("address", ip).Logger()
			go CheckProxySOCKS(logger, ip+":1080", resultsChan, &wg)
		}
		wg.Wait()
	}
}

func genRandomAddresses(len int, seenIPs map[string]bool) []string {
	var a []string
	for i := 0; i < len; i++ {
		ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
		if ok := seenIPs[ip]; ok {
			i--
			continue
		}
		a = append(a, ip)
	}
	return a
}

func doControlRequest() {
	res, err := http.DefaultClient.Get("https://en075sydjf92n6.x.pipedream.net")
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("control request successful: %v", res.Status)
}
