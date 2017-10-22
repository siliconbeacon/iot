// Code generated by protoc-gen-go. DO NOT EDIT.
// source: orientation.proto

/*
Package messages is a generated protocol buffer package.

It is generated from these files:
	orientation.proto
	weather.proto

It has these top-level messages:
	OrientationReadings
	OrientationReading
	GyroscopeReading
	MagnetometerReading
	AccelerometerReading
	AttitudeHeadingReference
	WeatherReadings
	WeatherReading
*/
package messages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type OrientationReadings struct {
	Device   string                     `protobuf:"bytes,1,opt,name=device" json:"device,omitempty"`
	BaseTime *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=base_time,json=baseTime" json:"base_time,omitempty"`
	Readings []*OrientationReading      `protobuf:"bytes,3,rep,name=readings" json:"readings,omitempty"`
}

func (m *OrientationReadings) Reset()                    { *m = OrientationReadings{} }
func (m *OrientationReadings) String() string            { return proto.CompactTextString(m) }
func (*OrientationReadings) ProtoMessage()               {}
func (*OrientationReadings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *OrientationReadings) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *OrientationReadings) GetBaseTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.BaseTime
	}
	return nil
}

func (m *OrientationReadings) GetReadings() []*OrientationReading {
	if m != nil {
		return m.Readings
	}
	return nil
}

type OrientationReading struct {
	RelativeTimeUs uint32                    `protobuf:"varint,1,opt,name=relative_time_us,json=relativeTimeUs" json:"relative_time_us,omitempty"`
	Gyroscope      *GyroscopeReading         `protobuf:"bytes,2,opt,name=gyroscope" json:"gyroscope,omitempty"`
	Magnetometer   *MagnetometerReading      `protobuf:"bytes,3,opt,name=magnetometer" json:"magnetometer,omitempty"`
	Accelerometer  *AccelerometerReading     `protobuf:"bytes,4,opt,name=accelerometer" json:"accelerometer,omitempty"`
	Ahrs           *AttitudeHeadingReference `protobuf:"bytes,5,opt,name=ahrs" json:"ahrs,omitempty"`
}

func (m *OrientationReading) Reset()                    { *m = OrientationReading{} }
func (m *OrientationReading) String() string            { return proto.CompactTextString(m) }
func (*OrientationReading) ProtoMessage()               {}
func (*OrientationReading) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *OrientationReading) GetRelativeTimeUs() uint32 {
	if m != nil {
		return m.RelativeTimeUs
	}
	return 0
}

func (m *OrientationReading) GetGyroscope() *GyroscopeReading {
	if m != nil {
		return m.Gyroscope
	}
	return nil
}

func (m *OrientationReading) GetMagnetometer() *MagnetometerReading {
	if m != nil {
		return m.Magnetometer
	}
	return nil
}

func (m *OrientationReading) GetAccelerometer() *AccelerometerReading {
	if m != nil {
		return m.Accelerometer
	}
	return nil
}

func (m *OrientationReading) GetAhrs() *AttitudeHeadingReference {
	if m != nil {
		return m.Ahrs
	}
	return nil
}

// degrees per second
type GyroscopeReading struct {
	XDps float32 `protobuf:"fixed32,1,opt,name=x_dps,json=xDps" json:"x_dps,omitempty"`
	YDps float32 `protobuf:"fixed32,2,opt,name=y_dps,json=yDps" json:"y_dps,omitempty"`
	ZDps float32 `protobuf:"fixed32,3,opt,name=z_dps,json=zDps" json:"z_dps,omitempty"`
}

func (m *GyroscopeReading) Reset()                    { *m = GyroscopeReading{} }
func (m *GyroscopeReading) String() string            { return proto.CompactTextString(m) }
func (*GyroscopeReading) ProtoMessage()               {}
func (*GyroscopeReading) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GyroscopeReading) GetXDps() float32 {
	if m != nil {
		return m.XDps
	}
	return 0
}

func (m *GyroscopeReading) GetYDps() float32 {
	if m != nil {
		return m.YDps
	}
	return 0
}

func (m *GyroscopeReading) GetZDps() float32 {
	if m != nil {
		return m.ZDps
	}
	return 0
}

// micro teslas
type MagnetometerReading struct {
	XUt float32 `protobuf:"fixed32,1,opt,name=x_ut,json=xUt" json:"x_ut,omitempty"`
	YUt float32 `protobuf:"fixed32,2,opt,name=y_ut,json=yUt" json:"y_ut,omitempty"`
	ZUt float32 `protobuf:"fixed32,3,opt,name=z_ut,json=zUt" json:"z_ut,omitempty"`
}

func (m *MagnetometerReading) Reset()                    { *m = MagnetometerReading{} }
func (m *MagnetometerReading) String() string            { return proto.CompactTextString(m) }
func (*MagnetometerReading) ProtoMessage()               {}
func (*MagnetometerReading) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MagnetometerReading) GetXUt() float32 {
	if m != nil {
		return m.XUt
	}
	return 0
}

func (m *MagnetometerReading) GetYUt() float32 {
	if m != nil {
		return m.YUt
	}
	return 0
}

