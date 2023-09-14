package handler

import (
	"backend-service/internal/core_backend/api/handler/request"
	"backend-service/internal/core_backend/api/presenter"
	"backend-service/internal/core_backend/entity"
	"backend-service/internal/core_backend/usecase/author"
	"backend-service/internal/core_backend/usecase/product"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthorHandler interface
type AuthorHandler interface {
	CreateAuthor(*gin.Context) APIResponse
	GetListAuthor(*gin.Context) APIResponse
	GetDetailAuthor(*gin.Context) APIResponse
	UpdateAuthor(*gin.Context) APIResponse
}

// authorHandler struct
type authorHandler struct {
	AuthorService   author.UseCase
	ProductService  product.UseCase
	AuthorPresenter presenter.ConvertAuthor
}

// NewAuthorHandler create handler
func NewAuthorHandler(tuc author.UseCase, puc product.UseCase, pap presenter.ConvertAuthor) AuthorHandler {
	return &authorHandler{
		AuthorService:   tuc,
		ProductService:  puc,
		AuthorPresenter: pap,
	}
}

// CreateAuthor	godoc
// CreateAuthor	API
//
//	@Summary		Create Author
//	@Description	Create author
//	@Tags			author
//	@Accept      	json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/author [post]
//	@Param			request body entity.Author	true	"Create an author"
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		500					{object}	APIResponse
func (h *authorHandler) CreateAuthor(c *gin.Context) APIResponse {
	var requestAuthor entity.Author
	if err := c.ShouldBind(&requestAuthor); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	authorInserted, code, err := h.AuthorService.CreateAuthor(&requestAuthor)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", authorInserted)
}

// GetListAuthor	godoc
// GetListAuthor	API
//
//	@Summary		Get List of Author
//	@Description	Get List of author
//	@Tags			author
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/author [get]
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		500					{object}	APIResponse
func (h *authorHandler) GetListAuthor(c *gin.Context) APIResponse {
	listAuthors, code, err := h.AuthorService.GetListAuthor()
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", listAuthors)
}

// GetDetailAuthor	godoc
// GetDetailAuthor	API
//
//	@Summary		Get Author information
//	@Description	Get Author information
//	@Tags			author
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/author/{author_id} [get]
//	@Param			author_id	path		string	true	"Author ID"
//	@Success		200					{object}	APIResponse{result=bool}
//	@Failure		500					{object}	APIResponse
func (h *authorHandler) GetDetailAuthor(c *gin.Context) APIResponse {
	var req request.AuthorRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	author, code, err := h.AuthorService.GetAuthorDetail(&req.AuthorID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	authorID := author.ID.Hex()
	listProduct, code, err := h.ProductService.GetProductForAuthor(&authorID)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", h.AuthorPresenter.ResponseAuthorDetail(author, listProduct))
}

// UpdateAuthor	godoc
// UpdateAuthor	API
//
//	@Summary		Update Author
//	@Description	Update author
//	@Tags			author
//	@Accept      	json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Router			/admin/author/{author_id} [put]
//	@Param			author_id	path		string	true	"Author ID"
//	@Param			request body entity.Author	true	"Update an author"
//	@Success		200					{object}	APIResponse{result=entity.Author}
//	@Failure		500					{object}	APIResponse
func (h *authorHandler) UpdateAuthor(c *gin.Context) APIResponse {
	var req request.AuthorRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}
	var requestAuthor entity.Author
	if err := c.ShouldBind(&requestAuthor); err != nil {
		return CreateResponse(err, http.StatusBadRequest, "", err.Error(), nil)
	}

	updatedAuthor, code, err := h.AuthorService.UpdateAuthor(&req.AuthorID, &requestAuthor)
	if err != nil {
		return CreateResponse(err, code, "", err.Error(), nil)
	}

	return HandlerResponse(code, "", "", updatedAuthor)
}
