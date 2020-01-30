package chat

// FilterWordType an enum for types of
type FilterWordType int

const (
	ImmediateBan FilterWordType = 1 // Bans user without warning
	Blacklisted  FilterWordType = 2 // blocks message
	Censored     FilterWordType = 3 // replaces message with a happy word
)

// FilterWord a type for a word to be checked by the filter
type FilterWord struct {
	FilterType FilterWordType // what to do about it
	Intensity  int            // how bad it is
	Word       string         // what it is
}

// CheckSentence checks if a sentence contains any words in a blocklist
func CheckSentence(sentence string, blocklist []FilterWord) {

}
