package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/proxy"
)

const (
	timeout = time.Duration(60 * time.Second)
)

//Result contain info about proxy
type Result struct {
	Addr string `csv:"address"`
	Res  bool   `csv:"result"`
}

//CheckProxySOCKS Check proxies on valid
func CheckProxySOCKS(logger zerolog.Logger, addr string, c chan Result, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	d := net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}

	//Sending request through proxy
	dialer, _ := proxy.SOCKS5("tcp", addr, nil, &d)

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Dial:              dialer.Dial,
		},
	}
	logger.Debug().Msg("starting request")
	response, err := httpClient.Get("https://en075sydjf92n6.x.pipedream.net")
	log.Println("finished request")
	if err != nil {
		logger.Err(err).Send()
		c <- Result{Addr: addr, Res: false}
		return
	}
	logger.Debug().Msg("success!")

	defer response.Body.Close()
	io.Copy(ioutil.Discard, response.Body)

	c <- Result{Addr: addr, Res: true}

	return nil

}
