package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var updateRate = make(chan int)
var stopRequests = make(chan struct{})

func doEvery(initialRate int) {
	log.Printf("Starting loop with rate %d\n", initialRate)
	ticker := time.NewTicker(time.Duration(initialRate) * time.Millisecond)

	for {
		select {
		case <- stopRequests:
			return
		case tick:= <-ticker.C:
			sendRequests(tick)
		case rate := <-updateRate:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(rate) * time.Millisecond)
			log.Printf("Ticker reconfigured to %d ms\n", rate)
		}
	}
}

func sendRequests(t time.Time) {
	fmt.Printf("%v: Making request to hserver!\n", t)
	go http.Get("http://hserver")


	fmt.Printf("%v: Making request to vserver!\n", t)
	go http.Get("http://vserver")
}


func rate(w http.ResponseWriter, req *http.Request) {
	newRate := req.URL.Query().Get("rate")
	if newRate != "" {
		i, err := strconv.Atoi(newRate)
		if err == nil {
			updateRate <- i
			fmt.Fprintf(w, "Request rate updated to %d\n", i)
		}
	}
}

func start(w http.ResponseWriter, req *http.Request) {
	newRate := req.URL.Query().Get("rate")
	if newRate != "" {
		i, err := strconv.Atoi(newRate)
		if err == nil {
			go doEvery(i)
			fmt.Fprintf(w, "Started with rate %d\n", i)
		}
	} else {
		go doEvery(1000)
		fmt.Fprintf(w, "Started with default rate of 1 second\n")
	}
}

func stop(w http.ResponseWriter, req *http.Request) {
	stopRequests <- struct{}{}
	fmt.Fprintf(w, "Stopped requests\n")
}


func main() {
	http.HandleFunc("/rate", rate)
	http.HandleFunc("/start", start)
	http.HandleFunc("/stop", stop)

	http.ListenAndServe(":8090", nil)

}