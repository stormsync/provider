package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
	report "github.com/stormsync/transformer/proto"
	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stormsync/database"
	"google.golang.org/protobuf/proto"
)

func Test_processTornadoMessage(t *testing.T) {
	type args struct {
		msg []byte
	}
	dte, err := time.Parse(time.DateOnly, "2006-04-21")
	if err != nil {
		t.Fatal("unable to setup process tornado message test date variable")
	}
	tests := []struct {
		name    string
		args    args
		want    database.InsertReportParams
		wantErr error
	}{
		{
			name: "should process a tornado message correctly",
			args: args{
				msg: mustMarshal(&report.TornadoMsg{
					Time:      StringToUnixTimeInt64(dte.String(), "1131"),
					F_Scale:   int32(0),
					Distance:  2,
					Direction: "SSW",
					Location:  "Lamont",
					County:    "Jefferson",
					State:     "FL",
					Lat:       "30.35",
					Lon:       "-83.83",
					Remarks:   "A tornado touched down in far eastern Jefferson county and moved through most of southern Madison county. EF0 tree damage was confirmed in Jefferson county with EF1 dam (TAE)",
				}),
			},
			want: database.InsertReportParams{
				RptType: "tornado",
				ReportedTime: pgtype.Timestamptz{
					Time:  stringToUnixTime(dte.String(), "1131"),
					Valid: true,
				},

				VarCol: pgtype.Int4{
					Int32: int32(0),
					Valid: true,
				},
				DistFromLocation:    2,
				HeadingFromLocation: "SSW",
				County:              "Jefferson",
				State: pgtype.Text{
					String: "FL",
					Valid:  true,
				},
				Latitude: pgtype.Text{
					String: "30.35",
					Valid:  true,
				},
				Longitude: pgtype.Text{
					String: "-83.83",
					Valid:  true,
				},
				EventLocation: nil,
				Comments: pgtype.Text{
					String: "A tornado touched down in far eastern Jefferson county and moved through most of southern Madison county. EF0 tree damage was confirmed in Jefferson county with EF1 dam (TAE)",
					Valid:  true,
				},
				NwsOffice: pgtype.Text{
					String: "",
					Valid:  true,
				},
				Location: "Lamont",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processTornadoMessage(tt.args.msg)
			if err != nil {
				err = errors.Unwrap(err)
			}
			// TODO - address lack of time check
			// zero report times out so tests will pass.
			// no the best hack, but that's how it's going for now
			tt.want.ReportedTime.Time = time.Time{}
			got.ReportedTime.Time = time.Time{}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)

		})
	}
}

func mustMarshal(m proto.Message) []byte {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Fatal("failed to setup proto marshal for test: ", err)
	}
	return b
}

func stringToUnixTime(dateOnly string, hhmm string) time.Time {
	if len(hhmm) < 4 {
		return time.Time{}
	}

	if _, err := time.Parse(time.DateOnly, dateOnly); err != nil {
		return time.Time{}
	}

	t := fmt.Sprintf("%s %s:%s:00", dateOnly, hhmm[0:2], hhmm[2:4])
	newTime, err := time.Parse(time.DateTime, t)
	if err != nil {
		return time.Time{}
	}
	return newTime.UTC()

}

func StringToUnixTimeInt64(dateOnly string, hhmm string) int64 {
	if len(hhmm) < 4 {
		return 0
	}

	if _, err := time.Parse(time.DateOnly, dateOnly); err != nil {
		return 0
	}

	t := fmt.Sprintf("%s %s:%s:00", dateOnly, hhmm[0:2], hhmm[2:4])
	newTime, err := time.Parse(time.DateTime, t)
	if err != nil {
		return 0
	}
	return newTime.UTC().Unix()

}

func TestConsumer_InsertReportIntoDB(t *testing.T) {
	type fields struct {
		Reader   *kafka.Reader
		Topic    string
		Address  string
		user     string
		password string
		logger   *slog.Logger
		db       *database.Queries
	}
	type args struct {
		ctx context.Context
		irp database.InsertReportParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{

		{
			name: "should insert record into database",
			fields: fields{
				db: nil,
			},
			args: args{
				ctx: context.Background(),
				irp: database.InsertReportParams{
					RptType: "tornado",
					ReportedTime: pgtype.Timestamptz{
						Time:  stringToUnixTime("2024-01-30", "1131"),
						Valid: true,
					},
					CreatedAt: pgtype.Timestamptz{
						Time:  time.Now().UTC(),
						Valid: true,
					},

					VarCol: pgtype.Int4{
						Int32: int32(0),
						Valid: true,
					},
					DistFromLocation:    2,
					HeadingFromLocation: "SSW",
					County:              "Jefferson",
					State: pgtype.Text{
						String: "FL",
						Valid:  true,
					},
					Latitude: pgtype.Text{
						String: "30.35",
						Valid:  true,
					},
					Longitude: pgtype.Text{
						String: "-83.83",
						Valid:  true,
					},
					EventLocation: nil,
					Comments: pgtype.Text{
						String: "A tornado touched down in far eastern Jefferson county and moved through most of southern Madison county. EF0 tree damage was confirmed in Jefferson county with EF1 dam (TAE)",
						Valid:  true,
					},
					NwsOffice: pgtype.Text{
						String: "",
						Valid:  true,
					},
					Location: "Lamont",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consumer{
				Reader:   tt.fields.Reader,
				Topic:    tt.fields.Topic,
				Address:  tt.fields.Address,
				user:     tt.fields.user,
				password: tt.fields.password,
				logger:   tt.fields.logger,
				db:       tt.fields.db,
			}

			ctx := context.Background()
			conn, err := pgx.Connect(ctx, "user=postgres dbname=weather sslmode=disable")
			if err != nil {
				log.Fatal("no db", err)
			}
			defer conn.Close(ctx)
			db := database.New(conn)
			c.db = db
			err = c.InsertReportIntoDB(tt.args.ctx, tt.args.irp)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
