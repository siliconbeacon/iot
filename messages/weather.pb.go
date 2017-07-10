// Code generated by protoc-gen-go. DO NOT EDIT.
// source: weather.proto

/*
Package messages is a generated protocol buffer package.

It is generated from these files:
	weather.proto

It has these top-level messages:
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

type WeatherReadings struct {
	Device   string                     `protobuf:"bytes,1,opt,name=device" json:"device,omitempty"`
	BaseTime *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=base_time,json=baseTime" json:"base_time,omitempty"`
	Readings []*WeatherReading          `protobuf:"bytes,3,rep,name=readings" json:"readings,omitempty"`
}

func (m *WeatherReadings) Reset()                    { *m = WeatherReadings{} }
func (m *WeatherReadings) String() string            { return proto.CompactTextString(m) }
func (*WeatherReadings) ProtoMessage()               {}
func (*WeatherReadings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *WeatherReadings) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *WeatherReadings) GetBaseTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.BaseTime
	}
	return nil
}

func (m *WeatherReadings) GetReadings() []*WeatherReading {
	if m != nil {
		return m.Readings
	}
	return nil
}

type WeatherReading struct {
	RelativeTimeUs             uint32  `protobuf:"varint,1,opt,name=relative_time_us,json=relativeTimeUs" json:"relative_time_us,omitempty"`
	TemperatureDegreesC        float64 `protobuf:"fixed64,2,opt,name=temperature_degrees_c,json=temperatureDegreesC" json:"temperature_degrees_c,omitempty"`
	RelativeHumidityPercentage float64 `protobuf:"fixed64,3,opt,name=relative_humidity_percentage,json=relativeHumidityPercentage" json:"relative_humidity_percentage,omitempty"`
}

func (m *WeatherReading) Reset()                    { *m = WeatherReading{} }
func (m *WeatherReading) String() string            { return proto.CompactTextString(m) }
func (*WeatherReading) ProtoMessage()               {}
func (*WeatherReading) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *WeatherReading) GetRelativeTimeUs() uint32 {
	if m != nil {
		return m.RelativeTimeUs
	}
	return 0
}

func (m *WeatherReading) GetTemperatureDegreesC() float64 {
	if m != nil {
		return m.TemperatureDegreesC
	}
	return 0
}

func (m *WeatherReading) GetRelativeHumidityPercentage() float64 {
	if m != nil {
		return m.RelativeHumidityPercentage
	}
	return 0
}

func init() {
	proto.RegisterType((*WeatherReadings)(nil), "com.siliconbeacon.iot.messages.WeatherReadings")
	proto.RegisterType((*WeatherReading)(nil), "com.siliconbeacon.iot.messages.WeatherReading")
}

func init() { proto.RegisterFile("weather.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4b, 0xfb, 0x30,
	0x18, 0xc6, 0xc9, 0x06, 0x63, 0xcb, 0xd8, 0xfe, 0x7f, 0x22, 0x4a, 0x29, 0xa2, 0x63, 0x17, 0x7b,
	0xca, 0x60, 0x1e, 0xbc, 0xca, 0xf4, 0xa0, 0x82, 0x20, 0x51, 0x11, 0xbc, 0x94, 0x34, 0x7d, 0xcd,
	0x02, 0x4b, 0x53, 0x92, 0x74, 0xe2, 0x57, 0xf2, 0x20, 0x7e, 0x44, 0x69, 0xd3, 0x0e, 0x77, 0xf1,
	0xf8, 0xe6, 0x79, 0xde, 0xe7, 0xfd, 0x3d, 0xc1, 0x93, 0x77, 0xe0, 0x7e, 0x0d, 0x96, 0x96, 0xd6,
	0x78, 0x43, 0x4e, 0x84, 0xd1, 0xd4, 0xa9, 0x8d, 0x12, 0xa6, 0xc8, 0x80, 0x0b, 0x53, 0x50, 0x65,
	0x3c, 0xd5, 0xe0, 0x1c, 0x97, 0xe0, 0xe2, 0x53, 0x69, 0x8c, 0xdc, 0xc0, 0xa2, 0x71, 0x67, 0xd5,
	0xdb, 0xc2, 0x2b, 0x0d, 0xce, 0x73, 0x5d, 0x86, 0x80, 0xf9, 0x17, 0xc2, 0xff, 0x5e, 0x42, 0x24,
	0x03, 0x9e, 0xab, 0x42, 0x3a, 0x72, 0x84, 0x07, 0x39, 0x6c, 0x95, 0x80, 0x08, 0xcd, 0x50, 0x32,
	0x62, 0xed, 0x44, 0x2e, 0xf0, 0x28, 0xe3, 0x0e, 0xd2, 0x3a, 0x23, 0xea, 0xcd, 0x50, 0x32, 0x5e,
	0xc6, 0x34, 0x1c, 0xa0, 0xdd, 0x01, 0xfa, 0xd4, 0x1d, 0x60, 0xc3, 0xda, 0x5c, 0x8f, 0xe4, 0x0e,
	0x0f, 0x6d, 0x1b, 0x1e, 0xf5, 0x67, 0xfd, 0x64, 0xbc, 0xa4, 0xf4, 0x6f, 0x70, 0xba, 0xcf, 0xc4,
	0x76, 0xfb, 0xf3, 0x6f, 0x84, 0xa7, 0xfb, 0x22, 0x49, 0xf0, 0x7f, 0x0b, 0x1b, 0xee, 0xd5, 0x36,
	0xb0, 0xa5, 0x95, 0x6b, 0xc8, 0x27, 0x6c, 0xda, 0xbd, 0xd7, 0x18, 0xcf, 0x8e, 0x2c, 0xf1, 0xa1,
	0x07, 0x5d, 0x82, 0xe5, 0xbe, 0xb2, 0x90, 0xe6, 0x20, 0x2d, 0x80, 0x4b, 0x45, 0xd3, 0x06, 0xb1,
	0x83, 0x5f, 0xe2, 0x75, 0xd0, 0xae, 0xc8, 0x25, 0x3e, 0xde, 0xa5, 0xaf, 0x2b, 0xad, 0x72, 0xe5,
	0x3f, 0xd2, 0x12, 0xac, 0x80, 0xc2, 0x73, 0x09, 0x51, 0xbf, 0x59, 0x8d, 0x3b, 0xcf, 0x4d, 0x6b,
	0x79, 0xd8, 0x39, 0x56, 0x67, 0xaf, 0xc3, 0xae, 0xd7, 0x67, 0x2f, 0x7e, 0x0c, 0xa5, 0x57, 0xa1,
	0xf4, 0xad, 0xf1, 0xf4, 0xbe, 0x15, 0xb3, 0x41, 0xf3, 0x8b, 0xe7, 0x3f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xee, 0xe8, 0x28, 0x02, 0xe5, 0x01, 0x00, 0x00,
}
