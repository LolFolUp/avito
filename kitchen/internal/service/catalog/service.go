package catalog

type Service struct {
	restaurants RestaurantRepository
	partner     RestaurantClient
}

func New(restaurants RestaurantRepository, partner RestaurantClient) *Service {
	return &Service{
		restaurants: restaurants,
		partner:     partner,
	}
}
