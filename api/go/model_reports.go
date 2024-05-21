package api

import (
	"time"
)

type Reports struct {
	Reports []Report `json:"reports,omitempty"`
}

type Report struct {
	// Date and time the report was generated in UTC time. The date format is YYYY-MM-DD and the time format is HH:MM:SS.SSSZ using 24 hour mode (3PM = 15, 3AM = 03).
	Time   time.Time `json:"Time,omitempty"`
	VarCol string    `json:"Size,omitempty"`
	// The direction, (NNW, NW, SSW, etc), from the known landmark provided as a reference in the location field.
	Direction string `json:"Direction,omitempty"`
	// The distance, in miles, from the known landmark provided as a reference in the location field.
	Distance int32 `json:"Distance,omitempty"`
	// A known landmark to provide as a reference for direction and distance to the hail when reported. This can be a town, city, physical, or man-made feature.
	Location string `json:"Location,omitempty"`
	// County of the referenced hail when it was reported.
	County string `json:"County,omitempty"`
	// Two letter state abbreviation of the referenced hail when it was reported.
	State string `json:"State,omitempty"`
	// Latitude of the referenced hail when it was reported.
	Lat string `json:"Lat,omitempty"`
	// Longitude of the referenced hail when it was reported.
	Lon string `json:"Lon,omitempty"`
	// Remarks and comments regarding the referenced hail when it was reported.
	Comments string `json:"Comments,omitempty"`
	// Reporting weather office.
	Office string `json:"Office,omitempty"`
}

func (r Reports) ToHailReports() HailReports {
	var hrs HailReports
	for _, rpt := range r.Reports {
		hrs.Reports = append(hrs.Reports, HailReport{
			Time:      rpt.Time,
			Size:      rpt.VarCol,
			Direction: rpt.Direction,
			Distance:  rpt.Distance,
			Location:  rpt.Location,
			County:    rpt.County,
			State:     rpt.State,
			Lat:       rpt.Lat,
			Lon:       rpt.Lon,
			Comments:  rpt.Comments,
			Office:    rpt.Office,
		})
	}
	return hrs
}

func (r Reports) ToWindReports() WindReports {
	var wrs WindReports
	for _, rpt := range r.Reports {
		wrs.Reports = append(wrs.Reports, WindReport{
			Time:      rpt.Time,
			Speed:     rpt.VarCol,
			Direction: rpt.Direction,
			Distance:  rpt.Distance,
			Location:  rpt.Location,
			County:    rpt.County,
			State:     rpt.State,
			Lat:       rpt.Lat,
			Lon:       rpt.Lon,
			Comments:  rpt.Comments,
			Office:    rpt.Office,
		})
	}
	return wrs
}

func (r Reports) ToTornadoReports() TornadoReports {
	var trs TornadoReports
	for _, rpt := range r.Reports {
		trs.Reports = append(trs.Reports, TornadoReport{
			Time:      rpt.Time,
			FScale:    rpt.VarCol,
			Direction: rpt.Direction,
			Distance:  rpt.Distance,
			Location:  rpt.Location,
			County:    rpt.County,
			State:     rpt.State,
			Lat:       rpt.Lat,
			Lon:       rpt.Lon,
			Comments:  rpt.Comments,
			Office:    rpt.Office,
		})
	}
	return trs
}
