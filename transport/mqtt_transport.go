package transport
import (
	"github.com/alivinco/iotmsglibgo"
	"strings"
	log "github.com/Sirupsen/logrus"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttAdapter , mqtt adapter .
type MqttTransport struct {
	client     MQTT.Client
	msgHandler MessageHandler
	useDomains bool
}

type MessageHandler func(topic string, iotMsg *iotmsglibgo.IotMsg, domain string)

// NewMqttAdapter constructor
//serverUri="tcp://localhost:1883"
func NewMqttTransport(serverURI string, clientID string, username string, password string,cleanSession bool,useDomains bool) *MqttTransport {
	mh := MqttTransport{}
	opts := MQTT.NewClientOptions().AddBroker(serverURI)
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(mh.onMessage)
	opts.SetCleanSession(cleanSession)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(mh.onConnectionLost)
	opts.SetOnConnectHandler(mh.onConnect)
	//create and start a client using the above ClientOptions
	mh.client = MQTT.NewClient(opts)
	mh.useDomains = useDomains
	return &mh
}

// SetMessageHandler message handler setter
func (mh *MqttTransport) SetMessageHandler(msgHandler MessageHandler) {
	mh.msgHandler = msgHandler
}

// Start , starts adapter async.
func (mh *MqttTransport) Start() error {
	log.Info("Connecting to MQTT broker ")
	if token := mh.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Stop , stops adapter.
func (mh *MqttTransport) Stop() {
	mh.client.Disconnect(250)
}

// Subscribe - subscribing for topic
func (mh *MqttTransport) Subscribe(topic string, qos byte) error {
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	log.Debug("Subscribing to topic:", topic)
	if token := mh.client.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return token.Error()
	}
	return nil
}

// Unsubscribe , unsubscribing from topic
func (mh *MqttTransport) Unsubscribe(topic string) error {
	log.Debug("Unsubscribing from topic:", topic)
	if token := mh.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (mh *MqttTransport) onConnectionLost(client MQTT.Client, err error) {
	log.Errorf("Connection lost with MQTT broker . Error : %v", err)
}

func (mh *MqttTransport) onConnect(client MQTT.Client) {
	log.Infof("Connection established with MQTT broker .")
}

//define a function for the default message handler
func (mh *MqttTransport) onMessage(client MQTT.Client, msg MQTT.Message) {
	log.Debugf(" New msg from TOPIC: %s", msg.Topic())
	// log.Debug("MSG: %s\n", msg.Payload())
	var domain , topic string
	if mh.useDomains {
		domain, topic = DetachDomainFromTopic(msg.Topic())
	}else {
		topic = msg.Topic()
	}
	iotMsg, err := iotmsglibgo.ConvertBytesToIotMsg(topic, msg.Payload(), nil)
	if err == nil {
		mh.msgHandler(topic, iotMsg, domain)
	} else {
		log.Error(err)

	}
}

// Publish iotMsg
func (mh *MqttTransport) Publish(topic string, iotMsg *iotmsglibgo.IotMsg, qos byte, domain string) error {
	bytm, err := iotmsglibgo.ConvertIotMsgToBytes(topic, iotMsg, nil)
	if domain != "" && mh.useDomains {
		topic = AddDomainToTopic(domain, topic)
	}
	if err == nil {
		log.Debug("Publishing msg to topic:", topic)
		mh.client.Publish(topic, qos, false, bytm)
		return nil
	}
	return err

}

// AddDomainToTopic , adds prefix to topic .
func AddDomainToTopic(domain string, topic string) string {
	// Check if topic is already prefixed with  "/" if yes then concat without adding "/"
	// 47 is code of "/"
	if topic[0] == 47 {
		return domain + topic
	}
	return domain + "/" + topic
}

// DetachDomainFromTopic detaches domain from topic
func DetachDomainFromTopic(topic string) (string, string) {
	spt := strings.Split(topic, "/")
	// spt[0] - domain
	var top string
	if strings.Contains(spt[1], "jim") {
		top = strings.Replace(topic, spt[0]+"/", "", 1)
	} else {
		top = strings.Replace(topic, spt[0], "", 1)
	}
	// returns domain , topic
	return spt[0], top

}

