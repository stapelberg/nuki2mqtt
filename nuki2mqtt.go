package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	listenAddress = flag.String("listen",
		":8319",
		"listen address for Nuki Bridge webhooks")

	mqttBroker = flag.String("mqtt_broker",
		"tcp://dr.lan:1883",
		"MQTT broker address for github.com/eclipse/paho.mqtt.golang")

	mqttTopic = flag.String("mqtt_topic",
		"zkj-nuki/webhook",
		"MQTT topic for publishing webhook contents")

	mqttAdvancedTopic = flag.String("mqtt_advanced_topic",
		"zkj-nuki/advancedhook",
		"MQTT topic for publishing webhook contents")
)

func nukiBridge() error {
	opts := mqtt.NewClientOptions().AddBroker(*mqttBroker)
	opts.SetClientID("nuki2mqtt")
	opts.SetConnectRetry(true)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connection failed: %v", token.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/nuki", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("[bridge] remote: %s, headers: %v, body: %q", r.RemoteAddr, r.Header, string(b))

		mqttClient.Publish(*mqttTopic, 0 /* qos */, true /* retained */, string(b))
	})

	mux.HandleFunc("/advanced", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("[advanced] remote: %s, headers: %v, body: %q", r.RemoteAddr, r.Header, string(b))

		mqttClient.Publish(*mqttAdvancedTopic, 0 /* qos */, true /* retained */, string(b))

		// Nuki Advanced API requires a HTTP 204 response
		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("http.ListenAndServe(%q)", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, mux); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if err := nukiBridge(); err != nil {
		log.Fatal(err)
	}
}
