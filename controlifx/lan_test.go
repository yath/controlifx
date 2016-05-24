package controlifx

import (
	"bytes"
	"testing"
)

func TestLanMessage_MarshalBinary(t *testing.T) {
	o := SendableLanMessage{
		header:LanHeader{
			frame:LanHeaderFrame{
				Size:0x1fff,
				Tagged:true,
				Source:0x1fffffff,
			},
			frameAddress:LanHeaderFrameAddress{
				Target:0x1fffffffffffffff,
				TargetMac:false,
				AckRequired:true,
				ResRequired:true,
				Sequence:0x1f,
			},
			protocolHeader:LanHeaderProtocolHeader{
				Type:0x1fff,
			},
		},
		payload:&LanHeaderFrame{
			Size:0x1fff,
			Tagged:true,
			Source:0x1fffffff,
		},
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0,
		0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
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
			TargetMac:false,
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
		t.Error("error", err)
	}

	expected := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}

func TestLanHeaderFrame_MarshalBinary(t *testing.T) {
	// Check little-endianness and sub-byte values (by changing [Tagged]).

	o := LanHeaderFrame{
		Size:0x1fff,
		Tagged:true,
		Source:0x1fffffff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected1 := []byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected1, b) {
		t.Errorf("expected '%#v', got '%#v'", expected1, b)
	}

	o = LanHeaderFrame{
		Size:0x1fff,
		Tagged:false,
		Source:0x1fffffff,
	}

	b, err = o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected2 := []byte{0xff, 0x1f, 0x18, 0x0, 0xff, 0xff, 0xff, 0x1f}

	if !bytes.Equal(expected2, b) {
		t.Errorf("expected '%#v', got '%#v'", expected2, b)
	}
}

func TestLanHeaderFrameAddress_MarshalBinary(t *testing.T) {
	// Check little-endianness and sub-byte values (by changing [{Ack,Res}Required]).

	o := LanHeaderFrameAddress{
		Target:0x1fffffffffffffff,
		TargetMac:false,
		AckRequired:true,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected1 := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f}

	if !bytes.Equal(expected1, b) {
		t.Errorf("expected '%#v', got '%#v'", expected1, b)
	}

	o = LanHeaderFrameAddress{
		Target:0x1fffffffffffffff,
		TargetMac:false,
		AckRequired:false,
		ResRequired:true,
		Sequence:0x1f,
	}

	b, err = o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected2 := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1f}

	if !bytes.Equal(expected2, b) {
		t.Errorf("expected '%#v', got '%#v'", expected2, b)
	}
}

func TestLanHeaderProtocolHeader_MarshalBinary(t *testing.T) {
	// Check little-endianness.

	o := LanHeaderProtocolHeader{
		Type:0x1fff,
	}

	b, err := o.MarshalBinary()
	if err != nil {
		t.Error("error", err)
	}

	expected := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0,
		0x0}

	if !bytes.Equal(expected, b) {
		t.Errorf("expected '%#v', got '%#v'", expected, b)
	}
}
