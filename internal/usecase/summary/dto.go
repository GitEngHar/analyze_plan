package summary

type ResultResourceAddress struct {
	Create  []string
	Delete  []string
	Update  []string
	Replace []string
}
type ResourceNames []string
type ResultResourceTypeToNames struct {
	Delete  map[string]ResourceNames
	Replace map[string]ResourceNames
}

type ResultResourceCount struct {
	Create  int
	Delete  int
	Update  int
	Replace int
}
