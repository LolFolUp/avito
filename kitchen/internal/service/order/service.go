package order

type Service struct {
	restaurants RestaurantRepository
	menu        MenuRepository
	orders      OrderRepository
	txManager   TxManager
}

func New(
	restaurants RestaurantRepository,
	menu MenuRepository,
	orders OrderRepository,
	txManager TxManager,
) *Service {
	return &Service{
		restaurants: restaurants,
		menu:        menu,
		orders:      orders,
		txManager:   txManager,
	}
}
