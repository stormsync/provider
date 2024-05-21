package api

import (
	"time"
)

type WindReport struct {
	// Date and time the report was generated in UTC time. The date format is YYYY-MM-DD and the time format is HH:MM:SS.SSSZ using 24 hour mode (3PM = 15, 3AM = 03).
	Time time.Time `json:"Time,omitempty"`
	// Number indicating the speed of wind gusts in miles per hour (MPH).
	Speed string `json:"Speed,omitempty"`
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
	// Remarks and comments regarding the referenced hail when it was reported. - \"UTILITY POLES WERE DOWNED IN THE 3200 BLOCK OF DEARBORN STREET AS WELL AS DAVID RAINES ROAD. (SHV)\"
	Comments string `json:"Comments,omitempty"`
	// Reporting weather office.
	Office string `json:"Office,omitempty"`
}
