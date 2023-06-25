package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/devfullcycle/imersao13/go/internal/infra/kafka"
	"github.com/devfullcycle/imersao13/go/internal/market/dto"
	"github.com/devfullcycle/imersao13/go/internal/market/entity"
	"github.com/devfullcycle/imersao13/go/internal/market/transformer"
)

func main() { //função Main, principal do Golang
	//Start variaveis

	orderIn := make(chan *entity.Order)  //Recebemos os dados referente a ordem de entrada, solicitações de compra e venda via Kafka
	orderOut := make(chan *entity.Order) //Enviamos os dados via Kafka
	wg := &sync.WaitGroup{}              //Espera a função acabar
	defer wg.Wait()                      //Executa por ultimo, depois da função

	//End variaveis
	//--------------------------------
	//Start criação da conexão Kafka

	kafkaMsgChan := make(chan *ckafka.Message) //Criando o canal do Kafka para receber as mensagens
	configMap := &ckafka.ConfigMap{            //Criar a configMap, configurações para o kafka funcionar (conectar)
		"bootstrap.servers": "host.docker.internal:9094", //Utilizando Docker
		"group.id":          "myGroup",                   //Id do grupo, serve para juntar grupos de consumidores
		"auto.offset.reset": "latest",                    //Seria o momento da leitura de dados earliest - Desde o inicio | latest -
	}

	//End Conexão Kafka
	//--------------------------------
	//Start Criação objeto de comunicação com o Kafka
	producer := kafka.NewKafkaProducer(configMap)            //Criando objeto responsavel por produzir informações para o kafka
	kafka := kafka.NewConsumer(configMap, []string{"input"}) //Criando objeto Responsavel por consumir as mensagem do kafka

	go kafka.Consume(kafkaMsgChan) //Tread 2 - Iniciando função responsavel por consumir as mensagem do kafka

	//End Criação objeto de comunicação com o Kafka
	//--------------------------------
	//Start BOOK

	book := entity.NewBook(orderIn, orderOut, wg) //Criando um book, responsavel por armazenar as informações de compra, venda e transações

	/*
		COMO FUNCIONA O BOOK -->
			1 - Recebe do canal do kafka
			2 - Joga no input
			3 - Processa
			4 - Joga no output
			5 - Publica no kafka
	*/

	go book.Trade() //Tread 3 - Tread responsavel por receber os dados e computar as vendas e compras e gerar as transações

	//End BOOK
	//--------------------------------
	//Start função de receber informações do Kafka e criar ordens no mocrosserviço

	go func() { //Criação de uma função anonima no Kafka para ficar consumindo o canal do kafka
		for msg := range kafkaMsgChan { //For para percorrer os dados no canal do kafka
			wg.Add(1)                                     //Para cada transação, adicionar uma tarefa no wg (watingGroup)
			fmt.Println(string(msg.Value))                //Println para exibir um logo das mensagem recebidas
			tradeInput := dto.TradeInput{}                //Criando uma variavel para receber futuramente os objeto em JSON
			err := json.Unmarshal(msg.Value, &tradeInput) //Unmarshal responsavel por trasformar os dados JSON para objeto
			if err != nil {                               //Verificando se não ou erro
				panic(err) //Se ouve disparar o erro
			}
			order := transformer.TransformInput(tradeInput) //Função responsavel por criar uma nova ordem
			orderIn <- order                                //Com a ordem criada, enviar a ordem para o canal de Entrada de ordens
		}
	}()

	//End função de receber informações do Kafka e criar ordens no microsserviço
	//--------------------------------
	//Start função de enviar informações para o Kafka das ordens processadas no microsserviço

	for res := range orderOut { //For para percorer o canal de ordens
		output := transformer.TransformOutput(res)             //Trasformar os objetos da responsta em dados Cru
		outputJson, err := json.MarshalIndent(output, "", " ") //Trasformar os objetos em dados JSON
		fmt.Println(string(outputJson))                        //Println para exibir o dado no console
		if err != nil {                                        //Verificar se deu erro
			fmt.Println(err) //se der erro, exibir no log
		}
		producer.Publish(outputJson, []byte("orders"), "output") //Publicar a informação no Kafka
	}
	//End função de enviar informações para o Kafka das ordens processadas no microsserviço
}
