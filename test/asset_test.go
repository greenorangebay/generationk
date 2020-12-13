// +build test
package generationk

import (
	genk "generationk/internal"
	"testing"
)

func TestAsset(t *testing.T) {
	ctx := genk.NewContext()
	dm := genk.NewCSVDataManager(ctx)

	abb := dm.ReadCSVFile("ABB.csv")
	c.AddAsset(&abb)

	eric := dm.ReadCSVFile("ABB.csv")
	c.AddAsset(&eric)

	want := 2
	v := len(c.Asset)

	if got := v; got != want {
		t.Errorf("ReadCSVFile(\"ABB.csv\") = %d, want %d", got, want)
	}

}