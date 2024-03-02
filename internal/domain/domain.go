package domain

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Token{},
	}
}
