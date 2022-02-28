package objectTypes
func include(includeDeprecate bool, isDeprecated bool) bool{
	include := true;
	switch includeDeprecate{
	case true:
		include = true;
	case false:
		include = !isDeprecated;
	}

	return include;
}