package internal

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Asset data type
type Asset struct {
	Name string
	Ohlc []OHLC
}

//Type is the type for ohlc
type Type struct {
	Open  string
	High  string
	Low   string
	close string
}

// OHLC data type
type OHLC struct {
	Time                   time.Time
	Open, High, Low, Close float64
	Volume                 int
}

type EndOfDataError struct {
	Description string
}

func (e *EndOfDataError) Error() string {
	return fmt.Sprintf("End of data: %s", e.Description)
}

type DataNotInCombatZone struct {
	Description string
}

func (e *DataNotInCombatZone) Error() string {
	return fmt.Sprintf("DataNotInCombatZone: %s", e.Description)
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	//fmt.Printf("date1 %v", date1)
	y2, m2, d2 := date2.Date()
	//fmt.Printf("date2 %v", date2)
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (a *Asset) Shift(time time.Time) (int, error) {
	//fmt.Printf("Shifting data\n")
	//fmt.Printf("Shifting! Last date %v\n", a.Ohlc[len(a.Ohlc)-1].Time)
	var i int
	for ok := true; ok; ok = a.Ohlc[0].Time.Before(time) && len(a.Ohlc) > 0 {
		a.Ohlc = a.Ohlc[1:]
		i++
	}

	fmt.Printf("Len: a.Ohlc %v\n", len(a.Ohlc))
	fmt.Printf("New value: %f\n", a.Ohlc[0].Close)

	return i, nil
}

//CloseArray is used to get the close series
func (a *Asset) CloseArray() []float64 {
	s := make([]float64, len(a.Ohlc))

	for i, ohlc := range a.Ohlc {
		s[i] = ohlc.Close
	}
	log.WithFields(log.Fields{
		"CloseArray() length": len(s),
	}).Debug("Asset: CloseArray()")
	return s
}

//Close is used to get the close value
func (a *Asset) Close() float64 {
	return a.Ohlc[0].Close
}

//CloseAtBar is used to get the close value
func (a *Asset) CloseAtBar(ix int) float64 {
	return a.Ohlc[ix].Close
}

// Portfolio structure
type Portfolio struct {
}