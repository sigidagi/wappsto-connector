package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

type WappConfig struct {
	Username string
	Password string
	Url      string
}

type Meta struct {
	id      string `json:"id"`
	type_   string `json:"type"`
	version string `json:"version"`
}

type Session struct {
	meta Meta `json:"meta"`
}

func main() {

	f := "wappsto.toml"
	if _, err := os.Stat(f); err != nil {
		log.Fatal(err)
	}

	var config WappConfig
	if _, err := toml.DecodeFile(f, &config); err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %v\n", config.Username)

	data, err := json.Marshal(map[string]string{
		"username": config.Username,
		"password": config.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	var endpoint = "/services/session"
	url := config.Url + endpoint
	log.Printf("url: %s\n", url)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	reqDump, err := httputil.DumpRequestOut(request, true)
	fmt.Printf("\nREQUEST:\n%s", string(reqDump))

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	respDump, err := httputil.DumpResponse(response, true)
	fmt.Printf("\nRESPONSE:\n%s", string(respDump))

	defer response.Body.Close()

	log.Printf("Response %v\n", json.NewDecoder(response.Body))

	var session Session
	json.NewDecoder(response.Body).Decode(&session)

	fmt.Printf("session uuid: %s\n", session.meta.id)
	fmt.Println("Done")

}
