package validators

import "regexp"

func IsEmailValid(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) < 3 || len(email) > 254 || !rxEmail.MatchString(email) { // This checks if the email input from user matches the conditions setup above, so we Matchstring
		return false // Is OR condition is true and rxEmail condition is false which is negated so true i.e email does not match regexp, then return false, here || if anyone is 0 when answer is zero(false)
		// If OR condition is false, in case of normal email and email matches regexp - true which is the Negated so it becomes false it geos below
	}
	return true
}
