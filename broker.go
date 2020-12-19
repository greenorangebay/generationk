package generationk

import (
	log "github.com/sirupsen/logrus"
)

//OrderType is used to describe an order
type OrderType int

const (
	//Buy order
	Buy OrderType = iota
	//Sell order
	Sell
	//SellShort order
	SellShort
	//Cover short order
	Cover
)

//Broker is used to send orders
type Broker struct {
	portfolio Portfolio
	channel   chan Event
}

//PlaceOrder is used to place an order with the broker
func (b *Broker) PlaceOrder(order Order) {
	log.WithFields(log.Fields{
		"ordertype": order.Ordertype,
		"asset":     (*order.Asset).Name,
		"time":      order.Time,
		"amount":    order.Amount,
	}).Debug("BROKER>PLACE BUY ORDER")

	switch order.Ordertype {
	case Buy:
		go b.buy(order)
	case Sell:
		go b.sell(order)
	case SellShort:
		go b.sellshort(order)
	case Cover:
		go b.cover(order)
	}
}

func getAmountForQty(order Order) float64 {
	return order.Asset.Close() * float64(order.Qty)
}

func getQtyForAmount(order Order) int {
	return int(order.Amount / order.Asset.Close())
}

func (b *Broker) buy(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> BUY")
	if order.Qty > 0 {
		err := b.portfolio.updateCash(getAmountForQty(order))
		if err != nil {
			b.channel <- Rejected{message: "Insufficient funds"}
		}
	}
	b.portfolio.AddHolding(Holding{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time})
	b.channel <- Fill{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time}
	log.Info("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sell(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> SELL")
	b.channel <- Fill{Qty: order.Qty, AssetName: order.Asset.Name, Price: order.Asset.Close(), Time: order.Time}
	log.Info("BROKER> Put FILL EVENT in queue")
}

func (b *Broker) sellshort(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> SELLSHORT")
}

func (b *Broker) cover(order Order) {
	log.WithFields(log.Fields{
		"Order": order,
	}).Info("BROKER> COVER")
}