package mqtt33

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

/*
	Qos byte
   消息服务质量，一共有三个：
   0：尽力而为。消息可能会丢，但绝不会重复传输
   1：消息绝不会丢，但可能会重复传输
   2：恰好一次。每条消息肯定会被传输一次且仅传输一次
*/
//func init() {
//MQTT.ERROR = log.New(os.Stdout,"【ERROR】",0)
//MQTT.CRITICAL = log.New(os.Stdout,"【CRITICAL】",0)
//MQTT.WARN = log.New(os.Stdout,"【WARN】",0)
//MQTT.DEBUG = log.New(os.Stdout,"【DEBUG】",0)
//tools.DebugInfo("Init mqtt")
//}


type MqConfig struct {
	Port      int
	Url       string
	Broker    string
	UserName  string
	Password  string
	ClientID  string
	Qos       byte
}
type MqServer struct {
	Id            string
	Cli           MQTT.Client
	Msghandlefunc MQTT.MessageHandler
	ConnHandle    MQTT.OnConnectHandler
	Config        MqConfig
	End           chan bool
}
func NewMqtt(config MqConfig) *MqServer {
	m1 := &MqServer{
		Id:        fmt.Sprintf("id-%s-%s", config.Url, config.ClientID),
		Config: config,
		End:       make(chan bool),
	}
	return m1
}
func (m *MqServer)SetConnHandler(fn func(_ MQTT.Client)){
	m.ConnHandle = fn
}
func (m *MqServer) SetHandlerReceMsg(fn func(client MQTT.Client, msg MQTT.Message)){
	m.Msghandlefunc=fn
}
func (m *MqServer) handleRece(topic, payload string) {
	fmt.Println(fmt.Printf("\n[%s]%s", topic, payload))

}

func (m *MqServer) Start() {
	m.initMqtt()
	defer m.DisConn()
	for {
		select {
		case <-m.End:
			m.DisConn()
			return
		default:
			time.Sleep( time.Second)
		}
	}
}
func (m *MqServer) initMqtt() {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", m.Config.Url, m.Config.Port))
	opts.SetClientID(m.Config.ClientID)
	opts.SetUsername(m.Config.UserName)
	opts.SetPassword(m.Config.Password)
	opts.SetDefaultPublishHandler(m.Msghandlefunc)
	opts.OnConnect = m.ConnHandle
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	m.Cli = MQTT.NewClient(opts)
	if token := m.Cli.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("xxx",token.Error())
	}
}
func (m *MqServer) DisConn() {
	if m.Cli.IsConnected() {
		m.Cli.Disconnect(100)
	}
}
func (m *MqServer) Sub(topic string) error {
	if m.Cli.IsConnected() {
		token := m.Cli.Subscribe(topic, m.Config.Qos, m.Msghandlefunc)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
func (m *MqServer) Pub(pushTopic string,msg interface{}) error {
	if m.Cli.IsConnected() {
		token := m.Cli.Publish(pushTopic, m.Config.Qos, false, msg)
		if token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	return nil
}
