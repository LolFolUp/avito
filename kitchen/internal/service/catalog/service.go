package catalog

type Service struct {
	restaurants RestaurantRepository
	menu        MenuRepository
	txManager   TxManager
}

func New(restaurants RestaurantRepository, menu MenuRepository, txManager TxManager) *Service {
	return &Service{
		restaurants: restaurants,
		menu:        menu,
		txManager:   txManager,
	}
}
