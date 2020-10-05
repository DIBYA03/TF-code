package bank

type ProviderName string

func (n ProviderName) String() string {
	return string(n)
}

const (
	// BBVA
	ProviderNameBBVA = ProviderName("bbva")

	// Plaid
	ProviderNamePlaid = ProviderName("plaid")

	// Stripe
	ProviderNameStripe = ProviderName("stripe")
)

var consumerBanks = map[ProviderName]ConsumerServiceProvider{}
var businessBanks = map[ProviderName]BusinessServiceProvider{}
var proxyBanks = map[ProviderName]ProxyServiceProvider{}

func AddBank(name ProviderName, c ConsumerServiceProvider, b BusinessServiceProvider, p ProxyServiceProvider) {
	consumerBanks[name] = c
	businessBanks[name] = b
	proxyBanks[name] = p
}

func GetConsumerBank(name ProviderName) (ConsumerServiceProvider, error) {
	bank, ok := consumerBanks[name]
	if !ok {
		return nil, NewErrorFromCode(ErrorCodeInvalidBankProvider)
	}

	return bank, nil
}

func GetBusinessBank(name ProviderName) (BusinessServiceProvider, error) {
	bank, ok := businessBanks[name]
	if !ok {
		return nil, NewErrorFromCode(ErrorCodeInvalidBankProvider)
	}

	return bank, nil
}

func GetProxyBank(name ProviderName) (ProxyServiceProvider, error) {
	bank, ok := proxyBanks[name]
	if !ok {
		return nil, NewErrorFromCode(ErrorCodeInvalidBankProvider)
	}

	return bank, nil
}

type ConsumerServiceProvider interface {
	ConsumerEntityService(APIRequest) ConsumerService
	BankAccountService(APIRequest, ConsumerID) (ConsumerBankAccountService, error)
	LinkedAccountService(APIRequest, ConsumerID) (ConsumerLinkedAccountService, error)
	LinkedCardService(APIRequest, ConsumerID) (ConsumerLinkedCardService, error)
	CardService(APIRequest, ConsumerID) (ConsumerCardService, error)
	MoneyTransferService(APIRequest, ConsumerID) (ConsumerMoneyTransferService, error)
}

type BusinessServiceProvider interface {
	BusinessEntityService(APIRequest) BusinessService
	BankAccountService(APIRequest, BusinessID) (BusinessBankAccountService, error)
	LinkedAccountService(APIRequest, BusinessID) (BusinessLinkedAccountService, error)
	LinkedCardService(APIRequest, BusinessID) (BusinessLinkedCardService, error)
	CardService(APIRequest, BusinessID, ConsumerID) (BusinessCardService, error)
	MoneyTransferService(APIRequest, BusinessID) (BusinessMoneyTransferService, error)
	LinkedPayeeService(APIRequest, BusinessID) (BusinessLinkedPayeeService, error)
}

type ProxyServiceProvider interface {
	ProxyService(APIRequest) ProxyService
}

type ConsumerService interface {
	/* Create a Consumer */
	Create(CreateConsumerRequest) (*IdentityStatusConsumerResponse, error)

	/* Update a Consumer */
	Update(UpdateConsumerRequest) (*IdentityStatusConsumerResponse, error)

	/* Get by user id
	Get(ConsumerID) (*ConsumerResponse, error) */

	/* Get Identity Status */
	Status(ConsumerID) (*IdentityStatusConsumerResponse, error)

	/* Verify an Identity */
	UploadIdentityDocument(ConsumerIdentityDocumentRequest) (*IdentityDocumentResponse, error)

	/* Update phone or email */
	UpdateContact(ConsumerID, ConsumerPropertyType, string) error

	/* Update address */
	UpdateAddress(ConsumerID, ConsumerPropertyType, AddressRequest) error

	/* Delete consumer */
	Delete(ConsumerID) error
}

type BusinessService interface {
	/* Create a business */
	Create(CreateBusinessRequest) (*IdentityStatusBusinessResponse, error)

	/* Update a Business */
	Update(UpdateBusinessRequest) (*IdentityStatusBusinessResponse, error)

	/* Get by member id
	GetMember(BusinessID, ConsumerID) (*BusinessMemberResponse, error) */

	/* Upload identity document */
	UploadIdentityDocument(BusinessIdentityDocumentRequest) (*IdentityDocumentResponse, error)

	/* Get Identity Status */
	Status(BusinessID) (*IdentityStatusBusinessResponse, error)

	/* Update phone or email */
	UpdateContact(BusinessID, BusinessPropertyType, string) error

	/* Update address */
	UpdateAddress(BusinessID, BusinessPropertyType, AddressRequest) error
}

