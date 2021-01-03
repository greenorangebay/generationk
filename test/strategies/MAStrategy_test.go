package generationk_test

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	K "github.com/0dayfall/generationk"
	indicators "github.com/0dayfall/generationk/indicators"
	"github.com/pkg/profile"

	log "github.com/sirupsen/logrus"
)

//Strategy strategy
type MACrossStrategy struct {
	ma50  indicators.SimpleMovingAverage
	close indicators.TimeSeries
}

//Setup is used to declare what indicators will be used
func (ma *MACrossStrategy) Once(ctx *K.Context) error {
	//Want access to the latest 5 closing prices
	ma.close = indicators.NewTimeSeries(indicators.Close, 5)
	//MA50
	ma.ma50 = indicators.NewSimpleMovingAverage(indicators.Close, 50)

	//Add indicators to context
	ctx.AddIndicator(&ma.close)
	ctx.AddIndicator(&ma.ma50)

	//The data needed to calculate MA
	ctx.SetInitPeriod(50)

	return nil
}

//Update gets called when updates arrive
func (ma *MACrossStrategy) Update(ctx *K.Context) {
	ctx.K++
}

//Tick get called when there is new data coming in
func (ma *MACrossStrategy) PerBar(ohlc K.OHLC, callback K.Callback) {

	if ma.close.Current() > ma.ma50.Current() {
		if !callback.IsOwning(callback.Assets()[0]) {
			err := callback.OrderSend(callback.Assets()[0], K.BuyOrder, K.MarketOrder, 0, 100)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if ma.close.Current() < ma.ma50.Current() {
		if callback.IsOwning(callback.Assets()[0]) {
			err := callback.OrderSend(callback.Assets()[0], K.SellOrder, K.MarketOrder, 0, 100)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

//OrderEvent gets called on order events
func (ma *MACrossStrategy) OrderEvent(orderEvent K.Event) {
	/*log.WithFields(log.Fields{
		"orderEvent": orderEvent,
	}).Info("MAStrategy_test> OrderEvent")*/
}

func readFolder(folderPath string) {

	files, err := filepath.Glob(folderPath + "*.csv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("files %s", files)
	//d.ReadCSVFilesAsync(files)
	portfolio := K.NewPortfolio()
	portfolio.SetBalance(100000)

	var wg sync.WaitGroup

	y := 0

	for _, fileName := range files {
		wg.Add(1)
		go func(localFilename string) {
			genk := K.NewGenerationK()
			genk.AddPortfolio(portfolio)

			strategy := new(MACrossStrategy)

			//Going to run with the following data thingie to collect the data
			//assetName := strings.TrimSuffix(filepath.Base(fileName), path.Ext(fileName))
			//genk.AddAsset(NewAsset(assetName, OHLC{}))
			//genk.AddAsset(NewAsset(assetName, OHLC{}))
			genk.AddStrategy(strategy)

			//genk.SetBalance(100000)
			now := time.Now()
			start := now.AddDate(-15, -9, -2)
			genk.AddStartDate(start)
			now = time.Now()
			end := now.AddDate(0, -3, -2)
			genk.AddEndDate(end)

			//genk.RunEventBased()
			dataManager := K.NewCSVDataManager(genk)
			//genk.AddDataManager(dataManager)

			//dataManager.ReadCSVFilesAsync([]string{"test/data/ABB.csv", "test/data/ASSAb.csv"})
			count := dataManager.ReadCSVFile(localFilename)

			log.WithFields(log.Fields{
				"count": count,
			}).Info("Number of lines processed")

			wg.Done()
		}(fileName)
		y++
	}
	wg.Wait()

	log.WithFields(log.Fields{
		"balance": portfolio.GetBalance(),
	}).Info("Balance")
}

func TestRun(t *testing.T) {
	defer profile.Start().Stop()
	//t.Parallel()
	readFolder("test/datarepo/")
}