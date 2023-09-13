package common

const (
	HashCost                  = 12
	TimeoutInSecondsOfSession = 30
)

const (
	// DataExist Mapping data existed.
	DataExist = 520
	// DoesNotExist Mapping data does not exist
	DoesNotExist = 521
	// InvalidStatusCode error code invalid status
	InvalidStatusCode = 523
	// NotFoundCode error code record not found
	NotFoundCode = 524
)

const VI_LANG = "vi"

const (
	StatusActive   = "Active"
	StatusInactive = "Inactive"
)

const (
	StatusTxPending = "Pending"
	StatusTxSuccess = "Active"
	StatusTxFailure = "Failed"
)

const (
	MessageErrorEmailAlreadyUsed     = "email already used"
	MessageErrorInvalidToken         = "invalid token provided"
	MessageErrorExistedEmail         = "email already exists"
	MessageErrorOrgNotFound          = "organization not found"
	MessageErrorExistedMapping       = "chip or productItem is already mapped"
	MessageErrorExistedOrganization  = "this organization's name tage is already taken"
	MessageErrorExistedTag           = "this tag is already added"
	MessageErrorWrongFormat          = "wrong format provided"
	MessageErrorExistedProductName   = "the product name already exists"
	MessageErrorNotAbleToClaim       = "product is not available to claim"
	MessageErrorNotFoundUser         = "user is not registered yet!"
	MessageErrorNotFoundOrganization = "organization's name tag doesn't exist"
	MessageErrorCreateOrgFail        = "error on creating organization, tagname is taken"
	MessageErrorCreateTagFail        = "error on creating tag!"
	MessageErrorCreateTemplateFail   = "error creating template!"
	MessageErrorUpdateTemplateFail   = "error creating template!"
	MessageErrorFailDetectUser       = "cannot detect user information"
	MessageErrorInvalidEntityID      = "invalid entity id provided"
	MessageErrorInvalidTemplateID    = "invalid template id provided"
	MessageErrorInvalidOrgTagName    = "Organization Tag Name has invalid characters (only allow a-z (lowercase characters), A-Z (uppercase characters), 0-9 (number), - (hyphen), _ (underscore))"
)
