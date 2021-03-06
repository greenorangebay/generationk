package generationk

import (
	"errors"
)

//Directon is used to describe an order
type Directon int

const (
	//BuyOrder order
	BuyOrder Direction = iota
	//SellOrder order
	SellOrder
	//ShortOrder order to short
	ShortOrder
	//CoverOrder to cover a shrot
	CoverOrder
	//Zero = 0
	ZERO = 0
	//Empty account
	EMPTY = 0.0
)

var UnknownDirection = errors.New("Unknown type of direction for order - should be Buy, Sell, Short or Cover")
var ZeroQuantity = errors.New("Quantity <0")

type OrderType int

const (
	// A market order is an order to buy or sell a stock at the market’s
	// current best available price. A market order typically ensures
	// an execution but it does not guarantee a specified price.
	MarketOrder OrderType = iota

	//A limit order is an order to buy or sell a stock with a restriction on
	// the maximum price to be paid or the minimum price to be received
	// (the “limit price”). If the order is filled, it will only be at the
	// specified limit price or better. However, there is no assurance of
	// execution. A limit order may be appropriate when you think you can
	// buy at a price lower than—or sell at a price higher than—the
	// current quote.
	// Maximum price for buys, or minimum price for sells, at which the
	// order should be filled.
	LimitOrder

	// A stop order is an order to buy or sell a stock at the market price once
	// the stock has traded at or through a specified price (the “stop price”).
	// If the stock reaches the stop price, the order becomes a market order and
	// is filled at the next available market price. If the stock fails to reach
	// the stop price, the order is not executed.
	// For sells, the order will be placed if market price falls below this value.
	// For buys, the order will be placed if market price rises above this value.
	StopOrder

	// A stop order is an order to buy or sell a stock at the market price once
	// the stock has traded at or through a specified price (the “stop price”).
	// If the stock reaches the stop price, a limit order is placed
	StopLimitOrder
)

//Broker is used to send orders
type Broker struct {
	portfolio *Portfolio
	callback  OrderStatus
	comission Comission
}

//SetComission is used to set a comission scheme
func (b *Broker) SetComission(comission Comission) {
	b.comission = comission
}

//SendOrder is used to place an order with the broker
func (b *Broker) SendOrder(order Order, orderstatus OrderStatus) (float64, error) {
	b.callback = orderstatus

	switch order.direction {

	case BuyOrder:

		if order.Qty <= ZERO {
			return 0, ZeroQuantity
		}

		cost, err := b.buy(order)

		if err != nil {
			b.rejected(err)
			return 0, err
		}

		return cost, nil

	case SellOrder:

		if order.Qty == 0 {
			return 0, ZeroQuantity
		}

		if order.Qty == -1 {
			qty, err := b.portfolio.GetQuantity(Holding{assetName: order.Asset})
			if err != nil {
				return 0.0, err
			}

			order.Qty = qty

		}

		cash, err := b.sell(order)

		if err != nil {
			b.rejected(err)
			return 0, err
		}

		return cash, nil

	case ShortOrder:
		b.sellshort(order)

	case CoverOrder:
		b.cover(order)

	default:
		return 0, UnknownDirection

	}

	return b.portfolio.GetBalance(), nil
}

//getAmountForQty is used to get the amount needed to buy a certain quantity
func getAmountForQty(order Order) float64 {

	return order.Price * float64(order.Qty)
}

//getQtyForAmount is used to get how many many stocks we get for a certain amount
/*func getQtyForAmount(order Order) int {
	return int(order.Amount / order.Asset.Close())
}*/

//accpeted is used to send an order event with an accepted event
func (b Broker) accepted(order Order) {
	if b.callback != nil {
		b.callback.OrderEvent(Accepted{})
	}
}

//rejected is used to send a rejected order event
func (b Broker) rejected(err error) {
	if b.callback != nil {
		b.callback.OrderEvent(Rejected{err: err})
	}

}

//fill is used to send a rejected order event
func (b Broker) fill(fill Fill) {
	if b.callback != nil {
		b.callback.OrderEvent(fill)
	}

}

//buy is used to buy a qty or a possible qty from an amount, if a comission is
// set is will be used and deducted from the account
func (b *Broker) buy(order Order) (float64, error) {

	amount := getAmountForQty(order)

	/*if order.Amount > EMPTY {
		amount = order.Amount
		qty := getQtyForAmount(order)
		order.Qty = qty
	}*/

	if b.comission != nil {
		amount += b.comission.GetComisson(order.Price, order.Qty)
	}

	err := b.portfolio.subtractFromBalance(amount)
	if err != nil {
		b.rejected(err)

		return 0, err
	}

	b.accepted(order)
	b.portfolio.AddHolding(Holding{
		qty:       order.Qty,
		assetName: order.Asset,
		price:     order.Price,
		time:      order.Time,
	})

	b.fill(Fill{
		Qty:       order.Qty,
		AssetName: order.Asset,
		Price:     order.Price,
		Time:      order.Time,
	})

	return -amount, nil
}

//sell is used to sell a holding and book and the profits or losses
func (b *Broker) sell(order Order) (float64, error) {

	amount := getAmountForQty(order)
	b.portfolio.addToBalance(amount)

	b.accepted(order)

	err := b.portfolio.RemoveHolding(Holding{
		qty:       order.Qty,
		assetName: order.Asset,
		price:     order.Price,
		time:      order.Time,
	})

	b.fill(Fill{
		Qty:       -order.Qty,
		AssetName: order.Asset,
		Price:     order.Price,
		Time:      order.Time,
	})

	return amount, err
}

//sellshort is not implemented
func (b *Broker) sellshort(order Order) {
}

//cover is used to cover a short, not implemented
func (b *Broker) cover(order Order) {
}
