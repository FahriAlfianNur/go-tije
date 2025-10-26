package mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fahri/go-tije/internal/config"
)

type Client struct {
	client mqtt.Client
}

func NewClient(cfg *config.MQTTConfig) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	
	opts.OnConnect = func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	}
	
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v\n", err)
	}
	
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect: %v", token.Error())
	}
	
	return &Client{client: client}, nil
}

func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, 1, handler)
	token.Wait()
	return token.Error()
}

func (c *Client) Publish(topic string, payload []byte) error {
	token := c.client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}