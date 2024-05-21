package api

type StormReports struct {
	HailReports *HailReport `json:"HailReports,omitempty"`

	TornadoReports *TornadoReport `json:"TornadoReports,omitempty"`

	WindReports *WindReport `json:"WindReports,omitempty"`
}
