package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stormsync/database"
)

func (s ServerAndDB) GetHailReports(c echo.Context) error {
	qp := c.QueryParams()

	// s.DB.GetAllHailReports()
	if len(qp) == 0 {
		dbReports, errResp := s.DB.GetReportsByType(c.Request().Context(), database.ReportTypeHail)
		if errResp != nil {
			return c.JSON(500, "error makin query to database")
		}
		return c.JSON(200, dbToReportModel(dbReports).ToHailReports())
	}

	if val, ok := qp["date"]; ok {
		dbReports, errResponse := s.getReportsByDateAndType(c, val[0], database.ReportTypeHail)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToHailReports())
	}
	var qState, qCounty string
	if val, ok := qp["state"]; ok {
		qState = val[0]
	}

	if val, ok := qp["county"]; ok {
		qCounty = val[0]
	}

	if qState != "" && qCounty != "" {
		dbReports, errResponse := s.getReportsByTypeByCountyAndState(c, qCounty, qState, database.ReportTypeHail)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToHailReports())
	}

	if qState != "" {
		dbReports, errResponse := s.getReportsByTypeByState(c, qState, database.ReportTypeHail)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToHailReports())
	}

	return c.JSON(400, ApiResponse{
		Code:    400,
		Message: "valid query params are state, county, location, and date",
	})

}
