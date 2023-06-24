package transformer

import (
	"github.com/devfullcycle/imersao13/go/internal/market/dto"
	"github.com/devfullcycle/imersao13/go/internal/market/entity"
)

// Aqui recebemos o JSON do Kafka do canal de entrada, e temos que trasformar o JSON em dados que nossa aplicação entenda
func TransformInput(input dto.TradeInput) *entity.Order { //função - Recebo os dados do Kafka (JSON cru) e moldo eles confirme a aplicação precisa
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000) //Criando uma ação
	investor := entity.NewInvestor(input.InvestorID)             //Criando um investidor com os dados enviados

	order := entity.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType) //Criando uma ordem de serviço

	if input.CurrentShares > 0 { //Verificando se o investidor já tem alguma ação da qual esta querendo vender ou comprar
		assetPosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares) //Cria um novo objeto de ações de investidor
		investor.AddAssetPosition(assetPosition)                                             //adiciona o valor da ação na carteira ja existente do investidor
	}
	return order //retornar a ordem
}

// Aqui retornamos o JSON para o Kafka no canal de saida, temos que trasformar nossos dados da aplicação, e um objeto que o Kafka entenda
func TransformOutput(order *entity.Order) *dto.OrderOutput {
	var transactionsOutput []*dto.TransactionOutput //Criando uma variavel com o tipo []*dto.TransactionOutput

	for _, t := range order.Transactions { //percorrendo as transações dessa ordem
		transactionOutput := &dto.TransactionOutput{ //criando o objeto de Trasação
			TransactionID: t.ID,
			BuyerID:       t.BuyingOrder.Investor.ID,
			SellerID:      t.SellingOrder.Investor.ID,
			AssetID:       t.SellingOrder.Asset.ID,
			Price:         t.Price,
			Shares:        t.Shares,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput) //Adicionando o objeto criado dentro do Array[] de transactionsOutput
	}

	output := &dto.OrderOutput{ //Criando objeto de retorno output
		OrderID:            order.ID,
		InvestorID:         order.Investor.ID,
		AssetID:            order.Asset.ID,
		OrderType:          order.OrderType,
		Status:             order.Status,
		Partial:            order.PendingShares,
		Shares:             order.Shares,
		TransactionsOutput: transactionsOutput,
	}

	return output //retornando o output
}
