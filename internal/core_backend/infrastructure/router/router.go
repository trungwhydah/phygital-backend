package router

import (
	"net/http"

	"backend-service/internal/core_backend/api/handler"
	"backend-service/internal/core_backend/api/middleware"
	_ "backend-service/internal/core_backend/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Phygital Core Backend API
//	@version		1.0
//	@description	A book management service API in Go using Gin framework..
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://fygito.com//support
//	@contact.email	support@kardialab.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for what is this security definition being used
func Initialize(handler handler.AppHandler, mdw middleware.MidddlewareServices) {
	router := gin.New()
	router.Use(corsMiddleware())
	router.Use(gin.Logger())

	router.GET("/dummy", func(c *gin.Context) {
		result := handler.DummyHandler.GetDummy(c)
		c.JSON(result.Code, result)
	})

	router.POST("/signup", func(c *gin.Context) {
		result := handler.UserHandler.UserSignUp(c)
		c.JSON(result.Code, result)
	})

	// Refactor Part
	// Support
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	router.GET("/swagger/*any", swaggerHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	router.POST("/upload", func(c *gin.Context) {
		result := handler.UploadHandler.Upload(c)
		c.JSON(result.Code, result)
	})

	// End-user part
	scanGroup := router.Group("/scan")
	{
		scanGroup.GET("/nfc/verify", func(c *gin.Context) {
			result := handler.ScanHandler.DecodeScan(c)
			c.Redirect(http.StatusFound, result.URL)
		})
		scanGroup.GET("/nfc/tap/:tag_id", func(c *gin.Context) {
			result := handler.ScanHandler.Tap(c)
			c.Redirect(http.StatusFound, result.URL)
		})
	}
	sessionGroup := router.Group("/session")
	{
		sessionGroup.GET("/verify", func(c *gin.Context) {
			result := handler.SessionHandler.VerifySession(c)
			c.JSON(result.Code, result)
		})
	}

	product := router.Group("/product")
	{
		product.GET("/seo", func(c *gin.Context) {
			result := handler.ProductHandler.GetProductByTagID(c)
			c.JSON(result.Code, result)
		})
	}
	productItem := router.Group("/product-item")
	{
		productItem.GET("/story", func(c *gin.Context) {
			result := handler.ProductItemHandler.GetStoryByTagID(c)
			c.JSON(result.Code, result)
		})
		productItem.GET("/:product_item_id/detail", func(c *gin.Context) {
			result := handler.ProductItemHandler.GetDetailProductItem(c)
			c.JSON(result.Code, result)
		})
		productItem.PUT("/:product_item_id/claim", mdw.AuthenMiddleware.UserAuth.Authenticate, func(c *gin.Context) {
			result := handler.ProductItemHandler.ClaimItem(c)
			c.JSON(result.Code, result)
		})
		productItem.PUT("/:product_item_id/toggle-claimable", func(c *gin.Context) {
			result := handler.ProductItemHandler.ToggleClaimableItem(c)
			c.JSON(result.Code, result)
		})
		productItem.POST("/:product_item_id/like", func(c *gin.Context) {
			result := handler.ProductItemHandler.LikeProductItem(c)
			c.JSON(result.Code, result)
		})
		//Test: NFT Metatdata
		productItem.GET("/metadata/:token_id", func(c *gin.Context) {
			result := handler.ProductItemHandler.GetMetadata(c)
			c.JSON(result.Code, result.Result)
		})
	}

	// Authenticate Part - authenticate and authorization required
	adminGroup := router.Group("/admin")
	adminGroup.Use(mdw.AuthenMiddleware.AdminAuth.Authenticate)
	{
		productGroup := adminGroup.Group("/product")
		{
			productGroup.GET("", func(c *gin.Context) {
				result := handler.ProductHandler.GetAllProducts(c)
				c.JSON(result.Code, result)
			})
			productGroup.POST("/create", func(c *gin.Context) {
				result := handler.ProductHandler.CreateProduct(c)
				c.JSON(result.Code, result)
			})
			productGroup.GET("/:product_id", func(c *gin.Context) {
				result := handler.ProductHandler.GetProductDetail(c)
				c.JSON(result.Code, result)
			})
			productGroup.PUT("/:product_id", func(c *gin.Context) {
				result := handler.ProductHandler.UpdateProductDetail(c)
				c.JSON(result.Code, result)
			})
			productGroup.DELETE("/:product_id", func(c *gin.Context) {
				result := handler.ProductHandler.DeteleProductByID(c)
				c.JSON(result.Code, result)
			})
		}

		mappingGroup := adminGroup.Group("/mapping")
		{
			mappingGroup.GET("", func(c *gin.Context) {
				result := handler.MappingHandler.GetAllMapping(c)
				c.JSON(result.Code, result)
			})
			mappingGroup.GET("/product/:product_id", func(c *gin.Context) {
				result := handler.MappingHandler.GetAllMappingForProduct(c)
				c.JSON(result.Code, result)
			})
			mappingGroup.PUT("/:tag_id", func(c *gin.Context) {
				result := handler.MappingHandler.UpdateMapping(c)
				c.JSON(result.Code, result)
			})
			mappingGroup.DELETE("/:tag_id", func(c *gin.Context) {
				result := handler.MappingHandler.Unmap(c)
				c.JSON(result.Code, result)
			})
			mappingGroup.POST("/batch/multiple-mapping", func(c *gin.Context) {
				result := handler.MappingHandler.MultipleMappingWithSingleProduct(c)
				c.JSON(result.Code, result)
			})
		}

		businessProductItem := adminGroup.Group("/product-item")
		{
			businessProductItem.POST("/:product_item_id/mint", func(c *gin.Context) {
				result := handler.ProductItemHandler.MintProductItem(c)
				c.JSON(result.Code, result)
			})
			businessProductItem.POST("/create", func(c *gin.Context) {
				result := handler.ProductItemHandler.CreateProductItem(c)
				c.JSON(result.Code, result)
			})
			businessProductItem.POST("/create-multiple", func(c *gin.Context) {
				result := handler.ProductItemHandler.CreateMultipleProductItems(c)
				c.JSON(result.Code, result)
			})
			businessProductItem.GET("", func(c *gin.Context) {
				result := handler.ProductItemHandler.GetAllProductItem(c)
				c.JSON(result.Code, result)
			})
			businessProductItem.GET("/organization/:org_tag_name", func(c *gin.Context) {
				result := handler.ProductItemHandler.GetAllProductItemInOrg(c)
				c.JSON(result.Code, result)
			})
		}

		organizationGroup := adminGroup.Group("/organization")
		{
			organizationGroup.PUT("/:org_id", func(c *gin.Context) {
				result := handler.OrganizationHandler.UpdateOrganization(c)
				c.JSON(result.Code, result)
			})
			organizationGroup.POST("/create", func(c *gin.Context) {
				result := handler.OrganizationHandler.CreateOrganization(c)
				c.JSON(result.Code, result)
			})
			organizationGroup.GET("", func(c *gin.Context) {
				result := handler.OrganizationHandler.GetAllOrganizations(c)
				c.JSON(result.Code, result)
			})
			organizationGroup.GET("/:org_tag_name", func(c *gin.Context) {
				result := handler.OrganizationHandler.GetOrganization(c)
				c.JSON(result.Code, result)
			})
		}
		templateGroup := adminGroup.Group("template")
		{
			templateGroup.GET("/:template_id", func(c *gin.Context) {
				result := handler.TemplateHandler.GetTemplate(c)
				c.JSON(http.StatusOK, result)
			})
			templateGroup.PUT("/:template_id", func(c *gin.Context) {
				result := handler.TemplateHandler.UpdateTemplate(c)
				c.JSON(http.StatusOK, result)
			})
			templateGroup.GET("/all", func(c *gin.Context) {
				result := handler.TemplateHandler.GetAllTemplates(c)
				c.JSON(http.StatusOK, result)
			})
			templateGroup.POST("/create", func(c *gin.Context) {
				result := handler.TemplateHandler.CreateTemplate(c)
				c.JSON(http.StatusOK, result)
			})
		}
		pageGroup := adminGroup.Group("web-page")
		{
			pageGroup.GET("/all", func(c *gin.Context) {
				result := handler.WebPageHandler.GetAllWebPages(c)
				c.JSON(200, result)
			})
			pageGroup.POST("/create", func(c *gin.Context) {
				result := handler.WebPageHandler.CreateWebPage(c)
				c.JSON(result.Code, result)
			})
			pageGroup.PUT("/:webpage_id", func(c *gin.Context) {
				result := handler.WebPageHandler.UpdateWebPage(c)
				c.JSON(result.Code, result)
			})
			pageGroup.DELETE("/:webpage_id", func(c *gin.Context) {
				result := handler.WebPageHandler.DeleteWebPage(c)
				c.JSON(http.StatusOK, result)
			})
		}
		tagManagementGroup := adminGroup.Group("/tag")
		{
			tagManagementGroup.POST("/create", func(c *gin.Context) {
				result := handler.TagHandler.CreateTag(c)
				c.JSON(result.Code, result)

			})
			tagManagementGroup.POST("/create/batch", func(c *gin.Context) {
				// Create multiple chip: chip number: from - to for organization
			})
		}
		userGroup := adminGroup.Group("/user")
		{
			userGroup.PUT("/role", func(c *gin.Context) {
				result := handler.UserHandler.UpdateRole(c)
				c.JSON(result.Code, result)
			})
			userGroup.PUT("/sync-wallet-address", func(c *gin.Context) {
				result := handler.UserHandler.SyncWalletAddress(c)
				c.JSON(result.Code, result)
			})
		}

		digitalAssetGroup := adminGroup.Group("/digital-asset")
		{
			digitalAssetGroup.GET("", func(c *gin.Context) {
				result := handler.DigitalAssetHandler.GetDigitalAssets(c)
				c.JSON(result.Code, result.Result)
			})
			digitalAssetGroup.PUT("/sync-metadata", func(c *gin.Context) {
				result := handler.DigitalAssetHandler.SyncDigitalAssetsMetadata(c)
				c.JSON(result.Code, result.Result)
			})
		}

		authorGroup := adminGroup.Group("/author")
		{
			authorGroup.GET("", func(c *gin.Context) {
				result := handler.AuthorHandler.GetListAuthor(c)
				c.JSON(result.Code, result.Result)
			})
			authorGroup.POST("", func(c *gin.Context) {
				result := handler.AuthorHandler.CreateAuthor(c)
				c.JSON(result.Code, result.Result)
			})
		}
	}

	publicAuthorGroup := router.Group("/author")
	{
		publicAuthorGroup.GET("/:author_id", func(c *gin.Context) {
			result := handler.AuthorHandler.GetDetailAuthor(c)
			c.JSON(result.Code, result.Result)
		})
	}

	digitalAssetGroup := router.Group("/digital-asset")
	{
		digitalAssetGroup.GET("collection/:collectionID", func(c *gin.Context) {
			result := handler.DigitalAssetHandler.GetDigitalAssetsByCollection(c)
			c.JSON(result.Code, result.Result)
		})

		digitalAssetGroup.GET(":org_tag_name/:token_id", func(c *gin.Context) {
			result := handler.DigitalAssetHandler.GetDigitalMetadataWithID(c)
			c.JSON(result.Code, result.Result)
		})
	}
	webpageGroup := router.Group("/web-page")
	{
		webpageGroup.GET("/:webpage_id", func(c *gin.Context) {
			result := handler.WebPageHandler.GetWebPage(c)
			c.JSON(http.StatusOK, result)
		})
	}

	pubsubGroup := router.Group("/pubsub")
	pubsubGroup.Use(mdw.AuthenMiddleware.PubsubAuth.Authenticate)
	{
		pubsubGroup.POST("/upsert-user", func(c *gin.Context) {
			result := handler.PubsubHandler.UpsertUser(c)
			c.JSON(result.Code, result)
		})
	}

	userGroup := router.Group("/user")
	userGroup.Use(mdw.AuthenMiddleware.UserAuth.Authenticate)
	{
		userGroup.PUT("/org", func(c *gin.Context) {
			result := handler.UserHandler.UpdateOrg(c)
			c.JSON(result.Code, result)
		})
		userGroup.PUT("", func(c *gin.Context) {
			result := handler.UserHandler.UpdateUserDetails(c)
			c.JSON(result.Code, result)
		})
		userGroup.GET("", func(c *gin.Context) {
			result := handler.UserHandler.GetUserDetails(c)
			c.JSON(result.Code, result)
		})
	}

	competitionGroup := router.Group("/competition")
	{
		competitionGroup.GET("/:org_tag_name", func(c *gin.Context) {
			result := handler.ProductItemHandler.GetGalleryOfProductItemsInOrg(c)
			c.JSON(result.Code, result)
		})
		competitionGroup.GET("/v2/:org_tag_name", func(c *gin.Context) {
			result := handler.ProductItemHandler.GetGalleryOfProductItemsInOrgV2(c)
			c.JSON(result.Code, result)
		})
	}
	router.Run(":8080")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
