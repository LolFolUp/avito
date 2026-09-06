package order

type Service struct {
	restaurants RestaurantRepository
	partner     RestaurantClient
	orders      OrderRepository
	txManager   TxManager
}

func New(
	restaurants RestaurantRepository,
	partner RestaurantClient,
	orders OrderRepository,
	txManager TxManager,
) *Service {
	return &Service{
		restaurants: restaurants,
		partner:     partner,
		orders:      orders,
		txManager:   txManager,
	}
}
