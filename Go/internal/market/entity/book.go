package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order         []*Order        //Ter todas as ordens que acontece
	Transactions  []*Transaction  //Todas as transações que estão acontecendo
	OrdersChan    chan *Order     //é um canal de compra e venda //todas as ordem de compra e venda vão cair aqui
	OrdersChanOut chan *Order     //é um canal de ordens de saida
	Wg            *sync.WaitGroup //
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book { //Criado um Book tudo zerado
	return &Book{
		Order:         []*Order{},       //ordens solicitadas
		Transactions:  []*Transaction{}, //transações comcluidas
		OrdersChan:    orderChan,        //canal de ordem (compra e venda)
		OrdersChanOut: orderChanOut,
		Wg:            wg, //await
	}
}

func (b *Book) Trade() { //Função para realizar a compra e venda
	//****
	//Criando duas filas
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)
	// buyOrders := NewOrderQueue()  //fila de compra criando novos objetos
	// sellOrders := NewOrderQueue() //fila de venda criando novos objetos

	// heap.Init(buyOrders)
	// heap.Init(sellOrders)

	//sempre que recebermos uma nova ordem, iremos pegar ela dentro do for, loop infinito
	for order := range b.OrdersChan { //Percorrendo o canal de Orders (compra e venda)
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}
		if order.OrderType == "BUY" { //Verifica se o tipo da ordem é de compra
			//se for ordem de compra
			buyOrders[asset].Push(order)      //Adicionando a ordem na nossa fila de ordem de compra
			if sellOrders[asset].Len() > 0 && //verifica se existe alguma ordem de venda
				sellOrders[asset].Orders[0].Price <= order.Price { //se tiver ordem de venda verifica se o preço da ordem de venda é menor ou igual o preço da ordem de compra
				//Se entrou aqui é porque pode ocorrer uma negociação
				sellOrder := sellOrders[asset].Pop().(*Order) //removo a ordem pois pode ocorrer uma negociação, e não causar conflito com alguma compra futura, (pop) vai retornar a ordem
				if sellOrder.PendingShares > 0 {              //verifica se a ordem retornada está pendente ainda ou já foi finalizada (liquidade/vendida)
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price) //Criar uma transação pois essa ordem tem como negociar
					b.AddTransaction(transaction, b.Wg)                                            //adicionando uma nova transação ao nosso livro
					//toda ordem tem um []array de transações pois podem ter varias compras/vendas em uma mesma proposta
					sellOrder.Transactions = append(sellOrder.Transactions, transaction) //adicionando a transação no []array de transação da ordem de venda
					//order apenas pois já estamos dentro do if do ordem de compra
					order.Transactions = append(order.Transactions, transaction) //adicionando a transação no []array de transação da ordem de venda
					//Canal vai receber as ordem comcluidas para enviar depois para o kafka
					b.OrdersChanOut <- sellOrder     //Enviando a ordem de vende para o canal de concluido
					b.OrdersChanOut <- order         //Enviando a ordem de compra para o canal concluido
					if sellOrder.PendingShares > 0 { //Verificar se a ordem de venda ainda tem ações disponiveis para serem vendas
						sellOrders[asset].Push(sellOrder) //se tiver a ordem de venda volta para fila
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders[asset].Push(order)
			if buyOrders[asset].Len() > 0 &&
				buyOrders[asset].Orders[0].Price >= order.Price {
				//abre a chance de ocorrer negocio
				buyOrder := buyOrders[asset].Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, order.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- order
					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	//defer signifca que vai ser executado apenas no final da função, como ultima coisa
	defer wg.Done() //avisa para o wg que finalizou

	sellingShare := transaction.SellingOrder.PendingShares //Pega a informação de quantas ações estão sendo vendidas 2
	buyingShare := transaction.BuyingOrder.PendingShares   //Pega a informação de quantas ações estão sendo compradas 4
	//preciso saber quem tem menos cota
	//pois quem
	minShare := sellingShare
	if buyingShare < minShare { //Verifica qual é o menor valor referentes as ações de compra e venda
		minShare = buyingShare
	}
	//minShare = 2

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShare) //Atualizar a carteira do vendedor, tirando a quantidade de ações vendidas
	transaction.UpdateBuyPendingShare(-minShare)                                                        //Realiza a subtração da quantidade das ações de venda pendentes

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, minShare) //Atualizar a carteira do vendedor, tirando a quantidade de ações vendidas
	transaction.UpdateSellPendingShare(-minShare)                                                     //Realiza a subtraça2o da quantidade das ações de compra pendentes

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price) //realiza o valor total dessa transação, quantidade * valor de venda

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	b.Transactions = append(b.Transactions, transaction)
}
