package dto

type TradeInput struct { //Ordem de entrada investidor querendo comprar ou vender ações
	OrderID       string  `json:"order_id"`       //ID da Ordem
	InvestorID    string  `json:"investor_id"`    //ID do investidor
	AssetID       string  `json:"asset_id"`       //ID da Ações
	CurrentShares int     `json:"current_shares"` //Quantidade de ações que o investidor tem
	Shares        int     `json:"shares"`         //Quantidade de ações que quer negociar
	Price         float64 `json:"price"`          //Preço por ações
	OrderType     string  `json:"order_type"`     //Tipo da ordem (compra/buy ou venda/sell)
} //OBS: Não temos Transacations aqui pois essa é uma ordem de solicitação de compra/venda, assim não tendo nenhuma troca de dinheiro ou ação

type OrderOutput struct { //Ordem de saida, de trasferencia quando uma ação foi comprada ou vendida
	OrderID            string               `json:"order_id"`    //ID da Ordem
	InvestorID         string               `json:"investor_id"` //ID do investidor
	AssetID            string               `json:"asset_id"`    //ID da Ações
	OrderType          string               `json:"order_type"`  //Tipo da ordem (compra/buy ou venda/sell)
	Status             string               `json:"status"`      //Status da Ordem, se está pendente, vendido, em negociação
	Partial            int                  `json:"partial"`     //Quantidade de ações para comprar ou vender faltantes
	Shares             int                  `json:"shares"`      //Quantidade ações negociadas
	TransactionsOutput []*TransactionOutput `json:"transactions"`
}

type TransactionOutput struct {
	TransactionID string  `json:"transaction_id"` //ID da transaçã0
	BuyerID       string  `json:"buyer_id"`       //ID do investidor comprador
	SellerID      string  `json:"seller_id"`      //ID do investidor vendedor
	AssetID       string  `json:"asset_id"`       //ID do asset negociado
	Price         float64 `json:"price"`          //O preço que foi negociado da ação
	Shares        int     `json:"shares"`         //Quantidade de ações vendidas na transação
}
