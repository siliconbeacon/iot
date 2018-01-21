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
	YawRad   float32 `protobuf:"fixed32,3,opt,name=yaw_rad,json=yawRad" json:"yaw_rad,omitempty"`
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

func (m *AttitudeHeadingReference) GetYawRad() float32 {
	if m != nil {
		return m.YawRad
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
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0xd3, 0xd1, 0x6e, 0xd3, 0x3c,
	0x14, 0x07, 0x70, 0xb5, 0xe9, 0xba, 0xf6, 0xec, 0xdb, 0xa7, 0xcd, 0x45, 0x10, 0x8a, 0x04, 0x55,
	0x6e, 0xe8, 0x55, 0x86, 0x3a, 0x24, 0xb8, 0x65, 0x9a, 0x04, 0x48, 0xb4, 0x13, 0x81, 0x0a, 0x69,
	0x37, 0x95, 0x9b, 0x9c, 0x65, 0x16, 0x49, 0x1c, 0xd9, 0xee, 0xd6, 0xe4, 0x2d, 0x78, 0x0d, 0x1e,
	0x80, 0xe7, 0x43, 0xb6, 0xe3, 0x6e, 0x63, 0x83, 0x5e, 0xfa, 0x6f, 0xfb, 0x77, 0x4e, 0x7c, 0x14,
	0x38, 0xe4, 0x82, 0x61, 0xa1, 0xa8, 0x62, 0xbc, 0x08, 0x4b, 0xc1, 0x15, 0x27, 0xcf, 0x63, 0x9e,
	0x87, 0x92, 0x65, 0x2c, 0xe6, 0xc5, 0x12, 0x69, 0xcc, 0x8b, 0x90, 0x71, 0x15, 0xe6, 0x28, 0x25,
	0x4d, 0x51, 0x0e, 0x5f, 0xa4, 0x9c, 0xa7, 0x19, 0x1e, 0x99, 0xd3, 0xcb, 0xd5, 0xc5, 0x91, 0x62,
	0x39, 0x4a, 0x45, 0xf3, 0xd2, 0x02, 0xc1, 0xaf, 0x16, 0x0c, 0xce, 0x6e, 0xd8, 0x08, 0x69, 0xc2,
	0x8a, 0x54, 0x92, 0xc7, 0xd0, 0x4d, 0xf0, 0x8a, 0xc5, 0xe8, 0xb7, 0x46, 0xad, 0x71, 0x3f, 0x6a,
	0x56, 0xe4, 0x0d, 0xf4, 0x97, 0x54, 0xe2, 0x42, 0x3b, 0x7e, 0x7b, 0xd4, 0x1a, 0xef, 0x4d, 0x86,
	0xa1, 0x2d, 0x12, 0xba, 0x22, 0xe1, 0x57, 0x57, 0x24, 0xea, 0xe9, 0xc3, 0x7a, 0x49, 0x66, 0xd0,
	0x13, 0x0d, 0xee, 0x7b, 0x23, 0x6f, 0xbc, 0x37, 0x99, 0x84, 0xff, 0x6e, 0x3e, 0xbc, 0xdf, 0x57,
	0xb4, 0x31, 0x82, 0x1f, 0x1e, 0x90, 0xfb, 0x07, 0xc8, 0x18, 0x0e, 0x04, 0x66, 0x54, 0xb1, 0x2b,
	0xdb, 0xe3, 0x62, 0x25, 0xcd, 0x17, 0xec, 0x47, 0xff, 0xbb, 0x5c, 0xb7, 0x33, 0x97, 0x64, 0x06,
	0xfd, 0xb4, 0x12, 0x5c, 0xc6, 0xbc, 0x74, 0x5f, 0xf2, 0x6a, 0x5b, 0x47, 0xef, 0xdd, 0x05, 0xd7,
	0xcf, 0x0d, 0x41, 0xbe, 0xc1, 0x7f, 0x39, 0x4d, 0x0b, 0x54, 0x3c, 0x47, 0x85, 0xc2, 0xf7, 0x0c,
	0x79, 0xbc, 0x8d, 0x9c, 0xde, 0xba, 0xe3, 0xd4, 0x3b, 0x10, 0x39, 0x87, 0x7d, 0x1a, 0xc7, 0x98,
	0xa1, 0x68, 0xe4, 0x8e, 0x91, 0x5f, 0x6f, 0x93, 0xdf, 0xdd, 0xbe, 0xe4, 0xe8, 0xbb, 0x14, 0xf9,
	0x04, 0x1d, 0x7a, 0x29, 0xa4, 0xbf, 0x63, 0xc8, 0xb7, 0x5b, 0x49, 0xa5, 0x98, 0x5a, 0x25, 0xf8,
	0xa1, 0xd1, 0xf0, 0x02, 0x05, 0x16, 0x31, 0x46, 0x46, 0x09, 0x3e, 0xc3, 0xc1, 0x9f, 0x2f, 0x44,
	0x06, 0xb0, 0xb3, 0x5e, 0x24, 0xa5, 0x9d, 0x42, 0x3b, 0xea, 0xac, 0x4f, 0x4b, 0xa9, 0xc3, 0xca,
	0x84, 0x6d, 0x1b, 0x56, 0x4d, 0x58, 0x9b, 0xd0, 0xb3, 0x61, 0x7d, 0x5a, 0xca, 0x60, 0x06, 0x83,
	0x07, 0x5e, 0x88, 0x1c, 0x42, 0x67, 0xbd, 0x58, 0xa9, 0x06, 0xf5, 0xd6, 0x73, 0xa5, 0xa3, 0x4a,
	0x47, 0x96, 0xf4, 0x2a, 0x1b, 0xd5, 0x3a, 0xb2, 0xa0, 0x57, 0xcf, 0x55, 0x70, 0x06, 0x8f, 0x1e,
	0x7a, 0x17, 0x0b, 0xe6, 0xe9, 0x06, 0x9c, 0xa6, 0x16, 0xcc, 0xd3, 0x0d, 0x68, 0xa3, 0x5a, 0x47,
	0x0e, 0x9c, 0xa6, 0xc1, 0x77, 0xf0, 0xff, 0xf6, 0x2a, 0xe4, 0x29, 0xf4, 0x04, 0xcf, 0xb2, 0x85,
	0xa0, 0x49, 0x03, 0xef, 0xea, 0x75, 0x44, 0x13, 0xf2, 0x0c, 0xfa, 0x25, 0x53, 0xf1, 0xa5, 0xd9,
	0xb3, 0x15, 0x7a, 0x26, 0xd0, 0x9b, 0x4f, 0x60, 0xb7, 0xa2, 0xd7, 0x66, 0xcb, 0x56, 0xea, 0x56,
	0xf4, 0x3a, 0xa2, 0xc9, 0xc9, 0xcb, 0xf3, 0x9e, 0x9b, 0xc5, 0xcf, 0xf6, 0xf0, 0x8b, 0x1d, 0xd4,
	0x89, 0x1d, 0xd4, 0x47, 0xae, 0xc2, 0x69, 0xb3, 0xb9, 0xec, 0x9a, 0x7f, 0xf1, 0xf8, 0x77, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x88, 0xa6, 0xdf, 0xfe, 0x33, 0x04, 0x00, 0x00,
}
