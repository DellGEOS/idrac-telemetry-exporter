package redfishmetricreport

type DellOem struct {
	Type 							string `json:"@odata.type"`
	ServiceTag 						string `json:"ServiceTag"`
	MetricReportDefinitionDigest	string `json:"MetricReportDefinitionDigest"`
	IDRACFirmwareVersion 			string `json:"iDRACFirmwareVersion"`
}

type Oem struct {
	Dell DellOem `json:"Dell"`
}

type DellOemMetricValue struct {
	Type string `json:"@odata.type"`
	ContextID string `json:"ContextID"`
	Label string `json:"Label"`
	Source string `json:"Source"`
	FQDD string `json:"FQDD"`
}

type OemMetricValue struct {
	Dell DellOemMetricValue `json:"Dell"`
}

type MetricReportDefinition struct {
	ODataId string `json:"@odata.id"`
}

type MetricValue struct {
	MetricId string `json:"MetricId"`
	Timestamp string `json:"Timestamp"`
	Value string `json:"MetricValue"`
	MetricProperty string `json:"MetricProperty"`
	Oem OemMetricValue `json:"Oem"`
}

type MetricReport struct {
	Type string `json: "@odata.type"`
	Context string `json:"@odata.context"`
	ODataId string `json:"@odata.id"`
	Id string `json:"Id"`
	Name string `json:"Name"`
	Timestamp string `json:"Timestamp"`
	MetricReportDef MetricReportDefinition `json:"MetricReportDefinition"`
	MetricValues []MetricValue `json:"MetricValues"`
	MetricValuesCount int `json:"MetricValues@odata.count"`
	OemSection Oem `json:"Oem"`
}

type MetricReportList struct {
	Type string `json: "@odata.type"`
	Context string `json:"@odata.context"`
	ODataId string `json:"@odata.id"`
	Name string `json:"Name"`
	Members []MetricReportDefinition `json:"Members"`
	MembersCount int `json:"Members@odata.count"`
}