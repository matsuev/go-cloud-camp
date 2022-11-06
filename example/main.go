package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cloud-camp/client"
	"log"
	"time"
)

type AppConfig struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}

func main() {
	cfgClient, err := client.Connect("http://localhost:8080/config", "example")
	if err != nil {
		log.Fatalln(err)
	}

	if err = cfgClient.DeleteConfig(context.Background()); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(1 * time.Second)

	newCfg := &AppConfig{
		Key1: "Value1",
		Key2: "Value2",
	}

	if err = cfgClient.CreateConfig(context.Background(), newCfg); err != nil {
		log.Fatalln(err)
	}

	readedCfg := &AppConfig{}

	if err = cfgClient.ReadAndDecodeConfig(context.Background(), readedCfg); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Initial config: %#v\n", readedCfg)

	if err = cfgClient.AssignRefreshCallback(2*time.Second, LogConfigRefresh); err != nil {
		log.Println(err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		cfgUpdate := &AppConfig{
			Key1: fmt.Sprintf("Value-%d", i*1000),
			Key2: fmt.Sprintf("Value-%d", i*2000),
		}
		err := cfgClient.UpdateConfig(context.Background(), cfgUpdate)
		if err != nil {
			log.Println(err)
		}
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Update completed. Press <return> key")

	var ss string
	fmt.Scanln(&ss)
}

func LogConfigRefresh(data []byte) {
	newCfg := &AppConfig{}

	err := json.Unmarshal(data, newCfg)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Updated config: %#v\n", newCfg)
}
