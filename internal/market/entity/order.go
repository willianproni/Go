package entity

type Order struct {
	ID            string         //Id
	Investor      *Investor      //Investidor dono da ordem
	Asset         *Asset         //Ação
	Shares        int            //Quantas ele quer 10
	PendingShares int            //Quantas estão pende nte 3
	Price         float64        //Preço
	OrderType     string         //Tipo da ordem (buy or sell)
	Status        string         //Status (pendente, confirmada. aberta)
	Transactions  []*Transaction //Transção (caso a ordem aconteça) ela é array pois posso comprar ações de 2 ou mais pessoas, dependendo da quantidade que o investidor querer
}

func NewOrder(orderID string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order { //Criando uma nova ordem
	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN", //Toda nova ordem vai ser aberta, e quando concluir a transação ela vai ser virar closed
		Transactions:  []*Transaction{},
	}
}
