package dialogflow

type EntityType struct {
	Name                  string
	DisplayName           string
	Kind                  string
	AutoExpansionMode     string
	Entities              []Entity
	EnableFuzzyExtraction bool
}

type Entity struct {
	Value    string
	Synonyms []string
}
