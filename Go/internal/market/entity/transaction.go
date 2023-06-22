package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct { //Criação de uma transação
	ID           string
	SellingOrder *Order    //temos que ter uma ordem de venda
	BuyingOrder  *Order    //e uma ordem de compra
	Shares       int       //quantidade de ações movimentadas
	Price        float64   //valor
	Total        float64   //valor total
	DateTime     time.Time //Horario que aconteceu a transação
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price //Calculo para saber o valor total das ações negociadas
	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

// metodo

func (t *Transaction) CalculateTotal(shares int, price float64) { //calcula o valor total da transação
	t.Total = float64(t.Shares) * t.Price
}

func (t *Transaction) UpdateBuyPendingShare(shares int) { //Atualiza a quantidade de ações pendente ainda para compra
	t.BuyingOrder.PendingShares += shares
}

func (t *Transaction) UpdateSellPendingShare(shares int) { //Atualiza a quantidade de ações pendente ainda para venda
	t.SellingOrder.PendingShares += shares
}

func (t *Transaction) CloseBuyOrder() { //Fecha a ordem de compra na transação
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) CloseSellOrder() { //Fecha a ordem de venda na transação
	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}
