package repository

type QueryBuilder struct {
	query string
}

func (ref *QueryBuilder) WithAnd(condition string) *QueryBuilder {
	if ref.query != "" {
		ref.query += " AND "
	}
	ref.query += condition
	return ref
}

func (ref *QueryBuilder) WithOr(condition string) *QueryBuilder {
	if ref.query != "" {
		ref.query += " OR "
	}
	ref.query += condition
	return ref
}

func (ref *QueryBuilder) WithAndSubQuery(subQuery *QueryBuilder) *QueryBuilder {
	if subQuery.query != "" {
		if ref.query != "" {
			ref.query += " AND "
		}
		ref.query += "(" + subQuery.query + ")"
	}
	return ref
}

func (ref *QueryBuilder) Build() string {
	return ref.query
}
