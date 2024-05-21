package hoarder

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stormsync/database"
	report "github.com/stormsync/transformer/proto"
	"google.golang.org/protobuf/proto"
)

func hailMsgToDBReport(msgValue []byte) (database.InsertReportParams, error) {
	var msg report.HailMsg
	if err := proto.Unmarshal(msgValue, &msg); err != nil {
		return database.InsertReportParams{}, nil
	}

	uTime := time.Unix(msg.GetTime(), 0).UTC()

	dbReport := database.InsertReportParams{
		RptType: database.ReportTypeHail,
		ReportedTime: pgtype.Timestamptz{
			Time:  uTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: msg.Size,
			Valid: true,
		},
		DistFromLocation:    msg.Distance,
		HeadingFromLocation: msg.Direction,
		County:              msg.County,
		State: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Latitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Longitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		EventLocation: msg.Location,
		Comments: pgtype.Text{
			String: msg.State,
			Valid:  true},
		NwsOffice: pgtype.Text{
			String: msg.State,
			Valid:  true}, // TODO: add parsing the office to the transformer svc.
	}
	return dbReport, nil
}
func windMsgToDBReport(msgValue []byte) (database.InsertReportParams, error) {
	var msg report.WindMsg
	if err := proto.Unmarshal(msgValue, &msg); err != nil {
		return database.InsertReportParams{}, nil
	}

	uTime := time.Unix(msg.GetTime(), 0).UTC()

	dbReport := database.InsertReportParams{
		RptType: database.ReportTypeHail,
		ReportedTime: pgtype.Timestamptz{
			Time:  uTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: msg.Speed,
			Valid: true,
		},
		DistFromLocation:    msg.Distance,
		HeadingFromLocation: msg.Direction,
		County:              msg.County,
		State: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Latitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Longitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		EventLocation: msg.Location,
		Comments: pgtype.Text{
			String: msg.State,
			Valid:  true},
		NwsOffice: pgtype.Text{
			String: msg.State,
			Valid:  true}, // TODO: add parsing the office to the transformer svc.
	}
	return dbReport, nil
}
func tornadoMsgToDBReport(msgValue []byte) (database.InsertReportParams, error) {
	var msg report.TornadoMsg
	if err := proto.Unmarshal(msgValue, &msg); err != nil {
		return database.InsertReportParams{}, nil
	}

	uTime := time.Unix(msg.GetTime(), 0).UTC()

	dbReport := database.InsertReportParams{
		RptType: database.ReportTypeTornado,
		ReportedTime: pgtype.Timestamptz{
			Time:  uTime,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		VarCol: pgtype.Int4{
			Int32: msg.F_Scale,
			Valid: true,
		},
		DistFromLocation:    msg.Distance,
		HeadingFromLocation: msg.Direction,
		County:              msg.County,
		State: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Latitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		Longitude: pgtype.Text{
			String: msg.State,
			Valid:  true},
		EventLocation: msg.Location,
		Comments: pgtype.Text{
			String: msg.State,
			Valid:  true},
		NwsOffice: pgtype.Text{
			String: msg.State,
			Valid:  true}, // TODO: add parsing the office to the transformer svc.
	}
	return dbReport, nil
}
