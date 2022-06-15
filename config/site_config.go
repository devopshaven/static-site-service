package config

type SiteConfig struct {
	Hosting Hosting `json:"hosting"`
}

type Hosting struct {
	Public         string          `json:"public"`
	Ignore         []string        `json:"ignore,omitempty"`
	Redirects      []Redirect      `json:"redirects,omitempty"`
	Rewrites       []Rewrite       `json:"rewrites,omitempty"`
	Headers        []HostingHeader `json:"headers,omitempty"`
	CleanUrls      bool            `json:"cleanUrls"`
	TrailingSlash  bool            `json:"trailingSlash"`
	AppAssociation string          `json:"appAssociation"`
}

type HostingHeader struct {
	SourceRegex string         `json:"source"`
	Headers     []KeyValuePair `json:"headers"`
}

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Redirect struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        int64  `json:"type"`
}

type Rewrite struct {
	Source       string  `json:"source"`
	Destination  *string `json:"destination,omitempty"`
	DynamicLinks *bool   `json:"dynamicLinks,omitempty"`
	Function     *string `json:"function,omitempty"`
	Run          *Run    `json:"run,omitempty"`
}

type Run struct {
	ServiceID string `json:"serviceId"`
	Region    string `json:"region"`
}
