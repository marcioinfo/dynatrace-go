package customerservice

import (
	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/google/uuid"
)

type CustomerService struct {
	ID         string
	ServiceID  string
	CustomerID string
	Name       string
	Document   string
	BirthDate  datetypes.CustomDate
	Email      string
	Phone      string
	Gender     string
	CreatedAt  datetypes.CustomDateTime
	UpdatedAt  datetypes.CustomDateTime
	DeletedAt  datetypes.CustomDateTime
}

func NewCustomerService(
	ID string,
	ServiceID string,
	CustomerID string,
	Name string,
	Document string,
	BirthDate datetypes.CustomDate,
	Email string,
	Phone string,
	Gender string,
) *CustomerService {
	return &CustomerService{
		ID:         ID,
		ServiceID:  ServiceID,
		CustomerID: CustomerID,
		Name:       Name,
		Document:   Document,
		BirthDate:  BirthDate,
		Email:      Email,
		Phone:      Phone,
		Gender:     Gender,
	}
}

func (c *CustomerService) InitID() {
	c.ID = uuid.New().String()
}
