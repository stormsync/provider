package api

import (
	"github.com/labstack/echo/v4"
	"github.com/stormsync/database"
)

func (s ServerAndDB) GetWindReports(c echo.Context) error {
	qp := c.QueryParams()

	// s.DB.GetAllWindReports()
	if len(qp) == 0 {
		dbReports, errResp := s.DB.GetReportsByType(c.Request().Context(), database.ReportTypeWind)
		if errResp != nil {
			return c.JSON(500, "error makin query to database")
		}
		return c.JSON(200, dbToReportModel(dbReports).ToWindReports())
	}

	if val, ok := qp["date"]; ok {
		dbReports, errResponse := s.getReportsByDateAndType(c, val[0], database.ReportTypeWind)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToWindReports())
	}
	var qState, qCounty string
	if val, ok := qp["state"]; ok {
		qState = val[0]
	}

	if val, ok := qp["county"]; ok {
		qCounty = val[0]
	}

	if qState != "" && qCounty != "" {
		dbReports, errResponse := s.getReportsByTypeByCountyAndState(c, qCounty, qState, database.ReportTypeWind)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToWindReports())
	}

	if qState != "" {
		dbReports, errResponse := s.getReportsByTypeByState(c, qState, database.ReportTypeWind)
		if errResponse.Code > 0 {
			return c.JSON(int(errResponse.Code), errResponse)
		}
		return c.JSON(200, dbToReportModel(dbReports).ToWindReports())
	}

	return c.JSON(400, ApiResponse{
		Code:    400,
		Message: "valid query params are state, county, location, and date",
	})

}
