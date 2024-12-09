package card

import (
	"context"
	"encoding/json"
	"payment-layer-card-api/common/helpers"
	card_token_entities "payment-layer-card-api/entities/card_token"
	"payment-layer-card-api/entities/cards"
	payment_gateway_entities "payment-layer-card-api/entities/payment_gateway"
	"payment-layer-card-api/entities/queue"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	cards_usecase "payment-layer-card-api/usecases/cards"
	card_dtos "payment-layer-card-api/usecases/cards/dtos"
	customer_token_usecase "payment-layer-card-api/usecases/customer_token"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"payment-layer-card-api/workflows/customer"

	"github.com/adhfoundation/layer-tools/log"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"go.elastic.co/apm/v2"
)

type CreateCardWorkFlow struct {
	paymentGateways                               []payment_gateway_entities.PaymentGatewayInterface
	customersUseCaseGetCustomerByID               *customers_usecase.GetCustomerByID
	customerTokenUseCaseGetByCustomerIDAndGateway *customer_token_usecase.GetByCustomerIDAndGateway
	cardsUseCaseVerifyCardAlreadyExists           *cards_usecase.VerifyCardAlreadyExists
	cardsUseCaseCreateCard                        *cards_usecase.CreateCard
	cardTokenUseCaseGetByCardIDAndGateway         *card_token_usecase.GetCardTokenByCardIDAndGateway
	cardTokenRepository                           card_token_entities.CardTokenRepositoryInterface
	customerWorkflow                              *customer.CreateCustomerWorkflow
	cardRetryCount                                int
	cardRetryDelay                                int
	queueService                                  queue.QueueInterface
	queueUrl                                      string
}

func NewCreateCardWorkFlow(
	paymentGateways []payment_gateway_entities.PaymentGatewayInterface,
	customersUseCaseGetCustomerByID *customers_usecase.GetCustomerByID,
	customerTokenUseCaseGetByCustomerIDAndGateway *customer_token_usecase.GetByCustomerIDAndGateway,
	cardsUseCaseVerifyCardAlreadyExists *cards_usecase.VerifyCardAlreadyExists,
	cardsUseCaseCreateCard *cards_usecase.CreateCard,
	cardTokenUseCaseGetByCardIDAndGateway *card_token_usecase.GetCardTokenByCardIDAndGateway,
	cardTokenRepository card_token_entities.CardTokenRepositoryInterface,
	customerWorkflow *customer.CreateCustomerWorkflow,
	cardRetryCount int,
	cardRetryDelay int,
	queueService queue.QueueInterface,
	queueUrl string,

) *CreateCardWorkFlow {
	return &CreateCardWorkFlow{
		paymentGateways:                               paymentGateways,
		customersUseCaseGetCustomerByID:               customersUseCaseGetCustomerByID,
		customerTokenUseCaseGetByCustomerIDAndGateway: customerTokenUseCaseGetByCustomerIDAndGateway,
		cardsUseCaseVerifyCardAlreadyExists:           cardsUseCaseVerifyCardAlreadyExists,
		cardsUseCaseCreateCard:                        cardsUseCaseCreateCard,
		cardTokenUseCaseGetByCardIDAndGateway:         cardTokenUseCaseGetByCardIDAndGateway,
		cardTokenRepository:                           cardTokenRepository,
		cardRetryCount:                                cardRetryCount,
		cardRetryDelay:                                cardRetryDelay,
		customerWorkflow:                              customerWorkflow,
		queueService:                                  queueService,
		queueUrl:                                      queueUrl,
	}
}

