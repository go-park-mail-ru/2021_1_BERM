package models

type OrderInfo struct {
	OrderName   string
	CustomerId  uint64
	ExecutorId  uint64
	Budget      uint64
	DeadLine    uint64
	Description string
	Category    string
}
