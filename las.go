package las

import (
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
	PointSourceId                                   uint16
}

// LAS reads a LAS file
type Reader struct {
	Lasfile                                           *os.File
	Version, SystemIdentifier, GeneratingSoftware     string
	Day, Year, HeaderSize                             uint16
	OffsetToPointData, VLRNum                         uint32
	PointFormat                                       uint8
	PointRecLen                                       uint16
	PointNum                                          uint32
	XScale, YScale, ZScale, XOffset, YOffset, ZOffset float64
	XMax, XMin, YMax, YMin, ZMax, ZMin                float64
}

func Open(filename string) (*Reader, error) {
	f, err := os.Open(filename)

	if err != nil {
		return &Reader{}, err
	}

	r := &Reader{
		Lasfile: f,
	}

	r.readHeader()

	return r, nil
}
func (r *Reader) readHeader() error {
	skipBytes(r.Lasfile, 24)

	maj := strconv.Itoa(int(readUInt8(r.Lasfile)))
	min := strconv.Itoa(int(readUInt8(r.Lasfile)))
	r.Version = fmt.Sprintf("%s.%s", maj, min)

	r.SystemIdentifier = readString(r.Lasfile, 32)
	r.GeneratingSoftware = readString(r.Lasfile, 32)

	r.Day = readUInt16(r.Lasfile)
	r.Year = readUInt16(r.Lasfile)

	r.HeaderSize = readUInt16(r.Lasfile)
	r.OffsetToPointData = readUInt32(r.Lasfile)

	r.VLRNum = readUInt32(r.Lasfile)

	r.PointFormat = readUInt8(r.Lasfile)
	r.PointRecLen = readUInt16(r.Lasfile)
	r.PointNum = readUInt32(r.Lasfile)
	skipBytes(r.Lasfile, 20)

	r.XScale = readFloat64(r.Lasfile)
	r.YScale = readFloat64(r.Lasfile)
	r.ZScale = readFloat64(r.Lasfile)
	r.XOffset = readFloat64(r.Lasfile)
	r.YOffset = readFloat64(r.Lasfile)
	r.ZOffset = readFloat64(r.Lasfile)

	r.XMax = readFloat64(r.Lasfile)
	r.XMin = readFloat64(r.Lasfile)
	r.YMax = readFloat64(r.Lasfile)
	r.YMin = readFloat64(r.Lasfile)
	r.ZMax = readFloat64(r.Lasfile)
	r.ZMin = readFloat64(r.Lasfile)

	remainingToFirstPoint := r.OffsetToPointData - 227
	skipBytes(r.Lasfile, int(remainingToFirstPoint))

	return nil
}

func (r Reader) ReadPoint() LidarPoint {
	lp := LidarPoint{}

	lp.X = (float64(readUInt32(r.Lasfile)) * r.XScale) + r.XOffset
	lp.Y = (float64(readUInt32(r.Lasfile)) * r.YScale) + r.YOffset
	lp.Z = (float64(readUInt32(r.Lasfile)) * r.ZScale) + r.ZOffset
	lp.Intensity = readUInt16(r.Lasfile)

	returnB := readByte(r.Lasfile)
	lp.Return = returnB >> 5
	lp.NumOfReturns = returnB << 3 >> 5
	lp.ScanDir = returnB << 6 >> 7
	lp.EdgeOfFlightLine = returnB << 7 >> 7

	lp.Classification = readUInt8(r.Lasfile)
	lp.ScanAngleRank = readInt8(r.Lasfile)
	lp.UserData = readUInt8(r.Lasfile)
	lp.PointSourceId = readUInt16(r.Lasfile)

	skipBytes(r.Lasfile, 12)

	return lp
}

func readByte(r io.Reader) byte {
	b := make([]byte, 1)
	r.Read(b)
	return b[0]
}

func readString(r io.Reader, l int) string {
	b := make([]byte, l)
	r.Read(b)
	return strings.TrimSpace(string(b))
}

func readInt8(r io.Reader) int8 {
	b := make([]byte, 1)
	r.Read(b)
	return int8(b[0])
}

func readUInt8(r io.Reader) uint8 {
	b := make([]byte, 1)
	r.Read(b)
	return uint8(b[0])
}

func readUInt16(r io.Reader) uint16 {
	b := make([]byte, 2)
	r.Read(b)
	return binary.LittleEndian.Uint16(b)
}

func readUInt32(r io.Reader) uint32 {
	b := make([]byte, 4)
	r.Read(b)
	return binary.LittleEndian.Uint32(b)
}

func readFloat64(r io.Reader) float64 {
	var f float64
	binary.Read(r, binary.LittleEndian, &f)
	return f
}

func skipBytes(r io.Reader, l int) {
	b := make([]byte, l)
	r.Read(b)
}
