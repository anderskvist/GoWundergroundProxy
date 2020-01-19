package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	log "github.com/anderskvist/GoHelpers/log"
	"github.com/anderskvist/GoWundergroundProxy/wunderground"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ini "gopkg.in/ini.v1"
)

var pubConnection mqtt.Client
var mqttClientID = "GoWundergroundProxy" + string(os.Getpid())
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
		pubConnection = connect(mqttClientID+"pub", uri)
		log.Debug("Connecting to MQTT (pub)")
		pubConnection.Publish("wunderground/GoWundergroundProxyOnline", 0, false, "true")
	}

	http.HandleFunc("/weatherstation/updateweatherstation.php", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	id := keys["ID"][0]

	wunderground := wunderground.Parse(keys)
	elements := reflect.ValueOf(&wunderground).Elem()
	types := elements.Type()

	for i := 0; i < elements.NumField(); i++ {
		pubConnection.Publish("wunderground/"+id+"/"+types.Field(i).Name, 0, false, elements.Field(i).String())
	}

	// send request to wunderground (we stole the hostname, so we need to resolve it externally)
}
