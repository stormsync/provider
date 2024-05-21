package hoarder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/stormsync/collector"
	"github.com/stormsync/database"
)

type Transformer struct {
	consumer      sarama.Consumer
	consumerTopic string
	queries       *database.Queries
}

func NewTransformer(ctx context.Context, consumerTopic string, brokers []string, dbConnString string) (*Transformer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatal("Could not create server: ", err)
	}
	conn, err := pgx.Connect(ctx, dbConnString)
	if err != nil {
		return nil, err
	}

	defer conn.Close(ctx)

	queries := database.New(conn)

	return &Transformer{
		consumer:      consumer,
		consumerTopic: consumerTopic,
		queries:       queries,
	}, nil
}

func (t *Transformer) GetMessages(ctx context.Context) error {
	fmt.Println("getting messages")
	if err := t.subscribe(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to consumerTopic %s: %w", t.consumerTopic, err)
	}
	return nil
}

func (t *Transformer) subscribe(ctx context.Context) error {
	fmt.Println("subscribing")
	partitionList, err := t.consumer.Partitions(t.consumerTopic) // get all partitions on the given consumerTopic
	if err != nil {
		fmt.Errorf("failed retrieving partitionList for consumerTopic %s: %w", t.consumerTopic, err)
	}

	initialOffset := sarama.OffsetOldest // get offset for the oldest message on the consumerTopic
	var errs error
	for _, partition := range partitionList {
		pc, err := t.consumer.ConsumePartition(t.consumerTopic, partition, initialOffset)
		if err != nil {
			log.Println("consume partition error: ", err)
			return fmt.Errorf("consumePartition error: %w", err)
		}
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				fmt.Println("msg recieved: ", message.Topic)
				if err := t.messageReceived(ctx, message); err != nil {
					errors.Join(errs, fmt.Errorf("failed to recieve message: %w", err))
				} else {
					log.Println(message.Topic)
				}
			}
		}(pc)
	}
	return errs
}

func (t *Transformer) messageReceived(ctx context.Context, message *sarama.ConsumerMessage) error {
	log.Println("Report Type: ", string(message.Key))
	var errs error
	var err error
	var rptType collector.ReportType
	for _, v := range message.Headers {
		if strings.EqualFold(string(v.Value), "reportType") {
			rptType, err = collector.FromString(string(v.Value))
			if err != nil {
				errors.Join(errs, fmt.Errorf("unable to process report due to unknown type: %w", err))
			}
		}
		continue
	}
	if errs != nil {
		return errs
	}
	dbReport, err := t.ToDBReport(message)
	if err != nil {
		return fmt.Errorf("failed to process message type %s: %w", rptType, err)
	}

	if _, err := t.queries.InsertReport(ctx, []database.InsertReportParams{dbReport}); err != nil {
		errors.Join(errs, err)
		log.Printf("failed to insert %s reported at %v\n", dbReport.RptType, dbReport.ReportedTime)
	}

	return errs
}

func (t *Transformer) ToDBReport(message *sarama.ConsumerMessage) (database.InsertReportParams, error) {
	var errs error
	var err error
	var rptType collector.ReportType
	for _, v := range message.Headers {
		if strings.EqualFold(string(v.Value), "reportType") {
			rptType, err = collector.FromString(string(v.Value))
			if err != nil {
				errors.Join(errs, fmt.Errorf("unable to process report due to unknown type: %w", err))
			}
		}
		continue
	}

	var dbReport database.InsertReportParams
	switch rptType {
	case collector.Hail:
		dbReport, err = hailMsgToDBReport(message.Value)
	case collector.Wind:
		dbReport, err = windMsgToDBReport(message.Value)
	case collector.Tornado:
		dbReport, err = tornadoMsgToDBReport(message.Value)
	}
	if err != nil {
		return dbReport, fmt.Errorf("failed to generate database report model for %s: %w", rptType, err)
	}
	return dbReport, nil
}

func (t *Transformer) Poll(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := t.GetMessages(ctx); err != nil {
				log.Println("failed to collect: ", err)
			} else {
				log.Println("messages collected: ")
			}

		}
	}
}
