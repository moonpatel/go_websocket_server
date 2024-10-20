package main

import (
	"flag"
	"log"
	"net/http"
	"runtime"
	"time"
)

var addr = flag.String("addr", ":8010", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func printRuntimeState() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("Sys = %v MiB", bToMb(m.Sys))
	log.Printf("NumGC = %v\n", m.NumGC)
	// Print the number of goroutines
	log.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())

	// Print the number of logical CPUs
	log.Printf("Number of CPUs: %d\n", runtime.NumCPU())

	// Print the number of cgo calls
	log.Printf("Number of cgo calls: %d\n", runtime.NumCgoCall())
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func startRuntimeStats(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				printRuntimeState()
			}
		}
	}()
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
	})

	// go startRuntimeStats(5)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
