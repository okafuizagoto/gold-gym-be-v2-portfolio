package goldgym

type InsertSubsAll struct {
	HeaderData SubscriptionAll      `json:"header"`
	DetailData []SubscriptionDetail `json:"detail"`
}
