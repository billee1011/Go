package common

type EventParams struct {
	Params   []interface{}
}

func CreateEventParams(param...interface{}) EventParams{
	return EventParams{Params:param}
}