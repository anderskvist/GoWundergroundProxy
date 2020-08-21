package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"time"

	log "github.com/anderskvist/GoHelpers/log"
	"github.com/anderskvist/GoWundergroundProxy/wunderground"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"
)

var pubConnection mqtt.Client
var mqttClientID = "GoWundergroundProxy-" + strconv.Itoa(os.Getpid())
var cfg *ini.File

func connect(clientID string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	log.Debug("Connecting to MQTT (pub)")
	client.Publish("wunderground/GoWundergroundProxyOnline", 1, false, "true")
	return client
}

func createClientOptions(clientID string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	return opts
}

func main() {
	cfg, _ = ini.Load(os.Args[1])

	mqttURL := cfg.Section("mqtt").Key("url").String()
	uri, err := url.Parse(mqttURL)
	if err != nil {
		log.Fatal(err)
	}

	if pubConnection == nil {
		pubConnection = connect(mqttClientID+"-pub", uri)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for sig := range c {
			log.Debug("ping")

			if sig == os.Interrupt {
				log.Debug(sig)
				pubConnection.Publish("wunderground/GoWundergroundProxyOnline", 1, false, "false")
				// we need to sleep just a bit (1ms seems to be enough), to allow us to send info to the mqtt broker
				time.Sleep(10 * time.Millisecond)
				os.Exit(0)
			}
		}
	}()

	http.HandleFunc("/weatherstation/updateweatherstation.php", handler)
	http.ListenAndServe(":80", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Got request")
	keys := r.URL.Query()
	id := keys["ID"][0]

	wunderground := wunderground.Parse(keys)
	elements := reflect.ValueOf(&wunderground).Elem()
	types := elements.Type()

	for i := 0; i < elements.NumField(); i++ {
		pubConnection.Publish("wunderground/"+id+"/"+types.Field(i).Name, 0, false, elements.Field(i).String())
	}

	// send request to wunderground (we stole the hostname, so we need to resolve it externally)
	if true {
		log.Debug("Sending request on to the real rtupdate.wunderground.com")
		req, err := http.NewRequest("GET", "http://169.47.111.60"+r.URL.RequestURI(), nil)
		req.Host = "rtupdate.wunderground.com"

		client := &http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		defer resp.Body.Close()
		log.Debug("HTTP Status: ", resp.StatusCode)
	}
}
