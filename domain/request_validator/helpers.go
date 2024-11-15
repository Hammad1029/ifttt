package requestvalidator

func GenerateAll(req *map[string]RequestParameter) error {
	for key, val := range *req {
		if regex, err := val.Generate(); err != nil {
			return err
		} else {
			val.Regex = regex
			(*req)[key] = val
		}
	}
	return nil
}
