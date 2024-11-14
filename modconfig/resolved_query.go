package modconfig

// ResolvedQuery contains the execute SQL, raw SQL and args string used to execute a query
type ResolvedQuery struct {
	Name       string
	ExecuteSQL string
	RawSQL     string
	Args       []any

	IsMetaQuery bool
}
