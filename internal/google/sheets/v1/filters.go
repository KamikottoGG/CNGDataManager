package sheets

func HasRequiredColumn(header ...string) RowFilter {
	return func(m map[string]string) bool {
		for _, col := range header {
			if m[col] == "" {
				return false
			}
		}
		return true
	}
}
