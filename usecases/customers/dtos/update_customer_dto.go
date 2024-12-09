package customer_dto

import "github.com/adhfoundation/layer-tools/datetypes"

type UpdateCustomerDTO struct {
	Name      *string               `json:"name" validate:"omitempty" example:"João Silva"`
	Email     *string               `json:"email" validate:"omitempty,email" example:"teste@gmail.com"`
	Phone     *string               `json:"phone" validate:"omitempty,numeric" example:"11999999999"`
	Gender    *string               `json:"gender" validate:"omitempty,oneof=male female" example:"male"`
	BirthDate *datetypes.CustomDate `json:"birth_date" validate:"omitempty" example:"2000-11-02"`
}
