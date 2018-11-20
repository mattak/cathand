package cathand

func AssertError(err error) {
	if err != nil {
		panic("Error: " + err.Error())
	}
}
