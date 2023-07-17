package mqtt33

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

//"192.168.3.212",
//"1883",

func TestName(t *testing.T) {
	q := NewMqtt(
		MqConfig{
			Port:     1883,
			Url:      "192.168.3.212",
			Broker:   "",
			UserName: "",
			Password: "",
			ClientID: "xxx",
			Qos:      0,
		},
	)
	q.SetHandlerReceMsg(func(client MQTT.Client, msg MQTT.Message) {
		println(string(msg.Topic()), string(msg.Payload()))
	})
	go q.Start()
	time.Sleep(2 * time.Second)
	q.Sub("xxx/v1")
	func() {
		for {
			data := map[string]string{
				"type": "real_data",
				"key":  "0",
				"data": "50.1",
			}
			v, _ := json.Marshal(data)
			q.Pub("xxx/v1", v)
			println(time.Now().String())
			time.Sleep(3 * time.Second)
		}
	}()
}
