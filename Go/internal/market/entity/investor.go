package entity

// Asset significa ação
type Investor struct { //Tenhos o Investidor
	ID            string                   //Id
	Name          string                   //Nome do meu Investidor
	AssetPosition []*InvestorAssetPosition //Quantidade de ações (Carteira de investimentos)
}

type InvestorAssetPosition struct { //Temos a posição do investirod
	AssetID string //Id da ação
	Shares  int    //quantidade
}

func NewInvestor(id string) *Investor { //Criado um novo investidor do 0, sem nenhuma ação (New User)
	return &Investor{ //Retornando um investidor vazio
		ID:            id,                         //Recebe um ID aleatório
		AssetPosition: []*InvestorAssetPosition{}, //Cria uma carteira de investimentos vazia
	}
}

//Metodos
//Funções para uma determinada entidade, aqui temos uma função para o Investor

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) { //Criando um metodo quee fala a posição na carteira de investimento
	i.AssetPosition = append(i.AssetPosition, assetPosition) //Adiciona uma nova posição
}

//append adiciona um novo valor o slicer array []

func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) { //atualiza o quantidade de ações
	assetPosition := i.GetAssetPosition(assetID) //Busca a ação pelo ID
	if assetPosition == nil {                    //valor null
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qtdShares)) //caso valor for vazio, cria um novo valor na carteira
	} else { //caso encontre
		assetPosition.Shares += qtdShares //faz uma soma com os valores existentes
	}
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition { //Procurando a posição
	for _, assetPosition := range i.AssetPosition { //Fazendo um for percorrendo a carteira de ações do investir i.
		if assetPosition.AssetID == assetID { //Verificando se o AssetID é igual o assetID de alguma ação na carteiras de investimento
			return assetPosition //se for igual retornar
		}
	}
	return nil //se não for igual retorna nil (null | em branco)
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition { //Função para criar uma nova posição de ação
	return &InvestorAssetPosition{ //retornar a posição
		AssetID: assetID,
		Shares:  shares,
	}
}
