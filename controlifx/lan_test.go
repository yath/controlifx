package controlifx

import (
	"fmt"
	"testing"
)

func TestLanHeaderFrame_MarshalBinary(t *testing.T) {
	// Check little-endianness and sub-byte values (by changing [Tagged])

	o := LanHeaderFrame{Size:0x1fff, Tagged:true, Source:0x1fffffff}
	b, _ := o.MarshalBinary()
	s := fmt.Sprintf("%#v", b)

	const Expected1 = "[]byte{0xff, 0x1f, 0x38, 0x0, 0xff, 0xff, 0xff, 0x1f}"

	if s != Expected1 {
		t.Errorf("expected '%s', got '%s'", Expected1, s)
	}

	o = LanHeaderFrame{Size:0x1fff, Tagged:false, Source:0x1fffffff}
	b, _ = o.MarshalBinary()
	s = fmt.Sprintf("%#v", b)

	const Expected2 = "[]byte{0xff, 0x1f, 0x18, 0x0, 0xff, 0xff, 0xff, 0x1f}"

	if s != Expected2 {
		t.Errorf("expected '%s', got '%s'", Expected2, s)
	}
}

func TestLanHeaderFrameAddress_MarshalBinary(t *testing.T) {
	// Check little-endianness and sub-byte values (by changing [{Ack,Res}Required])

	o := LanHeaderFrameAddress{Target:0x1fffffffffffffff, AckRequired:true, ResRequired:true, Sequence:0x1f}
	b, _ := o.MarshalBinary()
	s := fmt.Sprintf("%#v", b)

	const Expected1 = "[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x1f}"

	if s != Expected1 {
		t.Errorf("expected '%s', got '%s'", Expected1, s)
	}

	o = LanHeaderFrameAddress{Target:0x1fffffffffffffff, AckRequired:false, ResRequired:true, Sequence:0x1f}
	b, _ = o.MarshalBinary()
	s = fmt.Sprintf("%#v", b)

	const Expected2 = "[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1f}"

	if s != Expected2 {
		t.Errorf("expected '%s', got '%s'", Expected2, s)
	}
}

func TestLanHeaderProtocolHeader_MarshalBinary(t *testing.T) {
	o := LanHeaderProtocolHeader{Type:0x1fff}
	b, _ := o.MarshalBinary()
	s := fmt.Sprintf("%#v", b)

	const Expected = "[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0x1f, 0x0, 0x0}"

	if s != Expected {
		t.Errorf("expected '%s', got '%s'", Expected, s)
	}
}
