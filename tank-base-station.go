package main

import (
	"fmt"
	"time"
	_ "time/tzdata"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	doneFilling(client, msg.Payload())
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func main() {
	// var broker = "69.164.210.24"
	var broker = "10.30.30.23"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	/*
		opts.SetClientID("go_mqtt_client")
		opts.SetUsername("emqx")
		opts.SetPassword("public")
	*/
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.KeepAlive = 10

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}
	time.Local = loc

	fmt.Println(time.Now().Hour())

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// TODO - change to hours
	ticker := time.NewTicker(6 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				checkCycleTime(client)
				//fmt.Println("Tick at", t)
			}
		}
	}()

	sub(client, "tank/doneFilling")
	//sub(client, "topic/test")
	for {
	}

}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func publish(client mqtt.Client, topic string, timeToFill uint32) {
	token := client.Publish(topic, 0, false, timeToFill)
	token.Wait()
	time.Sleep(time.Second)
}

func doneFilling(client mqtt.Client, msgPayload []byte) {
	type T struct {
		Tank        string
		ElapsedFill uint32
	}
	fmt.Print(msgPayload)
	// TODO - save to DB
}

func checkCycleTime(client mqtt.Client) {
	const cycleTime uint32 = 60000
	// TODO - get from DB
	const timeRefilled uint32 = 0

	const tankACycleTime uint32 = cycleTime - timeRefilled

	publish(client, "tank/Cycle/A", tankACycleTime)
	publish(client, "tank/Cycle/B", cycleTime)

}
