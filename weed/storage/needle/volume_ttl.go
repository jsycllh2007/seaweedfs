package needle

import (
	"strconv"
)

const (
	//stored unit types
	Empty byte = iota
	Minute
	Hour
	Day
	Week
	Month
	Year
)

type TTL struct {
	Count byte
	Unit  byte
}

var EMPTY_TTL = &TTL{}

// translate a readable ttl to internal ttl
// Supports format example:
// 3m: 3 minutes
// 4h: 4 hours
// 5d: 5 days
// 6w: 6 weeks
// 7M: 7 months
// 8y: 8 years
func ReadTTL(ttlString string) (*TTL, error) {
	if ttlString == "" {
		return EMPTY_TTL, nil
	}
	ttlBytes := []byte(ttlString)
	unitByte := ttlBytes[len(ttlBytes)-1]
	countBytes := ttlBytes[0 : len(ttlBytes)-1]
	if '0' <= unitByte && unitByte <= '9' {
		countBytes = ttlBytes
		unitByte = 'm'
	}
	count, err := strconv.Atoi(string(countBytes))
	unit := toStoredByte(unitByte)
	return &TTL{Count: byte(count), Unit: unit}, err
}

// read stored bytes to a ttl
func LoadTTLFromBytes(input []byte) (t *TTL) {
	return &TTL{Count: input[0], Unit: input[1]}
}

// read stored bytes to a ttl
func LoadTTLFromUint32(ttl uint32) (t *TTL) {
	input := make([]byte, 2)
	input[1] = byte(ttl)
	input[0] = byte(ttl >> 8)
	return LoadTTLFromBytes(input)
}

// save stored bytes to an output with 2 bytes
func (t *TTL) ToBytes(output []byte) {
	output[0] = t.Count
	output[1] = t.Unit
}

func (t *TTL) ToUint32() (output uint32) {
	output = uint32(t.Count) << 8
	output += uint32(t.Unit)
	return output
}

func (t *TTL) String() string {
	if t == nil || t.Count == 0 {
		return ""
	}
	if t.Unit == Empty {
		return ""
	}
	countString := strconv.Itoa(int(t.Count))
	switch t.Unit {
	case Minute:
		return countString + "m"
	case Hour:
		return countString + "h"
	case Day:
		return countString + "d"
	case Week:
		return countString + "w"
	case Month:
		return countString + "M"
	case Year:
		return countString + "y"
	}
	return ""
}

func toStoredByte(readableUnitByte byte) byte {
	switch readableUnitByte {
	case 'm':
		return Minute
	case 'h':
		return Hour
	case 'd':
		return Day
	case 'w':
		return Week
	case 'M':
		return Month
	case 'y':
		return Year
	}
	return 0
}

func (t TTL) Minutes() uint32 {
	switch t.Unit {
	case Empty:
		return 0
	case Minute:
		return uint32(t.Count)
	case Hour:
		return uint32(t.Count) * 60
	case Day:
		return uint32(t.Count) * 60 * 24
	case Week:
		return uint32(t.Count) * 60 * 24 * 7
	case Month:
		return uint32(t.Count) * 60 * 24 * 31
	case Year:
		return uint32(t.Count) * 60 * 24 * 365
	}
	return 0
}
