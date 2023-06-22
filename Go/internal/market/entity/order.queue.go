package entity

type OrderQueue struct {
	Orders []*Order
}

// metodos

func (oq *OrderQueue) Less(i, j int) bool { //comparando dos valores se é menor e retornando um boolean
	return oq.Orders[i].Price < oq.Orders[j].Price
}

func (oq *OrderQueue) Swap(i, j int) { //invertendo posições
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

func (oq *OrderQueue) Len() int { //Retorna o tamanho da order
	return len(oq.Orders)
}

func (oq *OrderQueue) Push(x interface{}) { //Adiciona um novo dado
	oq.Orders = append(oq.Orders, x.(*Order))
}

func (oq *OrderQueue) Pop() interface{} {
	old := oq.Orders
	numberOrders := len(old)
	item := old[numberOrders-1]
	oq.Orders = old[0 : numberOrders-1]
	return item
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}

//Less - Comparar dois valores
//Swap - Inverte os valores
//Len - ver o tamanho dos dados
//Push - adiciona novos dados (appned)
//Pop - responsavel por remover uma posição
