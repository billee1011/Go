package utils

// GetHszFx 根据配牌信息获取换三张方向
func GetHszFx(fxValue string) (fx int) {
	switch fxValue {
	case "dui":
		fx = 1
	case "shun":
		fx = 0
	case "ni":
		fx = 2
	default:
		fx = -1
	}
	return
}
