package dialogflow

type Intent struct {
	Name                     string
	DisplayName              string
	Priority                 int32
	IsFallback               bool
	TrainingPhrases          []TrainingPhrase
	Action                   string
	InputContextNames        []string
	OutputContexts           []Context
	Parameters               []Parameter
	Messages                 []Message
	RootFollowupIntentName   string
	ParentFollowupIntentName string
	FollowupIntents          []Intent
	FollowupIntentInfo       []FollowupIntentInfo
}

type Context struct {
	Name          string
	LifespanCount int32
	Parameters    []string
}

type TrainingPhrase struct {
	Name  string
	Parts []TrainingPhrasePart
}

type TrainingPhrasePart struct {
	Text        string
	EntityType  string
	Alias       string
	UserDefined bool
}

type Parameter struct {
	Name                  string
	DisplayName           string
	Value                 string
	DefaultValue          string
	EntityTypeDisplayName string
	Mandatory             bool
	Prompts               []string
	IsList                bool
}

type Message struct {
	Text     string
	Platform string
}

type FollowupIntentInfo struct {
	FollowupIntentName       string
	ParentFollowupIntentName string
}
