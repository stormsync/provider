package consumer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stormsync/collector"
	"github.com/stormsync/database"
	report "github.com/stormsync/transformer/proto"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	Reader   *kafka.Reader
	Topic    string
	Address  string
	user     string
	password string
	logger   *slog.Logger
	db       *database.Queries
}

// NewConsumer generates a new kafka provider.
func NewConsumer(address, topic, user, pw, groupID string, logger *slog.Logger, db *database.Queries) (*Consumer, error) {
	mechanism, err := scram.Mechanism(scram.SHA256, user, pw)
	if err != nil {
		return nil, fmt.Errorf("failed to create scram.Mechanism for auth: %w", err)
	}
	readerConfig := kafka.ReaderConfig{
		GroupID:        groupID,
		CommitInterval: time.Second,
		Brokers:        []string{"organic-ray-9236-us1-kafka.upstash.io:9092"},
		Topic:          topic,
		Dialer: &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		}}
	reader := kafka.NewReader(readerConfig)

	return &Consumer{
		Reader:   reader,
		Topic:    topic,
		Address:  address,
		user:     user,
		password: pw,
		logger:   logger,
		db:       db,
	}, nil

}

// ReadMessage allows for the consumer to read a message from a topic
// and commit that message.
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	message, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		fmt.Errorf("failed to read message: %w", err)
		return kafka.Message{}, err
	}
	if err := c.Reader.CommitMessages(ctx, message); err != nil {
		fmt.Errorf("failed to commit message: %w", err)
	}

	return message, err
}

// GetMessage pulls a message off of the topic, transforms it,
// applies business logic if needed, marshals, and moves onto
// another topic for further processing.
func (c *Consumer) GetMessage(ctx context.Context) error {
	msg, err := c.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}
	c.logger.Debug("incoming message", "msg value", string(msg.Value))

	reportType, err := getReportTypeFromHeader(msg.Headers)
	if err != nil {
		return fmt.Errorf("unable to extract message type from header: %w", err)
	}

	c.logger.Debug("report type found", "type", reportType.String())

	irp, err := c.processMessage(reportType, msg.Value)
	if err != nil {
		return fmt.Errorf("failed to process message %s\n%s\nerror: %w", reportType, msg.Value, err)
	}
	c.logger.Info("Inserting Record", "type", irp.RptType)
	if _, err := c.db.InsertReport(ctx, []database.InsertReportParams{irp}); err != nil {
		// quietly ignore duplicate key errors
		// TODO - use upset style function in db to insert
		//  or use MERGE on pg 15
		if !strings.Contains(err.Error(), "duplicate key value violates") {
			c.logger.Debug("failed to write message to database", "irp", fmt.Sprintf("%#+v", irp))
			return fmt.Errorf("failed to insert into database: %w", err)
		}
	}

	return nil
}

func (c *Consumer) InsertReportIntoDB(ctx context.Context, irp database.InsertReportParams) error {
	c.logger.Debug("Inserting record", "irp type", irp.RptType, "state", irp.State)
	_, err := c.db.InsertReport(ctx, []database.InsertReportParams{irp})
	if err != nil {
		return fmt.Errorf("failed to write IRP record %#+v  : %w", irp, err)
	}
	return nil
}

// getReportTypeFromHeader extracts the report type from the message headers
// func getReportTypeFromHeader(hdrs []kafka.Header) (collector.ReportType, error) {
func getReportTypeFromHeader(hdrs []kafka.Header) (collector.ReportType, error) {
	var rptType collector.ReportType
	var err error

	if len(hdrs) == 0 {
		return rptType, errors.New("headers do not contain report type")
	}

	reportERR := errors.New("unable to find reportType key, cannot determine report type")
	for _, v := range hdrs {
		if strings.EqualFold(v.Key, "reportType") {
			rptType, err = collector.FromString(string(v.Value))
			if err != nil {
				return rptType, fmt.Errorf("unable to process report due to unknown type: %w", err)
			}
			reportERR = nil
			break
		}
	}
	return rptType, reportERR
}

// processMessage performs the logic to get a msg off the topic and
// get it unmarshaled and turned into a dbReport ready to insert into the database.
func (c *Consumer) processMessage(rptType collector.ReportType, msg []byte) (database.InsertReportParams, error) {
	var err error
	var irp database.InsertReportParams

	switch rptType {
	case collector.Hail:
		irp, err = processHailMessage(c.logger, msg)
	case collector.Wind:
		irp, err = processWindMessage(c.logger, msg)
	case collector.Tornado:
		irp, err = processTornadoMessage(c.logger, msg)
	default:
		err = fmt.Errorf("unknown report type %q", rptType.String())
	}

	return irp, err
}

