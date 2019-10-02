package queue

import (
	"github.com/streadway/amqp"
)

type FeedQueue struct {
	ConnectionString string
	Connection       *amqp.Connection
	Channel          *amqp.Channel
}

func (q *FeedQueue) ConnectAndPrepareQueue(connStr string) error {
	queueConn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}
	q.Connection = queueConn
	ch, err := q.Connection.Channel()
	if err != nil {
		return err
	}
	q.Channel = ch
	err = ch.ExchangeDeclare(
		"igfeed", // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	queue, err := ch.QueueDeclare(
		"igfeed_pic", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}
	err = ch.QueueBind(
		queue.Name,   // queue name
		"igfeed_pic", // routing key
		"igfeed",     // exchange
		false,
		nil)
	if err != nil {
		return err
	}
	return nil
}
