package controlifx

import (
	"bytes"
	"math"
	"reflect"
	"testing"
	_time "time"
)

func TestSendableLanMessage_MarshalBinary(t *testing.T) {
	o := SendableLanMessage{
		header:LanHeader{
			frame:LanHeaderFrame{
				Size:0x1fff,
				Tagged:true,
				Source:0x1fffffff,
			},
			frameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			protocolHeader:LanHeaderProtocolHeader{
				Type:0x1fff,
			},
		},
		payload:&LightSetPowerLanMessage{
			level:0xffff,
			duration:0x1fffffff,
		},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestSendableLanMessage_Payload(t *testing.T) {
	o := SendableLanMessage{}

	p := LightSetPowerLanMessage{
		level:0xffff,
		duration:0x1fffffff,
	}

	o.Payload(p)

	if p != o.payload {
		t.Error("payload was not set correctly")
	}
}

func TestReceivableLanMessage_UnmarshalBinary(t *testing.T) {
	o := ReceivableLanMessage{}

	b := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 3, 0x0, 0x0, 0x0, 0x1f, 0xff,
		0xff, 0xff, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := ReceivableLanMessage{
		header:LanHeader{
			frame:LanHeaderFrame{
				Size:0x1fff,
				Tagged:true,
				Source:0x1fffffff,
			},
			frameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			protocolHeader:LanHeaderProtocolHeader{
				Type:3,
			},
		},
		payload:&StateServiceLanMessage{
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
		frame:LanHeaderFrame{
			Size:0x1fff,
			Tagged:true,
			Source:0x1fffffff},
		frameAddress:LanHeaderFrameAddress{
			Target:0x1fffffffffffffff,
			AckRequired:true,
			ResRequired:true,
			Sequence:0x1f,
		},
		protocolHeader:LanHeaderProtocolHeader{
			Type:0x1fff,
		},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
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

	expected := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f}

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

	expected := []byte{0xff, 0x1f, 0x18, 0x0, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrame_UnmarshalBinary(t *testing.T) {
	o := LanHeaderFrame{}

	b := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f}

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
		Target:0x1fffffffffffffff,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrameAddress_MarshalBinary2(t *testing.T) {
	o := LanHeaderFrameAddress{
		Target:0x1fffffffffffffff,
		AckRequired:false,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error:", err)
	}

	expected := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrameAddress_UnmarshalBinary(t *testing.T) {
	o := LanHeaderFrameAddress{}

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3, 0x1f}

	if err := o.UnmarshalBinary(b); err != nil {
		t.Error("error:", err)
	}

	expected := LanHeaderFrameAddress{
		Target:0x1fffffffffffffff,
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
	o := label("hello world")

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
	o := port(0x1fffffff)

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
	o := powerLevel(0xffff)

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
	o := powerLevel(0x0)

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
	o := powerLevel(0x1fff)

	_, err := o.MarshalBinary()
	if err == nil {
		t.Errorf("non 0 or 65535 value was erroneously allowed")
	}
}

func TestTime_Time(t *testing.T) {
	o := time(1464000000000000000)

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
	o := time(math.MaxInt64 + 1)

	_, err := o.Time()
	if err == nil {
		t.Error("overflowing time was erroneously allowed")
	}
}

func TestTime_MarshalBinary(t *testing.T) {
	// 1464000000000000000
	o := time(0x14512c3e4f2c0000)

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
	o, err := NewReceivablePayloadOfType(3)
	if err != nil {
		t.Error("error:", o)
	}

	o, ok := o.(*StateServiceLanMessage)
	if !ok {
		t.Error("error: could not cast %T to StateServiceLanMessage", o)
	}

	expected := &StateServiceLanMessage{}

	if reflect.TypeOf(expected) != reflect.TypeOf(o) {
		t.Errorf("expected %T, got %T", expected, o)
	}
}

func TestNewReceivablePayloadOfType2(t *testing.T) {
	_, err := NewReceivablePayloadOfType(4)
	if err == nil {
		t.Error("invalid payload type did not error")
	}
}
