package controller

import "github.com/gin-gonic/gin"

// @BasePath /api/v1

// ExperimentsParamsController godoc
// @Summary return experiment data
// @Schemes
// @Description Consults configuration DB to fetch experiment data.
// @Tags experiments-params
// @Accept json
// @Produce json
// @Success 200 {string} Experiment
// @Router /experiment [get]
func ExperimentsParamsController(g *gin.Context) {

}

// UpdateExperimentParametersController godoc
// @Summary updates experiment data and return experiment
// @Schemes
// @Description Update configuration db with experiment data.
// @Tags experiments-params
// @Accept json
// @Param experiment body model.Experiment true "Experiment data"
// @Produce json
// @Success 200 {string} Experiment
// @Router /experiment [put]
func UpdateExperimentParametersController(g *gin.Context) {

}

// AddExperimentParametersController godoc
// @Summary creates experiment data and return experiment
// @Schemes
// @Description Create experiment in configuration db.
// @Tags experiments-params
// @Accept json
// @Param experiment body model.Experiment true "Experiment data"
// @Produce json
// @Success 200 {string} Experiment
// @Router /experiment [post]
func AddExperimentParametersController(g *gin.Context) {}
