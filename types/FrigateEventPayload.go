package types

type FrigateEventPayload struct {
	Type   string       `json:"type"`
	Before FrigateEvent `json:"before"`
	After  FrigateEvent `json:"after"`
}
