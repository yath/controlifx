package controlifx

/*
import (
	"bytes"
	"math"
	"reflect"
	"testing"
	_time "time"
)

func TestSendableLanMessage_MarshalBinary(t *testing.T) {
	o := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:0x1fff,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:0x1fff,
			},
		},
		Payload:&LightSetPowerLanMessage{
			Level:0xffff,
			Duration:0x1fffffff,
		},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x0, 0x34, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestSendableLanMessage_Payload(t *testing.T) {
	o := SendableLanMessage{}

	p := LightSetPowerLanMessage{
		Level:0xffff,
		Duration:0x1fffffff,
	}

	o.Payload = p

	if p != o.Payload {
		t.Error("payload was not set correctly")
	}
}

func TestReceivableLanMessage_UnmarshalBinary(t *testing.T) {
	o := ReceivableLanMessage{}

	b := []byte{0xff, 0x1f, 0x0, 0x34, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 3, 0x0, 0x0, 0x0, 0x1f, 0xff,
		0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := ReceivableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:0x1fff,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:3,
			},
		},
		Payload:&StateServiceLanMessage{
			Service:0x1f,
			Port:0x1fffffff,
		},
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanHeader_MarshalBinary(t *testing.T) {
	o := LanHeader{
		Frame:LanHeaderFrame{
			Size:0x1fff,
			Tagged:true,
			Source:0x1fffffff},
		FrameAddress:LanHeaderFrameAddress{
			Target:0x1fffffffffff,
			AckRequired:true,
			ResRequired:true,
			Sequence:0x1f,
		},
		ProtocolHeader:LanHeaderProtocolHeader{
			Type:0x1fff,
		},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x0, 0x34, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrame_MarshalBinary(t *testing.T) {
	o := LanHeaderFrame{
		Size:0x1fff,
		Tagged:true,
		Source:0x1fffffff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x0, 0x34, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrame_MarshalBinary2(t *testing.T) {
	o := LanHeaderFrame{
		Size:0x1fff,
		Tagged:false,
		Source:0x1fffffff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x0, 0x14, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrame_UnmarshalBinary(t *testing.T) {
	o := LanHeaderFrame{}

	b := []byte{0xff, 0x1f, 0x0, 0x34, 0xff, 0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LanHeaderFrame{
		Size:0x1fff,
		Tagged:true,
		Source:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanHeaderFrameAddress_MarshalBinary(t *testing.T) {
	o := LanHeaderFrameAddress{
		Target:0x1fffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrameAddress_MarshalBinary2(t *testing.T) {
	o := LanHeaderFrameAddress{
		Target:0x1fffffffffff,
		AckRequired:false,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrameAddress_UnmarshalBinary(t *testing.T) {
	o := LanHeaderFrameAddress{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LanHeaderFrameAddress{
		Target:0x1fffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanHeaderProtocolHeader_MarshalBinary(t *testing.T) {
	o := LanHeaderProtocolHeader{
		Type:0x1fff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0,
		0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderProtocolHeader_UnmarshalBinary(t *testing.T) {
	o := LanHeaderProtocolHeader{}

	b := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LanHeaderProtocolHeader{
		Type:0x1fff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLabel_MarshalBinary(t *testing.T) {
	o := Label("hello world")

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestPort_MarshalBinary(t *testing.T) {
	o := Port(0x1fffffff)

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestPowerLevel_MarshalBinary(t *testing.T) {
	o := PowerLevel(0xffff)

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestPowerLevel_MarshalBinary2(t *testing.T) {
	o := PowerLevel(0x0)

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestPowerLevel_MarshalBinary3(t *testing.T) {
	o := PowerLevel(0x1fff)

	_, err := o.MarshalBinary()
	if err == nil {
		t.Errorf("non 0 or 65535 value was erroneously allowed")
	}
}

func TestTime_Time(t *testing.T) {
	o := Time(1464000000000000000)

	time, err := o.Time()
	if err != nil {
		t.Error("error:", err)
	}

	expected := _time.Unix(0, 1464000000000000000)

	if time != expected {
		t.Errorf("expected '%#v', got '%#v'", expected, time)
	}
}

func TestTime_Time2(t *testing.T) {
	o := Time(math.MaxInt64 + 1)

	_, err := o.Time()
	if err == nil {
		t.Error("overflowing time was erroneously allowed")
	}
}

func TestTime_MarshalBinary(t *testing.T) {
	// 1464000000000000000
	o := Time(0x14512c3e4f2c0000)

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x00, 0x00, 0x2c, 0x4f, 0x3e, 0x2c, 0x51, 0x14}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestNewReceivablePayloadOfType(t *testing.T) {
	o, err := newReceivablePayloadOfType(3)
	if err != nil {
		t.Error("error:", o)
	}

	if _, ok := o.(*StateServiceLanMessage); !ok {
		t.Errorf("error: could not cast %T to StateServiceLanMessage", o)
	}
}

func TestNewReceivablePayloadOfType2(t *testing.T) {
	_, err := newReceivablePayloadOfType(4)
	if err == nil {
		t.Error("invalid payload type did not error")
	}
}

func TestLanDeviceMessageBuilder_Tagged(t *testing.T) {
	o := LanDeviceMessageBuilder{
		target:0x1,
	}

	if !o.Tagged() {
		t.Error("target was specified but Tagged() returned false")
	}
}

func TestLanDeviceMessageBuilder_Tagged2(t *testing.T) {
	o := LanDeviceMessageBuilder{}

	if o.Tagged() {
		t.Error("target was not specified but Tagged() returned true")
	}
}

func TestLanDeviceMessageBuilder_buildNormalMessageOfType(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.buildNormalMessageOfType(2)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:2,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestLanDeviceMessageBuilder_GetService(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetService()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:2,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateServiceLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateServiceLanMessage{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateServiceLanMessage{
		Service:0xff,
		Port:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetHostInfo(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetHostInfo()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:12,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateHostInfoLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateHostInfoLanMessage{}

	b := []byte{0xdb, 0x0f, 0x49, 0x40, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateHostInfoLanMessage{
		Signal:3.1415927,
		Tx:0x1fffffff,
		Rx:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetHostFirmware(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetHostFirmware()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:14,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateHostFirmwareLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateHostFirmwareLanMessage{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateHostFirmwareLanMessage{
		Build:0x1fffffffffffffff,
		Version:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetWifiInfo(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetWifiInfo()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:16,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateWifiInfoLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateWifiInfoLanMessage{}

	b := []byte{0xdb, 0x0f, 0x49, 0x40, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateWifiInfoLanMessage{
		Signal:3.1415927,
		Tx:0x1fffffff,
		Rx:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetWifiFirmware(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetWifiFirmware()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:18,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateWifiFirmwareLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateWifiFirmwareLanMessage{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateWifiFirmwareLanMessage{
		Build:0x1fffffffffffffff,
		Version:0x1fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetPower(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetPower()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:20,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestSetPowerLanMessage_MarshalBinary(t *testing.T) {
	o := SetPowerLanMessage{
		Level:0,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanDeviceMessageBuilder_SetPower(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	p := SetPowerLanMessage{
		Level:0xffff,
	}

	m := o.SetPower(p)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize + 2,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:21,
			},
		},
		Payload:p,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStatePowerLanMessage_UnmarshalBinary(t *testing.T) {
	o := StatePowerLanMessage{}

	b := []byte{0xff, 0xff}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StatePowerLanMessage{
		Level:0xffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetLabel(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetLabel()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:23,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestSetLabelLanMessage_MarshalBinary(t *testing.T) {
	o := SetLabelLanMessage{
		Label:"hello world",
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanDeviceMessageBuilder_SetLabel(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	p := SetLabelLanMessage{
		Label:"hello world",
	}

	m := o.SetLabel(p)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize + 32,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:24,
			},
		},
		Payload:p,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateLabelLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateLabelLanMessage{}

	b := []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c,
		0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateLabelLanMessage{
		Label:"hello world",
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetVersion(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetVersion()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:32,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateVersionLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateVersionLanMessage{}

	b := []byte{0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff, 0x2f, 0xff, 0xff,
		0xff, 0x3f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateVersionLanMessage{
		Vendor:0x1fffffff,
		Product:0x2fffffff,
		Version:0x3fffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetInfo(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetInfo()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:34,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateInfoLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateInfoLanMessage{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x2f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0x3f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateInfoLanMessage{
		Time:0x1fffffffffffffff,
		Uptime:0x2fffffffffffffff,
		Downtime:0x3fffffffffffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetLocation(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetLocation()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:48,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateLocationLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateLocationLanMessage{}

	b := []byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateLocationLanMessage{
		Location:[16]byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Label:"hello world",
		UpdatedAt:0x1fffffffffffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_GetGroup(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.GetGroup()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:51,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestStateGroupLanMessage_UnmarshalBinary(t *testing.T) {
	o := StateGroupLanMessage{}

	b := []byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20,
		0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := StateGroupLanMessage{
		Group:[16]byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		Label:"hello world",
		UpdatedAt:0x1fffffffffffffff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestEchoRequestLanMessage_MarshalBinary(t *testing.T) {
	o := EchoRequestLanMessage{
		Payload:[64]byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanDeviceMessageBuilder_EchoRequest(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	p := EchoRequestLanMessage{
		Payload:[64]byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	m := o.EchoRequest(p)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize + 64,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:58,
			},
		},
		Payload:p,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestEchoResponseLanMessage_UnmarshalBinary(t *testing.T) {
	o := EchoResponseLanMessage{}

	b := []byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := EchoResponseLanMessage{
		Payload:[64]byte{0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestHSBK_MarshalBinary(t *testing.T) {
	o := HSBK{
		Hue:0x1fff,
		Saturation:0x2fff,
		Brightness:0x3fff,
		Kelvin:0x1f40,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0xff, 0x2f, 0xff, 0x3f, 0x40, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestHSBK_MarshalBinary2(t *testing.T) {
	o := HSBK{
		Hue:0x1fff,
		Saturation:0x2fff,
		Brightness:0x3fff,
		Kelvin:2499,
	}

	_, err := o.MarshalBinary()
	if err == nil {
		t.Error("color temperature outside range was erroneously allowed")
	}
}

func TestHSBK_MarshalBinary3(t *testing.T) {
	o := HSBK{
		Hue:0x1fff,
		Saturation:0x2fff,
		Brightness:0x3fff,
		Kelvin:9001,
	}

	_, err := o.MarshalBinary()
	if err == nil {
		t.Error("color temperature outside range was erroneously allowed")
	}
}

func TestHSBK_UnmarshalBinary(t *testing.T) {
	o := HSBK{}

	b := []byte{0xff, 0x1f, 0xff, 0x2f, 0xff, 0x3f, 0xff, 0x4f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := HSBK{
		Hue:0x1fff,
		Saturation:0x2fff,
		Brightness:0x3fff,
		Kelvin:0x4fff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_LightGet(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.LightGet()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:101,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestLightSetColorLanMessage_MarshalBinary(t *testing.T) {
	o := LightSetColorLanMessage{
		Color:HSBK{
			Hue:0x1fff,
			Saturation:0x2fff,
			Brightness:0x3fff,
			Kelvin:0x1f40,
		},
		Duration:0x1fffffff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0x0, 0xff, 0x1f, 0xff, 0x2f, 0xff, 0x3f, 0x40, 0x1f,
		0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanDeviceMessageBuilder_LightSetColor(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	p := LightSetColorLanMessage{
		Color:HSBK{
			Hue:0x1fff,
			Saturation:0x2fff,
			Brightness:0x3fff,
			Kelvin:0x4fff,
		},
		Duration:0x1fffffff,
	}

	m := o.LightSetColor(p)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize + 13,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:102,
			},
		},
		Payload:p,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestLightStateLanMessage_UnmarshalBinary(t *testing.T) {
	o := LightStateLanMessage{}

	b := []byte{0xff, 0x1f, 0xff, 0x2f, 0xff, 0x3f, 0xff, 0x4f, 0xff, 0x1f,
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LightStateLanMessage{
		Color:HSBK{
			Hue:0x1fff,
			Saturation:0x2fff,
			Brightness:0x3fff,
			Kelvin:0x4fff,
		},
		Power:0x1fff,
		Label:"hello world",
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}

func TestLanDeviceMessageBuilder_LightGetPower(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	m := o.LightGetPower()

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:116,
			},
		},
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestLightSetPowerLanMessage_MarshalBinary(t *testing.T) {
	o := LightSetPowerLanMessage{
		Level:0xffff,
		Duration:0x1fffffff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanDeviceMessageBuilder_LightSetPower(t *testing.T) {
	o := LanDeviceMessageBuilder{
		source:0x1fffffff,
		target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	p := LightSetPowerLanMessage{
		Level:0xffff,
		Duration:0x1fffffff,
	}

	m := o.LightSetPower(p)

	expected := SendableLanMessage{
		Header:LanHeader{
			Frame:LanHeaderFrame{
				Size:LanHeaderSize + 6,
				Tagged:true,
				Source:0x1fffffff,
			},
			FrameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			ProtocolHeader:LanHeaderProtocolHeader{
				Type:117,
			},
		},
		Payload:p,
	}

	if !reflect.DeepEqual(expected, m) {
		t.Errorf("expected '%#v', got '%#v'", expected, m)
	}
}

func TestLightStatePowerLanMessage_UnmarshalBinary(t *testing.T) {
	o := LightStatePowerLanMessage{}

	b := []byte{0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LightStatePowerLanMessage{
		Level:0x1fff,
	}

	if !reflect.DeepEqual(expected, o) {
		t.Errorf("expected '%#v', got '%#v'", expected, o)
	}
}
*/
