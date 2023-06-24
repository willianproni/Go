package kafka

//Arquivo responsavel por produzir informações no Kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer struct { //Criando strutura do Producer
	ConfigMap *ckafka.ConfigMap //Configurações do Kafka importadas da biblioteca
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer { //Criando um novo Producer
	return &Producer{
		ConfigMap: configMap,
	}
}

func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap) //Criando um novo dado no kafka
	if err != nil {                                  //Verifica se deu erro ao tentar criar um novo Producer
		return err //Retorna o erro criado
	}

	message := &ckafka.Message{ //Criando a variavel de mensagem
		TopicPartition: ckafka.TopicPartition{ //Objeto de config do topic
			Topic:     &topic,              //Nome do topico
			Partition: ckafka.PartitionAny, //quantidade de partições
		},
		Key:   key,          //chave
		Value: msg.([]byte), //valor da mensagem (JSON, string, int...)
	}

	err = producer.Produce(message, nil) //Função que cria a mensagem no Kafka

	if err == nil { //Verifica se deu erro
		return err //Se der erro, retornar  erro
	}

	return nil //Se produzir sem problemas, retorna null
}
