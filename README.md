# mqtt33

golang mqtt包


获取包

```
go get -u  "github.com/eclipse/paho.mqtt.golang"
go get -u "github.com/ituserxxx/mqtt33"
```

使用示例 demo：

```go
package main

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	MQTT33 "github.com/ituserxxx/mqtt33"
	"time"
)

func main() {
	q := MQTT33.NewMqtt(MQTT33.MqConfig{
		Port:     1883,
		Url:      "192.168.3.212",
		Broker:   "",
		UserName: "",
		Password: "",
		ClientID: "xxx",
		Qos:      0,
	})
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
	
	/* out
	2023-07-17 14:19:43.2555954 +0800 CST m=+5.005994501
	xxx/v1 {"data":"50.1","key":"0","type":"real_data"}
	2023-07-17 14:19:46.2556641 +0800 CST m=+8.006063201
	xxx/v1 {"data":"50.1","key":"0","type":"real_data"}
	2023-07-17 14:19:49.2575503 +0800 CST m=+11.007949401
	xxx/v1 {"data":"50.1","key":"0","type":"real_data"}
	
	*/
}
```