func (m *MagnetometerReading) GetZUt() float32 {
	if m != nil {
		return m.ZUt
	}
	return 0
}

// milli earth-gravities
type AccelerometerReading struct {
	XMg float32 `protobuf:"fixed32,1,opt,name=x_mg,json=xMg" json:"x_mg,omitempty"`
	YMg float32 `protobuf:"fixed32,2,opt,name=y_mg,json=yMg" json:"y_mg,omitempty"`
	ZMg float32 `protobuf:"fixed32,3,opt,name=z_mg,json=zMg" json:"z_mg,omitempty"`
}

func (m *AccelerometerReading) Reset()                    { *m = AccelerometerReading{} }
func (m *AccelerometerReading) String() string            { return proto.CompactTextString(m) }
func (*AccelerometerReading) ProtoMessage()               {}
func (*AccelerometerReading) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AccelerometerReading) GetXMg() float32 {
	if m != nil {
		return m.XMg
	}
	return 0
}

func (m *AccelerometerReading) GetYMg() float32 {
	if m != nil {
		return m.YMg
	}
	return 0
}

func (m *AccelerometerReading) GetZMg() float32 {
	if m != nil {
		return m.ZMg
	}
	return 0
}

type AttitudeHeadingReference struct {
	RollRad  float32 `protobuf:"fixed32,1,opt,name=roll_rad,json=rollRad" json:"roll_rad,omitempty"`
	PitchRad float32 `protobuf:"fixed32,2,opt,name=pitch_rad,json=pitchRad" json:"pitch_rad,omitempty"`
	YawRaw   float32 `protobuf:"fixed32,3,opt,name=yaw_raw,json=yawRaw" json:"yaw_raw,omitempty"`
}

func (m *AttitudeHeadingReference) Reset()                    { *m = AttitudeHeadingReference{} }
func (m *AttitudeHeadingReference) String() string            { return proto.CompactTextString(m) }
func (*AttitudeHeadingReference) ProtoMessage()               {}
func (*AttitudeHeadingReference) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *AttitudeHeadingReference) GetRollRad() float32 {
	if m != nil {
		return m.RollRad
	}
	return 0
}

func (m *AttitudeHeadingReference) GetPitchRad() float32 {
	if m != nil {
		return m.PitchRad
	}
	return 0
}

func (m *AttitudeHeadingReference) GetYawRaw() float32 {
	if m != nil {
		return m.YawRaw
	}
	return 0
}

func init() {
	proto.RegisterType((*OrientationReadings)(nil), "com.siliconbeacon.iot.messages.OrientationReadings")
	proto.RegisterType((*OrientationReading)(nil), "com.siliconbeacon.iot.messages.OrientationReading")
	proto.RegisterType((*GyroscopeReading)(nil), "com.siliconbeacon.iot.messages.GyroscopeReading")
	proto.RegisterType((*MagnetometerReading)(nil), "com.siliconbeacon.iot.messages.MagnetometerReading")
	proto.RegisterType((*AccelerometerReading)(nil), "com.siliconbeacon.iot.messages.AccelerometerReading")
	proto.RegisterType((*AttitudeHeadingReference)(nil), "com.siliconbeacon.iot.messages.AttitudeHeadingReference")
}

