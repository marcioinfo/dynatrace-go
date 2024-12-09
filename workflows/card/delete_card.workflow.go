package card

import (
	"context"
	"os"
	integrationapi "payment-layer-card-api/adapters/integration_api"
	card_token_entities "payment-layer-card-api/entities/card_token"
	"payment-layer-card-api/entities/cards"
	payment_gateway_entities "payment-layer-card-api/entities/payment_gateway"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	customer_token_usecase "payment-layer-card-api/usecases/customer_token"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/adhfoundation/layer-tools/log"
)

type DeleteCardWorkflow struct {
	paymentGateways                               []payment_gateway_entities.PaymentGatewayInterface
	cardsUseCaseGetCardByIDAndCustomerID          *cards_usecase.GetCardByIDAndCustomerID
	cardTokenUseCaseGetByCardIDAndGateway         *card_token_usecase.GetCardTokenByCardIDAndGateway
	customerTokenUseCaseGetByCustomerIDAndGateway *customer_token_usecase.GetByCustomerIDAndGateway
	cardsUseCaseDeleteCard                        *cards_usecase.DeleteCardByID
	cardTokenRepository                           card_token_entities.CardTokenRepositoryInterface
}

func NewDeleteCardWorkflow(
	paymentGateways []payment_gateway_entities.PaymentGatewayInterface,
	cardsUseCaseGetCardByIDAndCustomerID *cards_usecase.GetCardByIDAndCustomerID,
	cardTokenUseCaseGetByCardIDAndGateway *card_token_usecase.GetCardTokenByCardIDAndGateway,
	customerTokenUseCaseGetByCustomerIDAndGateway *customer_token_usecase.GetByCustomerIDAndGateway,
	cardsUseCaseDeleteCard *cards_usecase.DeleteCardByID,
	cardTokenRepository card_token_entities.CardTokenRepositoryInterface,

) *DeleteCardWorkflow {
	return &DeleteCardWorkflow{
		paymentGateways:                               paymentGateways,
		cardsUseCaseGetCardByIDAndCustomerID:          cardsUseCaseGetCardByIDAndCustomerID,
		cardTokenUseCaseGetByCardIDAndGateway:         cardTokenUseCaseGetByCardIDAndGateway,
		customerTokenUseCaseGetByCustomerIDAndGateway: customerTokenUseCaseGetByCustomerIDAndGateway,
		cardsUseCaseDeleteCard:                        cardsUseCaseDeleteCard,
		cardTokenRepository:                           cardTokenRepository,
	}
}

func (c *DeleteCardWorkflow) Execute(ctx context.Context, cardId string, customerId string) *card_dtos.DeleteCardTokenDTO {
	var cardTokens []card_dtos.CardToken

	response := &card_dtos.DeleteCardTokenDTO{
		TokensDeleted: cardTokens,
	}
	card, err := c.findCard(ctx, customerId, cardId)
	if err != nil {
		response.Error = err
		return response
	}

	integrationAdapter := integrationapi.NewIntegrationApiAdapter(os.Getenv("INTEGRATION_API_URL"))

	_, errorDeleteCardIntegration := integrationAdapter.DeleteCardInIntegration(ctx, cardId, customerId)
	if errorDeleteCardIntegration != nil {
		response.Error = errors.NewError(errors.ExternalError, errorDeleteCardIntegration)
	}

	deletedTokens, errorDeleteCardTokens := c.deleteCardTokens(ctx, card)
	if errorDeleteCardTokens != nil {
		response.Error = errorDeleteCardTokens
		return response
	}

	response.TokensDeleted = deletedTokens

	_, errorDeleteCard := c.cardsUseCaseDeleteCard.Execute(ctx, cardId)
	if errorDeleteCard != nil {
		response.Error = errorDeleteCard
		return response
	}

	response.Error = nil

	return response
}

func (c *DeleteCardWorkflow) findCard(ctx context.Context, customerId string, cardId string) (*cards.Card, *errors.ErrorOutput) {
	card, err := c.cardsUseCaseGetCardByIDAndCustomerID.Execute(ctx, cardId, customerId)
	if err != nil && err.Code != errors.CardNotFound {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao buscar o Cartão (%v) para o CustomerID (%v). Erro (%v)", cardId, customerId, err)
		return nil, err
	}

	if card == nil {
		log.Info(ctx).Msgf("O cartão com CardID (%v) não foi encontrado para o CustomerID (%v)", cardId, customerId)
		return nil, errors.NewError(errors.CardNotFound, nil, "Cartão não encontrado.")
	}

	return card, nil
}

func (c *DeleteCardWorkflow) deleteCardTokens(ctx context.Context, card *cards.Card) ([]card_dtos.CardToken, *errors.ErrorOutput) {

	var cardTokens []card_dtos.CardToken
	for _, gateway := range c.paymentGateways {

		gatewayName := gateway.GatewayName()

		cardToken, errorGetByCardIDAndGateway := c.cardTokenUseCaseGetByCardIDAndGateway.Execute(ctx, card.ID, gatewayName)

		if errorGetByCardIDAndGateway != nil && errorGetByCardIDAndGateway.Code != errors.CardNotFound {
			log.Info(ctx).Msgf("Ocorreu um erro ao pesquisar o CardToken para o CardID (%v) no Gateway (%v)", card.ID, gatewayName)
			return cardTokens, errorGetByCardIDAndGateway
		}

		if (errorGetByCardIDAndGateway != nil && errorGetByCardIDAndGateway.Code == errors.CardNotFound) || (cardToken == nil) {
			log.Info(ctx).Msgf("Não foi encontrado o CardToken para o CardID (%v) no Gateway (%v)", card.ID, gatewayName)
			continue
		}

		customerToken, errorGetCustomerToken := c.customerTokenUseCaseGetByCustomerIDAndGateway.Execute(ctx, card.CustomerID, gatewayName)

		if errorGetCustomerToken != nil {
			log.Info(ctx).Msgf("Ocorreu um erro ao pesquisar o CardToken para o CardID (%v) no Gateway (%v)", card.ID, gatewayName)
			return cardTokens, errorGetCustomerToken
		}

		deleteCardTokenUseCase := card_token_usecase.NewDeleteCardToken(
			c.cardTokenRepository,
			gateway,
		)

		output, errorDeleteCardToken := deleteCardTokenUseCase.Execute(ctx, customerToken.CustomerToken, cardToken.CardToken)

		if errorDeleteCardToken != nil {
			log.Error(ctx, errorDeleteCardToken.LogMessageToError()).Msgf("Erro ao deletar CardToken para o CardID (%v) no Gateway (%v), retornou o erro (%v)", card.ID, gatewayName, errorDeleteCardToken.Error())
			return cardTokens, errorDeleteCardToken
		}

		cardTokenDTO := card_dtos.CardToken{
			Token:   output.Token,
			Gateway: output.Gateway,
		}

		cardTokens = append(cardTokens, cardTokenDTO)

		log.Info(ctx).Msgf("Foi deletado o token de cartão com o CardID (%v) do CustomerID (%v) no Gateway (%v)",
			card.ID,
			card.CustomerID,
			gatewayName,
		)

	}
	return cardTokens, nil
}
