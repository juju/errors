package errors

// Since variables are declared before the init block, in order to get the goPath
// we need to return it rather than just reference it.
func GoPath() string {
	return goPath
}

var TrimGoPath = trimGoPath
