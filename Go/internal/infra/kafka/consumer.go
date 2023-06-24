package kafka

//Arquivo responsavel por consumir informações no Kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Consumer struct {
	ConfigMap *ckafka.ConfigMap //Configuração do Kafka, vem da biblioteca
	Topics    []string          //Topics que quero escutar na aplicação
}

func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer { //Criando um Novo consumer Kafka
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

// Metodo
func (c *Consumer) Consume(msgChan chan *ckafka.Message) error { //Canal que recebe as mensagens e retorna para outro tread ler
	consumer, err := ckafka.NewConsumer(c.ConfigMap) //Consumir a variavel de consumo
	//OBS: NewConsumer é uma função propria da biblioteca do kafka, não sendo a criada na arquivo
	if err != nil { //Validar se retornou um Erro
		panic(err) //Apresentar o erro
	}
	err = consumer.SubscribeTopics(c.Topics, nil) //Verificar se está se escrevendo nos topics
	if err != nil {
		panic(err)
	}
	for { //loop infinito que fica ledo as mensagens do kafka
		msg, err := consumer.ReadMessage(-1) //lendo a mensagem
		if err == nil {                      //verificar se recebe erro
			msgChan <- msg //jogar a mensagem no canal, para outro tread ler os dados desse canal
		}
	}

}