func (c *CreateCardWorkFlow) Execute(ctx context.Context, input *card_dtos.CreateCardDTOInput) (*card_dtos.CreateCardDTOOutput, *errors.ErrorOutput) {
	if len(input.CVV) < 3 || len(input.CVV) > 4 {
		return nil, errors.NewError(errors.BadRequest, nil, "CVV precisa ter 3 ou 4 digitos")
	}

	if len(input.Number) < 13 || len(input.Number) > 19 {
		return nil, errors.NewError(errors.BadRequest, nil, "O número do cartão precisar ter entre 13 e 19 digitos")
	}

	dateValid, errDateValid := helpers.IsValidDate(input.ExpMonth, input.ExpYear)
	if errDateValid != nil {
		return nil, errDateValid
	}

	if !dateValid {
		return nil, errors.NewError(errors.BadRequest, nil, "Data inválida")
	}

	isExpired, errExpired := helpers.IsExpired(input.ExpMonth, input.ExpYear)
	if errExpired != nil {
		return nil, errExpired
	}

	if isExpired {
		return nil, errors.NewError(errors.BadRequest, nil, "Cartão está expirado")
	}

	if input.CustomerID == "" && input.Customer == nil {
		return nil, errors.NewError(errors.BadRequest, nil, "CustomerID ou Customer precisa ser informado")
	}

	if input.CustomerID == "" {
		inputToCreateCustomer := &customer_dto.CreateCustomerDTOInput{
			FirstName: input.Customer.FirstName,
			LastName:  input.Customer.LastName,
			Document:  input.Customer.Document,
			BirthDate: input.Customer.BirthDate,
			Gender:    input.Customer.Gender,
			Phone:     input.Customer.Phone,
			Email:     input.Customer.Email,
			ServiceID: input.Customer.ServiceID,
			Address: customer_dto.Address{
				State:      input.Customer.Address.State,
				City:       input.Customer.Address.City,
				Country:    input.Customer.Address.Country,
				Line1:      input.Customer.Address.Line1,
				Line2:      input.Customer.Address.Line2,
				District:   input.Customer.Address.District,
				PostalCode: input.Customer.Address.PostalCode,
			},
		}
		outputCreateCustomer, err := c.customerWorkflow.Execute(ctx, inputToCreateCustomer)
		if err != nil {
			return nil, err
		}
		input.CustomerID = outputCreateCustomer.ID
	}

	card, err := c.findOrCreateCard(ctx, input)
	if err != nil {
		return nil, err
	}

	cardTokens, errorCreateCardTokens := c.createCardTokens(ctx, card.ID, input)
	if errorCreateCardTokens != nil {
		return nil, errorCreateCardTokens
	}

	err_message := c.sendNotification(ctx, card, cardTokens)
	if err_message != nil {
		log.Error(ctx, err_message).Msgf("Ocorreu um erro ao enviar evento de criação de Cartão para fila.. Erro (%v)", err_message.Error())
	}

	cardOutput := &card_dtos.CreateCardDTOOutput{
		ID:         card.ID,
		CustomerID: card.CustomerID,
	}

	return cardOutput, nil
}

func (c *CreateCardWorkFlow) findOrCreateCard(ctx context.Context, input *card_dtos.CreateCardDTOInput) (*cards.Card, *errors.ErrorOutput) {
	cardExists, err := c.cardsUseCaseVerifyCardAlreadyExists.Execute(ctx, *input)

	if err != nil && err.Code != errors.NotFoundError {
		log.Error(ctx, err.LogMessageToError()).Msgf("Ocorreu um erro ao buscar o Cartão pelo Fingerprint para o CustomerID (%v). Erro (%v)", input.CustomerID, err)
		return nil, err
	}

	if cardExists != nil {
		log.Info(ctx).Msgf("O cartão com CardID (%v) foi encontrado para o CustomerID (%v)", cardExists.ID, input.CustomerID)
		return cardExists, nil
	}

	createCardDTOOutput, errorCreateCard := c.cardsUseCaseCreateCard.Execute(ctx, input)
	if errorCreateCard != nil {
		log.Error(ctx, errorCreateCard.LogMessageToError()).Msgf("Ocorreu um erro ao criar o Card do CustomerID (%v). Erro (%v)", input.CustomerID, errorCreateCard)
		return nil, errorCreateCard
	}

	log.Info(ctx).Msgf("O CardID (%v) foi criado para o CustomerID (%v)", createCardDTOOutput.ID, createCardDTOOutput.CustomerID)
	return createCardDTOOutput, nil
}

