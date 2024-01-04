package goutilspkg

func IsInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func RemoveListDups(list []string) []string {
	keys := make(map[string]bool)
	var listWithoutDups []string
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			listWithoutDups = append(listWithoutDups, entry)
		}
	}
	return listWithoutDups
}
