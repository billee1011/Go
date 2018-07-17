package util

//合并任意个数组
func MergeStringArray(strings[][]string) []string{
	return mergeStringArray(nil,0,0,0,strings)
}

/**
	递归合并数组，调用时result传入nil,lastIndex,maxIndex,rIndex=0
 */
func mergeStringArray(result []string, maxIndex int, lastIndex int, rIndex int, strings[][]string) []string {
	if result == nil {
		for _, v := range strings {
			maxIndex += len(v)
		}
		result = make([]string, maxIndex)
	}

	if maxIndex == rIndex{
		return result
	}

	for _, v := range strings[lastIndex] {
		result[rIndex] = v
		rIndex ++
	}
	lastIndex++
	return mergeStringArray(result,maxIndex,lastIndex,rIndex,strings)
}
