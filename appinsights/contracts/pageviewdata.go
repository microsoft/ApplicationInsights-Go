package contracts

// NOTE: This file was automatically generated.

// An instance of PageView represents a generic action on a page like a button
// click. It is also the base type for PageView.
type PageViewData struct {
	Domain
	EventData

	// Request URL with all query string parameters
	Url string `json:"url"`

	// Request duration in format: DD.HH:MM:SS.MMMMMM. For a page view
	// (PageViewData), this is the duration. For a page view with performance
	// information (PageViewPerfData), this is the page load time. Must be less
	// than 1000 days.
	Duration string `json:"duration"`
}

// Creates a new PageViewData instance with default values set by the schema.
func NewPageViewData() *PageViewData {
	return &PageViewData{
		EventData: EventData{
			Ver:          2,
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		},
	}
}

// Returns the name used when this is embedded within an Envelope container.
func (data *PageViewData) EnvelopeName() string {
	return "Microsoft.ApplicationInsights.PageView"
}

// Returns the base type when placed within a Data object container.
func (data *PageViewData) BaseType() string {
	return "PageViewData"
}
