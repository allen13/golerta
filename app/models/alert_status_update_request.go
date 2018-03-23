package models

// AlertStatusUpdateRequest holds a status update alert request
type AlertStatusUpdateRequest struct {
	Status                  string `json:"status"`
	Text                    string `json:"text"`
	AcknowledgementDuration int    `json:"acknowledgement_duration"`
}
