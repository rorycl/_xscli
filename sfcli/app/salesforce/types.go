package salesforce

import "time"

// SOQLResponse is the top-level envelope for a SOQL query response.
type SOQLResponse struct {
	TotalSize      int           `json:"totalSize"`
	Done           bool          `json:"done"`
	NextRecordsURL string        `json:"nextRecordsUrl"`
	Records        []Opportunity `json:"records"`
}

// Opportunity represents the data for a single Opportunity record.
type Opportunity struct {
	ID                string     `json:"Id"`
	Name              string     `json:"Name"`
	Amount            float64    `json:"Amount"`
	CloseDate         string     `json:"CloseDate"` // Kept as string for simplicity in DB
	StageName         string     `json:"StageName"`
	RecordType        RecordType `json:"RecordType"`
	LastModifiedDate  time.Time  `json:"LastModifiedDate"`
	PayoutReference   *string    `json:"Payout_Reference__c"` // Pointer to handle null values
}

// RecordType is a nested object within the Opportunity response.
type RecordType struct {
	Name string `json:"Name"`
}

// BatchRequest is the structure for a Salesforce Composite Batch API request.
type BatchRequest struct {
	AllOrNone bool         `json:"binary"`
	Requests  []Subrequest `json:"batchRequests"`
}

// Subrequest represents a single operation within a batch request.
type Subrequest struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Body   interface{} `json:"richInput"`
}

// BatchResponse is the structure of the response from the Composite Batch API.
type BatchResponse struct {
	HasErrors bool     `json:"hasErrors"`
	Results   []Result `json:"results"`
}

// Result represents the outcome of a single subrequest within a batch response.
type Result struct {
	StatusCode int         `json:"statusCode"`
	Result     interface{} `json:"result"`
}
