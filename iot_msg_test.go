package iotmsglibgo

import (
	"testing"
)

func TestNewIotMsg(t *testing.T) {
	msg := NewIotMsg(PayloadTypeJsonIotMsgV1, "binary", "switch", nil)
	msg.SetDefaultStr("test value", "")
	msg.SetStrProperty("prop1","value1")
	t.Log(msg.String())
	if msg.GetDefaultStr() != "test value" {
		t.Failed()
	}
	if msg.GetStrProperty("prop1") != "value1" {
		t.Failed()
	}
}
