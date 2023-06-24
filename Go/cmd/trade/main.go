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

func main() {
	orderIn := make(chan *entity.Order)  //Recebemos os dados via Kafka
	orderOut := make(chan *entity.Order) //Enviamos os dados via Kafka
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}
	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan) //Tread 2

	//Recebe do canal do kafka, joga no input, processa joga no output e depois publica no kafka
	book := entity.NewBook(orderIn, orderOut, wg)

	go book.Trade() //Tread 3

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			orderIn <- order
		}
	}()

	for res := range orderOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", " ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
