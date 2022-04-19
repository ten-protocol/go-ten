package playground

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func panicIfAnyErr(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}
