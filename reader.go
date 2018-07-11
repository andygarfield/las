package las

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type LidarPoint struct {
	X, Y, Z                                         float64
	Intensity                                       uint16
	Return, NumOfReturns, ScanDir, EdgeOfFlightLine byte
	Classification                                  byte
	ScanAngleRank                                   int8
	UserData                                        uint8
	PointSourceID                                   uint16
	GPSTime                                         float64
}

// Reader reads a LAS file
type Reader struct {
	r                                                 *bufio.Reader
	Version, SystemIdentifier, GeneratingSoftware     string
	Day, Year, HeaderSize                             uint16
	OffsetToPointData, VLRNum                         uint32
	PointFormat                                       uint8
	PointRecLen                                       uint16
	PointNum                                          uint32
	XScale, YScale, ZScale, XOffset, YOffset, ZOffset float64
	XMax, XMin, YMax, YMin, ZMax, ZMin                float64
}

func OpenFile(filename string) (*Reader, error) {
	f, err := os.Open(filename)

	if err != nil {
		return &Reader{}, err
	}

	r := &Reader{
		r: bufio.NewReader(f),
	}

	r.readHeader()

	return r, nil
}

func (r *Reader) readHeader() error {
	skipBytes(r.r, 24)

	majNum, err := readUInt8(r.r)
	minNum, err := readUInt8(r.r)

	maj := strconv.Itoa(int(majNum))
	min := strconv.Itoa(int(minNum))
	r.Version = fmt.Sprintf("%s.%s", maj, min)

	r.SystemIdentifier, err = readString(r.r, 32)
	r.GeneratingSoftware, err = readString(r.r, 32)

	r.Day, _ = readUInt16(r.r)
	r.Year, _ = readUInt16(r.r)

	r.HeaderSize, _ = readUInt16(r.r)
	r.OffsetToPointData, _ = readUInt32(r.r)

	r.VLRNum, _ = readUInt32(r.r)

	r.PointFormat, _ = readUInt8(r.r)
	r.PointRecLen, _ = readUInt16(r.r)
	r.PointNum, _ = readUInt32(r.r)
	skipBytes(r.r, 20)

	r.XScale, _ = readFloat64(r.r)
	r.YScale, _ = readFloat64(r.r)
	r.ZScale, _ = readFloat64(r.r)
	r.XOffset, _ = readFloat64(r.r)
	r.YOffset, _ = readFloat64(r.r)
	r.ZOffset, _ = readFloat64(r.r)

	r.XMax, _ = readFloat64(r.r)
	r.XMin, _ = readFloat64(r.r)
	r.YMax, _ = readFloat64(r.r)
	r.YMin, _ = readFloat64(r.r)
	r.ZMax, _ = readFloat64(r.r)
	r.ZMin, _ = readFloat64(r.r)

	remainingToFirstPoint := r.OffsetToPointData - 227
	skipBytes(r.r, int(remainingToFirstPoint))

	return err
}

func (r *Reader) ReadPoint() (LidarPoint, error) {
	lp := LidarPoint{}

	xi, err := readUInt32(r.r)
	yi, err := readUInt32(r.r)
	zi, err := readUInt32(r.r)

	lp.X = (float64(xi) * r.XScale) + r.XOffset
	lp.Y = (float64(yi) * r.YScale) + r.YOffset
	lp.Z = (float64(zi) * r.ZScale) + r.ZOffset
	lp.Intensity, err = readUInt16(r.r)

	returnB, err := readByte(r.r)
	lp.Return = returnB >> 5
	lp.NumOfReturns = returnB << 3 >> 5
	lp.ScanDir = returnB << 6 >> 7
	lp.EdgeOfFlightLine = returnB << 7 >> 7

	lp.Classification, err = readUInt8(r.r)
	lp.ScanAngleRank, err = readInt8(r.r)
	lp.UserData, err = readUInt8(r.r)
	lp.PointSourceID, err = readUInt16(r.r)

	if r.PointFormat == 1 {
		lp.GPSTime, err = readFloat64(r.r)
	}

	return lp, err
}

func readByte(r io.Reader) (byte, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	return b[0], err
}

func readString(r io.Reader, l int) (string, error) {
	b := make([]byte, l)
	_, err := r.Read(b)
	return strings.TrimSpace(string(b)), err
}

func readInt8(r io.Reader) (int8, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	return int8(b[0]), err
}

func readUInt8(r io.Reader) (uint8, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	return uint8(b[0]), err
}

func readUInt16(r io.Reader) (uint16, error) {
	b := make([]byte, 2)
	_, err := r.Read(b)
	return binary.LittleEndian.Uint16(b), err
}

func readUInt32(r io.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := r.Read(b)
	return binary.LittleEndian.Uint32(b), err
}

func readFloat64(r io.Reader) (float64, error) {
	var f float64
	err := binary.Read(r, binary.LittleEndian, &f)
	return f, err
}

func skipBytes(r io.Reader, l int) {
	b := make([]byte, l)
	r.Read(b)
}
