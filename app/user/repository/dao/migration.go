package dao

import "micro_todoList/app/user/repository/model"

// struct -> 转换成 mysql 的 table
func migration() {
	err := _db.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}
