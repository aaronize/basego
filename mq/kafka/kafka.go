package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/Shopify/sarama"
    cluster "github.com/bsm/sarama-cluster"
    "log"
    "sync"
    "time"
)

var cli *client

type Config struct {
    Consumer    *consumerConfig  `yaml:"consumer"`
    Producer    *producerConfig  `yaml:"producer"`
}

type consumerConfig struct {
    Enable      bool        `yaml:"enable"`
    Group       string      `yaml:"group"`
    Hosts       []string    `yaml:"hosts"`
    Topic       string      `yaml:"topic"`

    wg          *sync.WaitGroup
}

type producerConfig struct {
    Enable      bool        `yaml:"enable"`
    Hosts       []string    `yaml:"hosts"`
    Topic       string      `yaml:"topic"`
    PoolSize    int         `yaml:"pool_size"`
    Retry       int         `yaml:"retry"`

    wg          *sync.WaitGroup
}

func NewConfig() *Config {
    return &Config{
        Consumer: &consumerConfig{
            Enable: false,
            Topic: "",
            Hosts: nil,
        },
        Producer: &producerConfig{
            Enable: false,
            PoolSize: 300,
            Retry: 3,
        },
    }
}

type client struct {
    producer    *producer
    consumer    *consumer
}

func NewQueueClient() *client {
    return &client{
        producer: nil,
        consumer: nil,
    }
}

func Init(ctx context.Context, config *Config) error {
    fmt.Println("queue config: ", config.Consumer, config.Producer)
    wg := ctx.Value("WaitGroup").(*sync.WaitGroup)
    config.Consumer.wg = wg
    config.Producer.wg = wg

    cli = NewQueueClient()

    go func() {
        <-ctx.Done()
        cli.Close()
    }()

    if config.Consumer.Enable {
        if err := cli.runConsumer(config.Consumer); err != nil {
            return err
        }
    }
    if config.Producer.Enable {
        if err := cli.runProducer(config.Producer); err != nil {
            return err
        }
    }

    return nil
}

func RegisterConsumeProcessor(processor func([]byte)) error {
    if cli.consumer == nil {
        return fmt.Errorf("未启用Consumer，无法注册处理函数")
    }
    cli.consumer.processor = processor

    return nil
}

func (c *client) runConsumer(config *consumerConfig) error {
    consumer := &consumer{
        group: config.Group,
        hosts: config.Hosts,
        topics: []string{config.Topic},

        wg: config.wg,
        quit: make(chan bool),
    }

    if err := consumer.generateConsumerCluster(); err != nil {
        return err
    }
    c.consumer = consumer

    return nil
}

func (c *client) runProducer(config *producerConfig) error {
    producer := &producer{
        hosts: config.Hosts,
        topic: config.Topic,
        retry: config.Retry,
        inputChan: make(chan *sarama.ProducerMessage, config.PoolSize),

        wg: config.wg,
        quit: make(chan bool),
    }

    if err := producer.generateProducer(); err != nil {
        return err
    }
    c.producer = producer

    return nil
}

func (c *client) Close() {
    if c.consumer != nil {
        c.consumer.close()
    }
    if c.producer != nil {
        c.producer.close()
    }
}

type producer struct {
    hosts       []string
    topic       string
    retry       int
    inputChan   chan *sarama.ProducerMessage

    quit        chan bool

    wg      *sync.WaitGroup
}


func (p *producer) generateProducer() error {
    pconf := sarama.NewConfig()
    pconf.Producer.Retry.Max = p.retry
    pconf.Producer.Retry.Backoff = 3 * time.Second
    pconf.Producer.RequiredAcks = 1

    producer, err := sarama.NewAsyncProducer(p.hosts, pconf)
    if err != nil {
        return err
    }
    p.wg.Add(1)

    go func(prod sarama.AsyncProducer) {
        for {
            select {
            case err := <-prod.Errors():
                if err != nil {
                    log.Printf("Input Queue Error, %s\n", err.Error())
                }
            case success := <-prod.Successes():
                if success != nil {
                    bt, _ := success.Value.Encode()
                    log.Printf("Input Queue Success, %s\n", string(bt))
                }
            case <-p.quit:
                return
            }
        }
    }(producer)

    go func(prod sarama.AsyncProducer) {
        for {
            select {
            case msg := <- p.inputChan:
                prod.Input() <- msg
            case <-p.quit:
                log.Printf("Closing Message Queue Producer... \n")
                prod.AsyncClose()

                p.wg.Done()
                log.Printf("Message Queue Producer was closed!\n")
                return
            }
        }
    }(producer)

    return nil
}

func (p *producer) close() {
    if p.quit == nil {
        return
    }
    close(p.quit)
}

// consumer
type consumer struct {
    group   string
    hosts   []string
    topics  []string

    processor func([]byte)

    quit    chan bool

    wg      *sync.WaitGroup
}


func (c *consumer) generateConsumer() error {
    qconf := sarama.NewConfig()
    qconf.Consumer.Return.Errors = true

    master, err := sarama.NewConsumer(c.hosts, qconf)
    if err != nil {
        return err
    }

    consumer, err := master.ConsumePartition(c.topics[0], 0, sarama.OffsetNewest)
    if err != nil {
        return err
    }

    // listen
    go func() {
        for {
            select {
            case err := <- consumer.Errors():
                fmt.Printf("consumer return error: %s\n", err.Error())
            case msg := <-consumer.Messages():
                log.Printf("+++ [Q]Record consume message: Topic: %s, Offset: %d, Partition: %d\n", msg.Topic, msg.Offset, msg.Partition)
                if c.processor == nil {
                    log.Printf("+++ ERROR not register Consumer Processor!!! \n")
                    continue
                }
                c.processor(msg.Value)
            case <-c.quit:
                log.Println("Closing Message Queue Consumer...")
                consumer.AsyncClose()
                return
            }
        }
    }()

    return nil
}

func (c *consumer) generateConsumerCluster() error {
    conf := cluster.NewConfig()
    conf.Consumer.Return.Errors = true
    conf.Group.Return.Notifications = true

    consumer, err := cluster.NewConsumer(c.hosts, c.group, c.topics, conf)
    if err != nil {
        return err
    }
    c.wg.Add(1)

    // Listen
    go func() {
        for {
            select {
            case err := <-consumer.Errors():
                log.Printf("+++ [Q]Consumer return error: %s\n", err.Error())
            case ntf := <-consumer.Notifications():
                log.Printf("+++ [Q]Rebalanced: %+v\n", ntf)
            case msg := <-consumer.Messages():
                log.Printf("+++ [Q]Record consumed message meta: Topic: %s, Offset: %d, Partition: %d\n", msg.Topic, msg.Offset, msg.Partition)
                if c.processor == nil {
                    log.Printf("+++ [Q]ERROR not register Consumer Processor!!! \n")
                    continue
                }
                c.processor(msg.Value)
            case <-c.quit:
                log.Println("Closing Message Queue Consumer...")
                _ = consumer.Close()

                c.wg.Done()
                log.Println("Message Queue Consumer was closed!")
                return
            }
        }
    }()

    return nil
}

func (c *consumer) close() {
    if c.quit == nil {
        return
    }
    close(c.quit)
}

//
func Put2Topic(topic string, msg interface{}) error {
    msgBytes, err := json.Marshal(msg)
    if err != nil {
        return err
    }
    message := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.ByteEncoder(msgBytes),
    }

    cli.producer.inputChan <- message

    return nil
}

func Put(message interface{}) error {
    return Put2Topic(cli.producer.topic, message)
}

