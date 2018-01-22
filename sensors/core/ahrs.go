package core

import (
	"math"
	"time"
)

type AhrsReading struct {
	Timestamp time.Time
	PitchRad  float32
	YawRad    float32
	RollRad   float32
}

const (
	degreesPerSecToRadiansPerSec = float32(0.0174533)
)

type MadgwickAhrs struct {
	beta           float32
	q0             float32
	q1             float32
	q2             float32
	q3             float32
	invSampleFreq  float32
	anglesComputed bool
}

func NewMadgwickAhrs(dataRate DataRate) *MadgwickAhrs {
	return &MadgwickAhrs{
		beta:           0.1,
		q0:             1.0,
		q1:             0.0,
		q2:             0.0,
		q3:             0.0,
		invSampleFreq:  float32(dataRate.Period / time.Second),
		anglesComputed: false,
	}
}

func (d *MadgwickAhrs) ComputeAhrs(gyro *GyroReading, accel *AccelReading, mag *MagReading) *AhrsReading {
	// convert gryo to radians / sec
	gx := gyro.Xdps * degreesPerSecToRadiansPerSec
	gy := gyro.Ydps * degreesPerSecToRadiansPerSec
	gz := gyro.Zdps * degreesPerSecToRadiansPerSec

	// rate of change of quaternion from gyro
	qDot1 := 0.5 * (-d.q1*gx - d.q2*gy - d.q3*gz)
	qDot2 := 0.5 * (d.q0*gx + d.q2*gz - d.q3*gy)
	qDot3 := 0.5 * (d.q0*gy - d.q1*gz + d.q3*gx)
	qDot4 := 0.5 * (d.q0*gz + d.q1*gy - d.q2*gx)

	// normalize accelerometer measurement
	recipNorm := float32(1.0 / math.Sqrt(float64(accel.Xmg*accel.Xmg+accel.Ymg*accel.Ymg+accel.Zmg*accel.Zmg)))
	ax := accel.Xmg * recipNorm
	ay := accel.Ymg * recipNorm
	az := accel.Zmg * recipNorm

	// normalize magnetometer measurement
	recipNorm = float32(1.0 / math.Sqrt(float64(mag.XuT*mag.XuT+mag.YuT*mag.YuT+mag.ZuT*mag.ZuT)))
	mx := mag.XuT * recipNorm
	my := mag.YuT * recipNorm
	mz := mag.ZuT * recipNorm

	// Auxiliary variables to avoid repeated arithmetic
	_2q0mx := 2.0 * d.q0 * mx
	_2q0my := 2.0 * d.q0 * my
	_2q0mz := 2.0 * d.q0 * mz
	_2q1mx := 2.0 * d.q1 * mx
	_2q0 := 2.0 * d.q0
	_2q1 := 2.0 * d.q1
	_2q2 := 2.0 * d.q2
	_2q3 := 2.0 * d.q3
	_2q0q2 := 2.0 * d.q0 * d.q2
	_2q2q3 := 2.0 * d.q2 * d.q3
	q0q0 := d.q0 * d.q0
	q0q1 := d.q0 * d.q1
	q0q2 := d.q0 * d.q2
	q0q3 := d.q0 * d.q3
	q1q1 := d.q1 * d.q1
	q1q2 := d.q1 * d.q2
	q1q3 := d.q1 * d.q3
	q2q2 := d.q2 * d.q2
	q2q3 := d.q2 * d.q3
	q3q3 := d.q3 * d.q3

	// Reference direction of Earth's magnetic field
	hx := mx*q0q0 - _2q0my*d.q3 + _2q0mz*d.q2 + mx*q1q1 + _2q1*my*d.q2 + _2q1*mz*d.q3 - mx*q2q2 - mx*q3q3
	hy := _2q0mx*d.q3 + my*q0q0 - _2q0mz*d.q1 + _2q1mx*d.q2 - my*q1q1 + my*q2q2 + _2q2*mz*d.q3 - my*q3q3
	_2bx := float32(math.Sqrt(float64(hx*hx + hy*hy)))
	_2bz := -_2q0mx*d.q2 + _2q0my*d.q1 + mz*q0q0 + _2q1mx*d.q3 - mz*q1q1 + _2q2*my*d.q3 - mz*q2q2 + mz*q3q3
	_4bx := 2.0 * _2bx
	_4bz := 2.0 * _2bz

	// Gradient decent algorithm corrective step
	s0 := -_2q2*(2.0*q1q3-_2q0q2-ax) + _2q1*(2.0*q0q1+_2q2q3-ay) - _2bz*d.q2*(_2bx*(0.5-q2q2-q3q3)+_2bz*(q1q3-q0q2)-mx) + (-_2bx*d.q3+_2bz*d.q1)*(_2bx*(q1q2-q0q3)+_2bz*(q0q1+q2q3)-my) + _2bx*d.q2*(_2bx*(q0q2+q1q3)+_2bz*(0.5-q1q1-q2q2)-mz)
	s1 := _2q3*(2.0*q1q3-_2q0q2-ax) + _2q0*(2.0*q0q1+_2q2q3-ay) - 4.0*d.q1*(1-2.0*q1q1-2.0*q2q2-az) + _2bz*d.q3*(_2bx*(0.5-q2q2-q3q3)+_2bz*(q1q3-q0q2)-mx) + (_2bx*d.q2+_2bz*d.q0)*(_2bx*(q1q2-q0q3)+_2bz*(q0q1+q2q3)-my) + (_2bx*d.q3-_4bz*d.q1)*(_2bx*(q0q2+q1q3)+_2bz*(0.5-q1q1-q2q2)-mz)
	s2 := -_2q0*(2.0*q1q3-_2q0q2-ax) + _2q3*(2.0*q0q1+_2q2q3-ay) - 4.0*d.q2*(1-2.0*q1q1-2.0*q2q2-az) + (-_4bx*d.q2-_2bz*d.q0)*(_2bx*(0.5-q2q2-q3q3)+_2bz*(q1q3-q0q2)-mx) + (_2bx*d.q1+_2bz*d.q3)*(_2bx*(q1q2-q0q3)+_2bz*(q0q1+q2q3)-my) + (_2bx*d.q0-_4bz*d.q2)*(_2bx*(q0q2+q1q3)+_2bz*(0.5-q1q1-q2q2)-mz)
	s3 := _2q1*(2.0*q1q3-_2q0q2-ax) + _2q2*(2.0*q0q1+_2q2q3-ay) + (-_4bx*d.q3+_2bz*d.q1)*(_2bx*(0.5-q2q2-q3q3)+_2bz*(q1q3-q0q2)-mx) + (-_2bx*d.q0+_2bz*d.q2)*(_2bx*(q1q2-q0q3)+_2bz*(q0q1+q2q3)-my) + _2bx*d.q1*(_2bx*(q0q2+q1q3)+_2bz*(0.5-q1q1-q2q2)-mz)

	// normalise step magnitude
	recipNorm = float32(1.0 / math.Sqrt(float64(s0*s0+s1*s1+s2*s2+s3*s3)))
	s0 *= recipNorm
	s1 *= recipNorm
	s2 *= recipNorm
	s3 *= recipNorm

	// Apply feedback step
	qDot1 -= d.beta * s0
	qDot2 -= d.beta * s1
	qDot3 -= d.beta * s2
	qDot4 -= d.beta * s3

	// Integrate rate of change of quaternion to yield quaternion
	d.q0 += qDot1 * d.invSampleFreq
	d.q1 += qDot2 * d.invSampleFreq
	d.q2 += qDot3 * d.invSampleFreq
	d.q3 += qDot4 * d.invSampleFreq

	// Normalise quaternion
	recipNorm = float32(1.0 / math.Sqrt(float64(d.q0*d.q0+d.q1*d.q1+d.q2*d.q2+d.q3*d.q3)))
	d.q0 *= recipNorm
	d.q1 *= recipNorm
	d.q2 *= recipNorm
	d.q3 *= recipNorm

	return &AhrsReading{
		Timestamp: gyro.Timestamp,
		RollRad:   float32(math.Atan2(float64(d.q0*d.q1+d.q2*d.q3), float64(0.5-d.q1*d.q1-d.q2*d.q2))),
		PitchRad:  float32(math.Asin(float64(-2.0 * (d.q1*d.q3 - d.q0*d.q2)))),
		YawRad:    float32(math.Atan2(float64(d.q1*d.q2+d.q0*d.q3), float64(0.5-d.q2*d.q2-d.q3*d.q3))),
	}
}
