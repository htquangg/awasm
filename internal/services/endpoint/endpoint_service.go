package endpoint

type (
	EndpointRepo interface{}

	EndpointService struct {
		endpointRepo EndpointRepo
	}
)

func NewEndpointService(endpointRepo EndpointRepo) *EndpointService {
	return &EndpointService{
		endpointRepo: endpointRepo,
	}
}
