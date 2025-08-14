package routers

import (
	providerControllers "provider-report-api/internal/modules/provider-detail/controllers"
	"provider-report-api/internal/modules/provider-detail/repositories"
	"provider-report-api/internal/modules/provider-detail/services"

	"github.com/gin-gonic/gin"
)

// Dependencies struct to hold all required dependencies
type Dependencies struct {
	ProviderService *services.ProviderService
	TemplateService *services.TemplateService
	ScheduleService *services.ScheduleService
	LogService      *services.LogService
	FieldRepo       *repositories.FieldRepository
}

func InitializeRoutes(r *gin.Engine, deps *Dependencies) {
	// Define global prefix for all routes
	globalRoute := r.Group("/tpa-api")

	reportRoute := globalRoute.Group("/report")

	// Define group path for provider-detail module
	providerRoute := reportRoute.Group("/provider-detail")
	{
		providerControllers.InitializeProviderDetailController(
			providerRoute,
			deps.ProviderService,
			deps.TemplateService,
			deps.ScheduleService,
			deps.LogService,
			deps.FieldRepo,
		)
	}
}