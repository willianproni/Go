package entity

type Asset struct { //Criação do entidade Asset ação
	ID           string //tem um ID
	Name         string //Um nome
	MarketVolume int
}

func NewAsset(id string, Name string, MarketVolume int) *Asset {
	return &Asset{
		ID:           id,
		Name:         Name,
		MarketVolume: MarketVolume,
	}
}
