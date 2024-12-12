package service

import "back-end/internal/repository"

var (
	GlobalDefaultPassword    string
	GlobalAdmin              *AdminService
	GlobalCustomer           *repository.CustomerRepository
	GlobalEmployee           *repository.EmployeeRepository
	GlobalDish               *repository.DishRepository
	GlobalCategory           *repository.CategoryRepository
	GlobalCartItem           *repository.CartItemRepository
	GlobalOrder              *repository.OrderRepository
	GlobalOrderItem          *repository.OrderItemRepository
	GlobalOrderStatusHistory *repository.OrderStatusHistoryRepository
	GlobalReview             *repository.ReviewRepository
)
