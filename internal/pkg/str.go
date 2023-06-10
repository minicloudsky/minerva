package pkg

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func FindStr(strSlice []string, target string) bool {
	for _, s := range strSlice {
		if s == target {
			return true
		}
	}

	return false
}

func RemoveStr(oldStrSLice []string, target string) []string {
	newStrSlice := make([]string, 0)
	for i, s := range oldStrSLice {
		if s == target {
			newStrSlice = append(newStrSlice, oldStrSLice[i+1:]...)
			return newStrSlice
		}
		newStrSlice = append(newStrSlice, s)
	}

	return newStrSlice
}
