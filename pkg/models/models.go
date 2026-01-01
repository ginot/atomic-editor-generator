package models

// Project represents the entire application project
type Project struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Brand       Brand        `json:"brand"`
	GlobalStyles GlobalStyles `json:"globalStyles"`
	ThirdParty  ThirdParty   `json:"thirdParty"`
}

type Brand struct {
	Colors     map[string]string    `json:"colors"`
	Typography Typography           `json:"typography"`
	Spacing    map[string]string    `json:"spacing"`
	Breakpoints map[string]string   `json:"breakpoints"`
}

type Typography struct {
	FontFamily  map[string]string `json:"fontFamily"`
	FontSizes   map[string]string `json:"fontSizes"`
	FontWeights map[string]interface{} `json:"fontWeights"`
}

type GlobalStyles struct {
	Reset     bool `json:"reset"`
	Normalize bool `json:"normalize"`
}

type ThirdParty struct {
	Analytics     *ThirdPartyService `json:"analytics,omitempty"`
	CookieConsent *ThirdPartyService `json:"cookieConsent,omitempty"`
	Recaptcha     *ThirdPartyService `json:"recaptcha,omitempty"`
	Fonts         *FontsConfig       `json:"fonts,omitempty"`
}

type ThirdPartyService struct {
	Provider  string `json:"provider"`
	ID        string `json:"id,omitempty"`
	ScriptSrc string `json:"scriptSrc,omitempty"`
	DataKey   string `json:"dataKey,omitempty"`
	Version   string `json:"version,omitempty"`
	SiteKey   string `json:"siteKey,omitempty"`
}

type FontsConfig struct {
	Google []string `json:"google,omitempty"`
}

// Clump represents a group of related pages
type Clump struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	BaseRoute string   `json:"baseRoute"`
	Pages     []string `json:"pages"`
}

// Page represents a single page in the application
type Page struct {
	ID     string            `json:"id"`
	Clump  string            `json:"clump"`
	Route  string            `json:"route"`
	Title  string            `json:"title"`
	Meta   map[string]string `json:"meta"`
	Layout string            `json:"layout"`
}

// Layout defines the structure of a page
type Layout struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Structure []LayoutSection  `json:"structure"`
}

type LayoutSection struct {
	Section   string   `json:"section"`
	Organism  string   `json:"organism,omitempty"`
	Organisms []string `json:"organisms,omitempty"`
	Position  string   `json:"position,omitempty"`
}

// Atom represents the smallest UI element
type Atom struct {
	ID       string                 `json:"id"`
	Subatom  string                 `json:"subatom"`
	Config   map[string]interface{} `json:"config,omitempty"`
	Styles   map[string]interface{} `json:"styles,omitempty"`
	States   map[string]map[string]interface{} `json:"states,omitempty"`
}

type Atoms struct {
	Images   []Atom `json:"images,omitempty"`
	Headings []Atom `json:"headings,omitempty"`
	Links    []Atom `json:"links,omitempty"`
	Buttons  []Atom `json:"buttons,omitempty"`
	Inputs   []Atom `json:"inputs,omitempty"`
	Text     []Atom `json:"text,omitempty"`
}

// Molecule represents a combination of atoms
type Molecule struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Atoms      map[string]string      `json:"atoms,omitempty"`
	Responsive []ResponsiveConfig     `json:"responsive,omitempty"`
	Styles     map[string]interface{} `json:"styles,omitempty"`
	States     map[string]map[string]interface{} `json:"states,omitempty"`
	Events     map[string]Event       `json:"events,omitempty"`
}

type ResponsiveConfig struct {
	Breakpoint string            `json:"breakpoint"`
	Atoms      map[string]string `json:"atoms"`
}

type Event struct {
	Action string                 `json:"action"`
	Target string                 `json:"target,omitempty"`
	Method string                 `json:"method,omitempty"`
	Params []string               `json:"params,omitempty"`
	Condition string              `json:"condition,omitempty"`
	ClassName string              `json:"className,omitempty"`
}

// Organism represents a complex UI section
type Organism struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Atoms          map[string]string      `json:"atoms,omitempty"`
	Molecules      interface{}            `json:"molecules,omitempty"` // can be map or array
	Sections       []OrganismSection      `json:"sections,omitempty"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Styles         map[string]interface{} `json:"styles,omitempty"`
	Layout         map[string]interface{} `json:"layout,omitempty"` // can contain strings or nested maps
	States         map[string]map[string]interface{} `json:"states,omitempty"`
	Events         map[string]Event       `json:"events,omitempty"`
	Behavior       *Behavior              `json:"behavior,omitempty"`
	ControlStyles  map[string]map[string]interface{} `json:"controlStyles,omitempty"`
	IndicatorStyles *IndicatorStyles      `json:"indicatorStyles,omitempty"`
}

type OrganismSection struct {
	Type      string   `json:"type"`
	Molecules []string `json:"molecules"`
}

type Behavior struct {
	Type         string `json:"type"`
	Autoplay     bool   `json:"autoplay,omitempty"`
	Interval     int    `json:"interval,omitempty"`
	Loop         bool   `json:"loop,omitempty"`
	Controls     bool   `json:"controls,omitempty"`
	Indicators   bool   `json:"indicators,omitempty"`
	Transition   string `json:"transition,omitempty"`
	PauseOnHover bool   `json:"pauseOnHover,omitempty"`
}

type IndicatorStyles struct {
	Container map[string]interface{} `json:"container"`
	Dot       map[string]interface{} `json:"dot"`
	DotActive map[string]interface{} `json:"dotActive"`
}

// AtomicStructure is the root structure containing all elements
type AtomicStructure struct {
	Project   Project     `json:"project"`
	Clump     Clump       `json:"clump"`
	Page      Page        `json:"page"`
	Layout    Layout      `json:"layout"`
	Atoms     Atoms       `json:"atoms"`
	Molecules []Molecule  `json:"molecules"`
	Organisms []Organism  `json:"organisms"`
}
