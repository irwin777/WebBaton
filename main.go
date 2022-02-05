package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type batons = struct {
	sync.Mutex
	pin string
	lpt time.Time
}

const (
	pin1     = "2"
	pin2     = "65"
	pin3     = "67"
	pin4     = "66"
	pin5     = "21"
	pin6     = "18"
	pin7     = "19"
	pin8     = "7"
	pin9     = "8"
	pin10    = "200"
	pin11    = "9"
	pin12    = "10"
	pin13    = "201"
	pin14    = "20"
	pin15    = "198"
	pin16    = "199"
	pushtime = 300
)

var pinmap map[int]*batons

func initbat() {
	pinmap = make(map[int]*batons)
	for i := 1; i < 17; i++ {
		pinmap[i] = new(batons)
	}
	pinmap[1].pin = pin1
	pinmap[2].pin = pin2
	pinmap[3].pin = pin3
	pinmap[4].pin = pin4
	pinmap[5].pin = pin5
	pinmap[6].pin = pin6
	pinmap[7].pin = pin7
	pinmap[8].pin = pin8
	pinmap[9].pin = pin9
	pinmap[10].pin = pin10
	pinmap[11].pin = pin11
	pinmap[12].pin = pin12
	pinmap[13].pin = pin13
	pinmap[14].pin = pin14
	pinmap[15].pin = pin15
	pinmap[16].pin = pin16

	for _, value := range pinmap {
		ioutil.WriteFile("/sys/class/gpio/export", []byte(value.pin), 0644)
		ioutil.WriteFile("/sys/class/gpio/gpio"+value.pin+"/direction", []byte("out"), 0644)
		ioutil.WriteFile("/sys/class/gpio/gpio"+value.pin+"/value", []byte("1"), 0644)
		fmt.Println("Init pin", value.pin)
		time.Sleep(50 * time.Millisecond)
	}
}

func pusher(n int) {
	r, ok := pinmap[n]
	if !ok {
		return
	}
	tn := time.Now()
	r.Lock()
	timeDif := int(tn.Sub(r.lpt).Seconds())
	if timeDif < 1 {
		time.Sleep(1 * time.Second)
	}
	ioutil.WriteFile("/sys/class/gpio/gpio"+r.pin+"/value", []byte("0"), 0644)
	time.Sleep(pushtime * time.Millisecond)
	ioutil.WriteFile("/sys/class/gpio/gpio"+r.pin+"/value", []byte("1"), 0644)
	r.lpt = tn
	r.Unlock()
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Use POST methods.\n")
	case "POST":
		ramp, err := strconv.Atoi(r.URL.Query().Get("ramp"))
		if err != nil || ramp < 1 || ramp > 16 {
			fmt.Fprintf(w, "Sorry it doesn't work.")
			return
		}
		go pusher(ramp)
		fmt.Fprintf(w, "<h1>Pushed %d!</h1>", ramp)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	initbat()
	http.HandleFunc("/push", handler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