func (c *CreateCardWorkFlow) createCardTokens(ctx context.Context, cardID string, input *card_dtos.CreateCardDTOInput) ([]*card_token_dto.CardTokenOutputDTO, *errors.ErrorOutput) {
	countErrors := 0
	var lastError *errors.ErrorOutput
	var cardTokens []*card_token_dto.CardTokenOutputDTO

	for _, gateway := range c.paymentGateways {

		gatewayName := gateway.GatewayName()

		resultGetByCardIDAndGateway, errorGetByCardIDAndGateway := c.cardTokenUseCaseGetByCardIDAndGateway.Execute(ctx, cardID, gatewayName)

		if errorGetByCardIDAndGateway != nil && errorGetByCardIDAndGateway.Code != errors.NotFoundError {
			log.Info(ctx).Msgf("Ocorreu um erro ao pesquisar o CardToken para o CardID (%v) no Gateway (%v)", cardID, gatewayName)
			return nil, errorGetByCardIDAndGateway
		}

		if (errorGetByCardIDAndGateway != nil && errorGetByCardIDAndGateway.Code == errors.NotFoundError) || (resultGetByCardIDAndGateway == nil) {
			log.Info(ctx).Msgf("Não foi encontrado o CardToken para o CardID (%v) no Gateway (%v)", cardID, gatewayName)
		}

		if resultGetByCardIDAndGateway != nil {
			customerTokenId := resultGetByCardIDAndGateway.ID
			log.Info(ctx).Msgf("Foi encontrado o CardToken (%v) para o CardID (%v) no Gateway (%v)", customerTokenId, cardID, gatewayName)
			continue
		}

		customer, errorGetCustomer := c.customersUseCaseGetCustomerByID.Execute(ctx, input.CustomerID)

		if errorGetCustomer != nil {
			return nil, errorGetCustomer
		}

		customerToken, errorGetCustomerToken := c.customerTokenUseCaseGetByCustomerIDAndGateway.Execute(ctx, input.CustomerID, gatewayName)
		if errorGetCustomerToken != nil {
			return nil, errorGetCustomerToken
		}

		customerAddress := card_token_dto.HolderAddress{
			State:      input.BillingAddress.State,
			City:       input.BillingAddress.City,
			Country:    input.BillingAddress.Country,
			Line1:      input.BillingAddress.Line1,
			Line2:      input.BillingAddress.Line2,
			District:   input.BillingAddress.District,
			PostalCode: input.BillingAddress.PostalCode,
		}

		cardTokenInputDTO := &card_token_dto.CardTokenInputDTO{
			CardID:          cardID,
			CustomerID:      customer.ID,
			CustomerEmail:   customer.Email,
			CustomerPhone:   customer.Phone,
			CustomerAddress: &customerAddress,
			CustomerToken:   customerToken.CustomerToken,
			Holder:          input.Holder,
			HolderDocument:  customer.Document,
			Number:          input.Number,
			ExpMonth:        input.ExpMonth,
			ExpYear:         input.ExpYear,
			CVV:             input.CVV,
			Brand:           input.Brand,
		}

		createCardTokenUseCase := card_token_usecase.NewCreateCardToken(
			c.cardTokenRepository,
			gateway,
			c.cardRetryCount,
			c.cardRetryDelay,
		)

		createdToken, errorCreateCardToken := createCardTokenUseCase.Execute(ctx, cardID, cardTokenInputDTO)

		if errorCreateCardToken != nil {
			log.Error(ctx, errorCreateCardToken.LogMessageToError()).Msgf("Erro ao criar CardToken para o CostumerID (%v) no Gateway (%v), retornou o erro (%v)", cardID, gatewayName, errorCreateCardToken.Error())
			countErrors++
			lastError = errorCreateCardToken
		} else {
			log.Info(ctx).Msgf("Foi criado o token de cartão com o CardID (%v) do CustomerID (%v) no Gateway (%v)",
				cardID,
				input.CustomerID,
				gatewayName,
			)

			cardTokens = append(cardTokens, createdToken)
		}
	}

	if countErrors >= len(c.paymentGateways) {
		log.Info(ctx).Msgf("Ocorreu erro na busca ou criação de todos os CardTokens para o CustomerId (%v)", cardID)
		return nil, lastError
	}

	return cardTokens, nil
}

type CardMessage struct {
	ID          string              `json:"id"`
	CustomerID  string              `json:"customer_id"`
	Holder      string              `json:"holder"`
	Brand       string              `json:"brand"`
	Fingerprint string              `json:"fingerprint,omitempty"`
	FirstDigits string              `json:"first_digits"`
	LastDigits  string              `json:"last_digits"`
	ExpMonth    string              `json:"exp_month"`
	ExpYear     string              `json:"exp_year"`
	Tokens      []*CardTokenMessage `json:"tokens,omitempty"`
}

type CardTokenMessage struct {
	Gateway string `json:"gateway"`
	Token   string `json:"token"`
}

func (c *CreateCardWorkFlow) sendNotification(ctx context.Context, card *cards.Card, cardTokens []*card_token_dto.CardTokenOutputDTO) error {
	tx := apm.TransactionFromContext(ctx)
	if tx != nil {
		card.ApmLink.SpanId = tx.TraceContext().Span
		card.ApmLink.TraceId = tx.TraceContext().Trace
	}

	cardMessage := CardMessage{
		ID:          card.ID,
		CustomerID:  card.CustomerID,
		Holder:      card.Holder,
		Brand:       card.Brand,
		Fingerprint: card.Fingerprint,
		FirstDigits: card.FirstDigits,
		LastDigits:  card.LastDigits,
		ExpMonth:    card.ExpMonth,
		ExpYear:     card.ExpYear,
	}

	for _, token := range cardTokens {
		tokenMessage := CardTokenMessage{
			Gateway: token.Gateway,
			Token:   token.Token,
		}
		cardMessage.Tokens = append(cardMessage.Tokens, &tokenMessage)
	}

	createCardSQSMessage, err := json.Marshal(cardMessage)
	if err != nil {
		log.Error(ctx, err).Msgf("Ocorreu um erro ao transformar cartão em JSON. Erro (%v)", err)
		return err
	}

	err = c.queueService.SendMessageWithContext(ctx, c.queueUrl, string(createCardSQSMessage))
	if err != nil {
		return errors.NewError(errors.QueueMessageError, err)
	}

	return nil
}
