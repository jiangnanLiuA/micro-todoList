package dao

import (
	"micro_todoList/app/task/repository/db/model"
)

func migration() {
	_db.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&model.Task{})
}
