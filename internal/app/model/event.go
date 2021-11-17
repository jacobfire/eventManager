package model

type Event struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Time string `json:"time"`
	//DateFrom string `json:"dateFrom"`
	//DateTo string `json:"dateTo"`
	//TimeFrom string `json:"timeFrom"`
	//TimeTo string `json:"timeTo"`
	Timezone string `json:"timezone"`
	Duration int32 `json:"duration"`
	Notes []string `json:"notes"`
}

type ExtendedEvent struct {
	Notes string
	Event
}