func init() { proto.RegisterFile("orientation.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 508 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd3, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x06, 0x60, 0xb5, 0xe9, 0xba, 0xf6, 0x8c, 0xa1, 0xcd, 0x45, 0x10, 0x8a, 0x04, 0x55, 0x6e,
	0xe8, 0x55, 0x86, 0x3a, 0x24, 0xb8, 0x65, 0x9a, 0x04, 0x48, 0xb4, 0x13, 0x81, 0x0a, 0x69, 0x37,
	0x91, 0x9b, 0x9c, 0x65, 0x16, 0x49, 0x1c, 0xd9, 0xee, 0xda, 0xe4, 0x2d, 0x78, 0x0d, 0x1e, 0x80,
	0xe7, 0x43, 0xb6, 0x93, 0x6e, 0x63, 0x63, 0xbd, 0xf4, 0x6f, 0xfb, 0x3b, 0x27, 0x3e, 0x0a, 0x1c,
	0x72, 0xc1, 0x30, 0x57, 0x54, 0x31, 0x9e, 0xfb, 0x85, 0xe0, 0x8a, 0x93, 0x97, 0x11, 0xcf, 0x7c,
	0xc9, 0x52, 0x16, 0xf1, 0x7c, 0x81, 0x34, 0xe2, 0xb9, 0xcf, 0xb8, 0xf2, 0x33, 0x94, 0x92, 0x26,
	0x28, 0x87, 0xaf, 0x12, 0xce, 0x93, 0x14, 0x8f, 0xcc, 0xe9, 0xc5, 0xf2, 0xe2, 0x48, 0xb1, 0x0c,
	0xa5, 0xa2, 0x59, 0x61, 0x01, 0xef, 0x4f, 0x0b, 0x06, 0x67, 0xd7, 0x6c, 0x80, 0x34, 0x66, 0x79,
	0x22, 0xc9, 0x53, 0xe8, 0xc6, 0x78, 0xc5, 0x22, 0x74, 0x5b, 0xa3, 0xd6, 0xb8, 0x1f, 0xd4, 0x2b,
	0xf2, 0x0e, 0xfa, 0x0b, 0x2a, 0x31, 0xd4, 0x8e, 0xdb, 0x1e, 0xb5, 0xc6, 0x7b, 0x93, 0xa1, 0x6f,
	0x8b, 0xf8, 0x4d, 0x11, 0xff, 0x7b, 0x53, 0x24, 0xe8, 0xe9, 0xc3, 0x7a, 0x49, 0x66, 0xd0, 0x13,
	0x35, 0xee, 0x3a, 0x23, 0x67, 0xbc, 0x37, 0x99, 0xf8, 0x0f, 0x37, 0xef, 0xdf, 0xed, 0x2b, 0xd8,
	0x18, 0xde, 0x2f, 0x07, 0xc8, 0xdd, 0x03, 0x64, 0x0c, 0x07, 0x02, 0x53, 0xaa, 0xd8, 0x95, 0xed,
	0x31, 0x5c, 0x4a, 0xf3, 0x05, 0xfb, 0xc1, 0xe3, 0x26, 0xd7, 0xed, 0xcc, 0x25, 0x99, 0x41, 0x3f,
	0x29, 0x05, 0x97, 0x11, 0x2f, 0x9a, 0x2f, 0x79, 0xb3, 0xad, 0xa3, 0x8f, 0xcd, 0x85, 0xa6, 0x9f,
	0x6b, 0x82, 0xfc, 0x80, 0x47, 0x19, 0x4d, 0x72, 0x54, 0x3c, 0x43, 0x85, 0xc2, 0x75, 0x0c, 0x79,
	0xbc, 0x8d, 0x9c, 0xde, 0xb8, 0xd3, 0xa8, 0xb7, 0x20, 0x72, 0x0e, 0xfb, 0x34, 0x8a, 0x30, 0x45,
	0x51, 0xcb, 0x1d, 0x23, 0xbf, 0xdd, 0x26, 0x7f, 0xb8, 0x79, 0xa9, 0xa1, 0x6f, 0x53, 0xe4, 0x0b,
	0x74, 0xe8, 0xa5, 0x90, 0xee, 0x8e, 0x21, 0xdf, 0x6f, 0x25, 0x95, 0x62, 0x6a, 0x19, 0xe3, 0xa7,
	0x5a, 0xc3, 0x0b, 0x14, 0x98, 0x47, 0x18, 0x18, 0xc5, 0xfb, 0x0a, 0x07, 0xff, 0xbe, 0x10, 0x19,
	0xc0, 0xce, 0x3a, 0x8c, 0x0b, 0x3b, 0x85, 0x76, 0xd0, 0x59, 0x9f, 0x16, 0x52, 0x87, 0xa5, 0x09,
	0xdb, 0x36, 0x2c, 0xeb, 0xb0, 0x32, 0xa1, 0x63, 0xc3, 0xea, 0xb4, 0x90, 0xde, 0x0c, 0x06, 0xf7,
	0xbc, 0x10, 0x39, 0x84, 0xce, 0x3a, 0x5c, 0xaa, 0x1a, 0x75, 0xd6, 0x73, 0xa5, 0xa3, 0x52, 0x47,
	0x96, 0x74, 0x4a, 0x1b, 0x55, 0x3a, 0xb2, 0xa0, 0x53, 0xcd, 0x95, 0x77, 0x06, 0x4f, 0xee, 0x7b,
	0x17, 0x0b, 0x66, 0xc9, 0x06, 0x9c, 0x26, 0x16, 0xcc, 0x92, 0x0d, 0x68, 0xa3, 0x4a, 0x47, 0x0d,
	0x38, 0x4d, 0xbc, 0x9f, 0xe0, 0xfe, 0xef, 0x55, 0xc8, 0x73, 0xe8, 0x09, 0x9e, 0xa6, 0xa1, 0xa0,
	0x71, 0x0d, 0xef, 0xea, 0x75, 0x40, 0x63, 0xf2, 0x02, 0xfa, 0x05, 0x53, 0xd1, 0xa5, 0xd9, 0xb3,
	0x15, 0x7a, 0x26, 0xd0, 0x9b, 0xcf, 0x60, 0xb7, 0xa4, 0xab, 0x50, 0xd0, 0x55, 0x5d, 0xa9, 0x5b,
	0xd2, 0x55, 0x40, 0x57, 0x27, 0xaf, 0xcf, 0x7b, 0xcd, 0x2c, 0x7e, 0xb7, 0x87, 0xdf, 0xec, 0xa0,
	0x4e, 0xec, 0xa0, 0x3e, 0x73, 0xe5, 0x4f, 0xeb, 0xcd, 0x45, 0xd7, 0xfc, 0x8b, 0xc7, 0x7f, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x64, 0xbc, 0x90, 0x53, 0x33, 0x04, 0x00, 0x00,
}