type ConsumerBankAccountService interface {
	Create(CreateConsumerBankAccountRequest) (*CreateConsumerBankAccountResponse, error)

	Get(AccountBankID) (*GetBankAccountResponse, error)
	Patch(AccountBankID, PatchBankAccountRequest) (*GetBankAccountResponse, error)
	Close(AccountBankID, AccountCloseReason) (*GetBankAccountResponse, error)

	GetParticipants(AccountBankID) ([]AccountParticipantResponse, error)
	AddParticipants(AccountBankID, []AccountParticipantRequest) ([]AccountParticipantResponse, error)
	AddParticipant(AccountBankID, AccountParticipantRequest) (*AccountParticipantResponse, error)
	RemoveParticipant(AccountBankID, ConsumerID) error

	Block(AccountBlockRequest) (*AccountBlockResponse, error)
	Unblock(AccountUnblockRequest) error
	GetAllBlocks(AccountBankID) ([]AccountBlockResponse, error)

	GetStatementByID(AccountBankID, AccountStatementBankID) (*GetAccountStatementDocument, error)
	GetStatements(AccountBankID) ([]AccountStatementResponse, error)
}

type BusinessBankAccountService interface {
	Create(CreateBusinessBankAccountRequest) (*CreateBusinessBankAccountResponse, error)

	Get(AccountBankID) (*GetBankAccountResponse, error)
	Patch(AccountBankID, PatchBankAccountRequest) (*GetBankAccountResponse, error)
	Close(AccountBankID, AccountCloseReason) (*GetBankAccountResponse, error)

	GetParticipants(AccountBankID) ([]AccountParticipantResponse, error)
	AddParticipants(AccountBankID, []AccountParticipantRequest) ([]AccountParticipantResponse, error)
	AddParticipant(AccountBankID, AccountParticipantRequest) (*AccountParticipantResponse, error)
	RemoveParticipant(AccountBankID, ConsumerID) error

	Block(AccountBlockRequest) (*AccountBlockResponse, error)
	Unblock(AccountUnblockRequest) error
	GetAllBlocks(AccountBankID) ([]AccountBlockResponse, error)

	GetStatementByID(AccountBankID, AccountStatementBankID) (*GetAccountStatementDocument, error)
	GetStatements(AccountBankID) ([]AccountStatementResponse, error)
}

type LinkedAccountService interface {
	Link(*LinkedBankAccountRequest) (*LinkedBankAccountResponse, error)
	Get(LinkedAccountBankID) (*LinkedBankAccountResponse, error)
	GetAll() ([]LinkedBankAccountResponse, error)
	Unlink(LinkedAccountBankID) error
}

type ConsumerLinkedAccountService LinkedAccountService
type BusinessLinkedAccountService LinkedAccountService

type LinkedCardService interface {
	Link(*LinkedCardRequest) (*LinkedCardResponse, error)
	Get(LinkedCardBankID) (*LinkedCardResponse, error)
	GetAll() ([]LinkedCardResponse, error)
	Unlink(LinkedCardBankID) error
}

type ConsumerLinkedCardService LinkedCardService
type BusinessLinkedCardService LinkedCardService

type BusinessLinkedPayeeService interface {
	Link(*LinkedPayeeRequest) (*LinkedPayeeResponse, error)
	Get(BankPayeeID) (*LinkedPayeeResponse, error)
	Unlink(BankPayeeID) error
}

type CardService interface {
	Create(CreateCardRequest) (*GetCardResponse, error)
	Get(CardBankID) (*GetCardResponse, error)
	GetAll() ([]GetCardResponse, error)
	GetLimit(CardBankID) (*GetCardLimitResponse, error)

	Reissue(ReissueCardRequest) (*GetCardResponse, error)
	Activate(ActivateCardRequest) (*GetCardResponse, error)
	SetPIN(SetCardPINRequest) error
	Cancel(CancelCardRequest) error
	CancelInternal(CardBankID) error

	Block(CardBlockRequest) ([]CardBlockResponse, error)
	Unblock(CardUnblockRequest) error
	GetAllBlocks(CardBankID) ([]CardBlockResponse, error)
}

type ConsumerCardService CardService
type BusinessCardService CardService

type MoneyTransferService interface {
	Submit(*MoneyTransferRequest) (*MoneyTransferResponse, error)
	Get(MoneyTransferBankID) (*MoneyTransferResponse, error)
	GetAll() ([]MoneyTransferResponse, error)
	Cancel(MoneyTransferBankID) (*MoneyTransferResponse, error)
}

type ConsumerMoneyTransferService MoneyTransferService
type BusinessMoneyTransferService MoneyTransferService

type ProxyService interface {
	// Gets bank access token
	GetAccessToken() (*string, error)

	// Get base api url
	GetBaseAPIURL() string

	// Get bank id for a consumer
	GetConsumerBankID(ConsumerID) (*ConsumerBankID, error)

	// Get consumer id for a consumer
	GetConsumerID(ConsumerBankID) (*ConsumerID, error)

	// Get bank id for a business
	GetBusinessBankID(BusinessID) (*BusinessBankID, error)

	// Get business ID using bank ID
	GetBusinessID(BusinessBankID) (*BusinessID, error)
}
