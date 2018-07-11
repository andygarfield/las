package las

import (
	"fmt"
	"testing"
)

func TestOpenFile(t *testing.T) {
	_, err := OpenFile("./testdata/FL_PinellasCo_2007_000398.las")
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}
}

func TestReadPoint(t *testing.T) {
	las, err := OpenFile("./testdata/FL_PinellasCo_2007_000398.las")
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}

	p, err := las.ReadPoint()
	if err != nil {
		fmt.Println(err)
		t.Failed()
	}

	fmt.Printf("%f\n", p.X)
	fmt.Printf("%f\n", p.Y)
	fmt.Printf("%f\n", p.Z)
	fmt.Println(p.Intensity)
	fmt.Println(p.Return)
	fmt.Println(p.NumOfReturns)
	fmt.Println(p.ScanDir)
	fmt.Println(p.EdgeOfFlightLine)
	fmt.Println(p.Classification)
	fmt.Println(p.ScanAngleRank)
	fmt.Println(p.UserData)
	fmt.Println(p.PointSourceID)
}