func processHailMessage(logger *slog.Logger, msg []byte) (database.InsertReportParams, error) {
	var irp database.InsertReportParams
	if msg == nil {
		return irp, errors.New("line cannot be nil")
	}

	hailMsg := report.HailMsg{}

	if err := proto.Unmarshal(msg, &hailMsg); err != nil {
		return irp, nil
	}
	logger.Debug("unmarshalled hail message", "type", hailMsg.Type)

	rptTime := time.Unix(hailMsg.Time, 0)

	irp = database.InsertReportParams{
		RptType: database.ReportTypeHail,
		ReportedTime: pgtype.Timestamptz{
			Time:  rptTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: hailMsg.Size,
			Valid: true,
		},
		DistFromLocation:    hailMsg.GetDistance(),
		HeadingFromLocation: hailMsg.GetDirection(),
		County:              hailMsg.GetCounty(),
		State: pgtype.Text{
			String: hailMsg.GetState(),
			Valid:  true,
		},
		Latitude: pgtype.Text{
			String: hailMsg.GetLat(),
			Valid:  true,
		},
		Longitude: pgtype.Text{
			String: hailMsg.GetLon(),
			Valid:  true,
		},
		EventLocation: nil,
		Comments: pgtype.Text{
			String: hailMsg.GetRemarks(),
			Valid:  true,
		},
		Location: hailMsg.GetLocation(),
		NwsOffice: pgtype.Text{
			String: "",
			Valid:  true,
		},
	}

	return irp, nil
}

func processWindMessage(logger *slog.Logger, msg []byte) (database.InsertReportParams, error) {
	logger.Debug("processing wind message")
	var irp database.InsertReportParams
	if msg == nil {
		return irp, errors.New("line cannot be nil")
	}
	windMsg := report.WindMsg{}
	if err := proto.Unmarshal(msg, &windMsg); err != nil {
		return irp, err
	}
	rptTime := time.Unix(windMsg.Time, 0)
	logger.Debug("unmarshalled wind message", "type", windMsg.Type)

	irp = database.InsertReportParams{
		RptType: database.ReportTypeWind,
		ReportedTime: pgtype.Timestamptz{
			Time:  rptTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: windMsg.Speed,
			Valid: true,
		},
		DistFromLocation:    windMsg.GetDistance(),
		HeadingFromLocation: windMsg.GetDirection(),
		County:              windMsg.GetCounty(),
		State: pgtype.Text{
			String: windMsg.GetState(),
			Valid:  true,
		},
		Latitude: pgtype.Text{
			String: windMsg.GetLat(),
			Valid:  true,
		},
		Longitude: pgtype.Text{
			String: windMsg.GetLon(),
			Valid:  true,
		},
		EventLocation: nil,
		Comments: pgtype.Text{
			String: windMsg.GetRemarks(),
			Valid:  true,
		},
		Location: windMsg.GetLocation(),
		NwsOffice: pgtype.Text{
			String: "",
			Valid:  true,
		},
	}

	return irp, nil
}

func processTornadoMessage(logger *slog.Logger, msg []byte) (database.InsertReportParams, error) {
	var irp database.InsertReportParams
	if msg == nil {
		return irp, errors.New("line cannot be nil")
	}
	tornadoMsg := report.TornadoMsg{}
	if err := proto.Unmarshal(msg, &tornadoMsg); err != nil {
		return irp, err
	}
	logger.Debug("unmarshalled tornado message", "type", tornadoMsg.Type)
	rptTime := time.Unix(tornadoMsg.Time, 0)

	irp = database.InsertReportParams{
		RptType: database.ReportTypeTornado,
		ReportedTime: pgtype.Timestamptz{
			Time:  rptTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: tornadoMsg.F_Scale,
			Valid: true,
		},
		DistFromLocation:    tornadoMsg.GetDistance(),
		HeadingFromLocation: tornadoMsg.GetDirection(),
		County:              tornadoMsg.GetCounty(),
		State: pgtype.Text{
			String: tornadoMsg.GetState(),
			Valid:  true,
		},
		Latitude: pgtype.Text{
			String: tornadoMsg.GetLat(),
			Valid:  true,
		},
		Longitude: pgtype.Text{
			String: tornadoMsg.GetLon(),
			Valid:  true,
		},
		EventLocation: nil,
		Comments: pgtype.Text{
			String: tornadoMsg.GetRemarks(),
			Valid:  true,
		},
		Location: tornadoMsg.GetLocation(),
		NwsOffice: pgtype.Text{
			String: "",
			Valid:  true,
		},
	}

	return irp, nil
}
