package wappsto

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"

	//"fmt"
	//"github.com/BurntSushi/toml"
	//"github.com/Shopify/sarama"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"net/http/httputil"
	"time"
	"wappsto-kafka-connector/internal/config"
	"wappsto-kafka-connector/internal/connector"
)

type Meta struct {
	Id      string `json:"id"`
	Type_   string `json:"type"`
	Version string `json:"version"`
}

type Session struct {
	Meta Meta `json:"meta"`
}

type Stream struct {
	Subscription []string `json:"subscription"`
	Meta         Meta     `json:"meta"`
}

// Global variables
var session Session
var stream Stream
var done chan interface{}

func Setup(conf config.Config) error {

	var endpoint = "/services/session"
	url := "https://" + config.C.Wappsto.Server + endpoint
	data, _ := json.Marshal(map[string]string{
		"username": conf.Wappsto.Username,
		"password": conf.Wappsto.Password,
	})

	log.Printf("url: %s\n", url)

	result, err := requestWappsto(http.MethodPost, url, data)
	if err != nil {
		return errors.Wrap(err, "wappsto connection error")
	}

	log.Printf("Success with Wappsto connection '%s', user: %s\n", url, conf.Wappsto.Username)

	err = mapstructure.Decode(result, &session)
	if err != nil {
		return errors.Wrap(err, "wappsto identification error; check username and password")
	}

	return nil
}

func requestWappsto(method, url string, data []byte) (interface{}, error) {

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")

		if (session.Meta != Meta{}) {
			request.Header.Set("X-Session", session.Meta.Id)
		}
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	reqDump, err := httputil.DumpRequestOut(request, true)
	log.Printf("\nREQUEST:\n%s", string(reqDump))

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	/*
	 *respDump, err := httputil.DumpResponse(response, true)
	 *fmt.Printf("\nRESPONSE:\n%s", string(respDump))
	 */

	defer response.Body.Close()

	var out interface{}
	json.NewDecoder(response.Body).Decode(&out)
	return out, nil
}

func receiveHandler(connection *websocket.Conn) {
	// once received handler
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
		err = connector.PushMessageToQueue("wappsto", msg)
		if err != nil {
			log.Printf("Error sending message to Kafka: %s\n", err.Error())
		}
	}
}

// Public.
func HandleWappstoStream() error {

	done = make(chan interface{}) // Channel to indicate that the receiverHandler is done

	endpoint := "/services/stream"
	url := "https://" + config.C.Wappsto.Server + endpoint
	data, _ := json.Marshal(map[string][]string{
		"subscription": {"/network"},
	})

	result, err := requestWappsto(http.MethodPost, url, data)
	if err != nil {
		return errors.Wrap(err, "Failed creating Wappsto stream")
	}

	err = mapstructure.Decode(result, &stream)
	if err != nil {
		return errors.Wrap(err, "Failed to get Wappsto stream ID")
	}

	endpoint = "/services/websocket/" + stream.Meta.Id + "?X-Session=" + session.Meta.Id
	url = "wss://" + config.C.Wappsto.Server + endpoint

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to create Wappsto websocket")
	}

	log.Println("Stream to/from Wappsto has been opened; waiting for messages")

	defer conn.Close()

	go receiveHandler(conn)

	for {
		select {
		case <-done:
			log.Println("Wappsto stream is closed! Exiting!")
			// TODO handle stream reconnection.
			return errors.New("Waapsto stream closed!")
		}
	}
}
