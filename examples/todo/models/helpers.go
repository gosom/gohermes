package models

func AutoMigrateModels() []interface{} {
	ans := []interface{}{
		CustomUser{}, Todo{},
	}
	return ans
}
