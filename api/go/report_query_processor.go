package api

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stormsync/database"
)

func (s ServerAndDB) getReportsByTypeByState(c echo.Context, state string, reportType database.ReportType) ([]database.Report, ApiResponse) {
	rpts, err := s.DB.GetReportsByTypeByState(c.Request().Context(), database.GetReportsByTypeByStateParams{
		State: pgtype.Text{
			String: state,
			Valid:  true,
		},
		RptType: reportType,
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}

func (s ServerAndDB) getReportsByTypeByCountyAndState(c echo.Context, county, state string, reportType database.ReportType) ([]database.Report, ApiResponse) {
	rpts, err := s.DB.GetReportsByTypeByCountyAndState(c.Request().Context(), database.GetReportsByTypeByCountyAndStateParams{
		County: county,
		State: pgtype.Text{
			String: state,
			Valid:  true,
		},
		RptType: reportType,
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}

func (s ServerAndDB) getReportsByDateAndType(c echo.Context, queryDate string, reportType database.ReportType) ([]database.Report, ApiResponse) {
	qt, err := time.Parse("2024-01-02 15:04 -07:00", queryDate)
	if err != nil {
		return nil, ApiResponse{
			Code:    400,
			Message: "date value not valid format, use YYYY-MM-DD HH:MM +00:00 ",
		}
	}

	rpts, err := s.DB.GetAllReportsByTypeByDate(c.Request().Context(), database.GetAllReportsByTypeByDateParams{
		ReportedTime: pgtype.Timestamptz{
			Time:  qt,
			Valid: true,
		},
		RptType: reportType,
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}

func dbToReportModel(rpts []database.Report) Reports {
	var reports Reports
	for _, row := range rpts {
		hr := Report{}
		if row.CreatedAt.Valid {
			hr.Time = row.CreatedAt.Time
		}
		if row.VarCol.Valid {
			hr.VarCol = strconv.FormatInt(int64(row.VarCol.Int32), 10)
		}
		if row.Comments.Valid {
			hr.Comments = row.Comments.String
		}
		if row.NwsOffice.Valid {
			hr.Office = row.NwsOffice.String
		}
		hr.Direction = row.HeadingFromLocation
		hr.Distance = row.DistFromLocation
		hr.Location = row.Location
		hr.County = row.County
		hr.State = row.State.String
		hr.Lat = row.Latitude.String
		hr.Lon = row.Longitude.String
		reports.Reports = append(reports.Reports, hr)
	}
	return reports
}

func (s ServerAndDB) getReportsByState(c echo.Context, state string) ([]database.Report, ApiResponse) {
	rpts, err := s.DB.GetReportsByState(c.Request().Context(), pgtype.Text{
		String: state,
		Valid:  true,
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}

func (s ServerAndDB) getReportsByCountyAndState(c echo.Context, county, state string) ([]database.Report, ApiResponse) {
	rpts, err := s.DB.GetReportsByCountyAndState(c.Request().Context(), database.GetReportsByCountyAndStateParams{
		County: county,
		State: pgtype.Text{
			String: state,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}

func (s ServerAndDB) getReportsByDate(c echo.Context, queryDate string) ([]database.Report, ApiResponse) {
	qt, err := time.Parse("2024-01-02 15:04 -07:00", queryDate)
	if err != nil {
		return nil, ApiResponse{
			Code:    400,
			Message: "date value not valid format, use YYYY-MM-DD HH:MM +00:00 ",
		}
	}

	rpts, err := s.DB.GetAllReportsByDate(c.Request().Context(), pgtype.Timestamptz{
		Time:  qt,
		Valid: true,
	})
	if err != nil {
		return nil, ApiResponse{
			Code:    500,
			Message: "error making query to database",
		}
	}
	return rpts, ApiResponse{}
}
