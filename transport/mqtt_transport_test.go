package transport

import "testing"
import (
	iotm "github.com/alivinco/iotmsglibgo"
	"time"
)

func TestMqttTransport_Publish(t *testing.T) {
	msg := iotm.NewIotMsg(iotm.MsgTypeCmd, "binary", "switch", nil)
	msg.SetDefaultStr("test value", "")
	mt := NewMqttTransport("tcp://localhost:1883","iotmsg_transport_test","","",true,false)
	mt.SetMessageHandler(func(topic string, iotMsg *iotm.IotMsg, domain string){
		t.Log("New message")
	})
	mt.Start()
	mt.Subscribe("test/iotmsgt",0)
	mt.Publish("test/iotmsgt",msg,0,"")
	time.Sleep(time.Second*1)
	mt.Stop()


}

func Test(t *testing.T) {

}