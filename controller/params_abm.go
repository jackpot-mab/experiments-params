package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"jackpot-mab/experiments-params/db"
	"jackpot-mab/experiments-params/model"
	"log"
	"net/http"
)

type ExperimentParamsController struct {
	DAO db.ExperimentsDAO
}

// @BasePath /api/v1

// GetExperiments godoc
// @Summary return data from all experiments
// @Schemes
// @Description Consults configuration DB to fetch experiments data.
// @Tags experiments-params
// @Accept json
// @Produce json
// @Success 200 {string} []Experiment
// @Router /experiments [get]
func (e *ExperimentParamsController) GetExperiments(g *gin.Context) {
	g.JSON(http.StatusOK, e.DAO.GetAllExperiments())
}

// GetExperiment godoc
// @Summary return experiment data
// @Schemes
// @Description Consults configuration DB to fetch experiment data.
// @Tags experiments-params
// @Accept json
// @Produce json
// @Param experiment_id path string true "Experiment ID"
// @Success 200 {string} Experiment
// @Router /experiment/{experiment_id} [get]
func (e *ExperimentParamsController) GetExperiment(g *gin.Context) {
	experimentId := g.Param("experiment_id")
	g.JSON(http.StatusOK, e.DAO.GetExperiment(experimentId))
}

// UpdateExperiment godoc
// @Summary updates experiment data and return experiment
// @Schemes
// @Description Update configuration db with experiment data.
// @Tags experiments-params
// @Accept json
// @Param experiment body model.Experiment true "Experiment data"
// @Produce json
// @Success 200 {string} Experiment
// @Router /experiment [put]
func (e *ExperimentParamsController) UpdateExperiment(g *gin.Context) {

}

// AddExperiment godoc
// @Summary creates experiment data and return experiment
// @Schemes
// @Description Create experiment in configuration db.
// @Tags experiments-params
// @Accept json
// @Param experiment body model.Experiment true "Experiment data"
// @Produce json
// @Success 200 {string} Experiment
// @Router /experiment [post]
func (e *ExperimentParamsController) AddExperiment(g *gin.Context) {
	var experiment model.Experiment
	if err := g.BindJSON(&experiment); err != nil {
		log.Print("error occurred", err)
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !model.Validate(experiment) {
		log.Print("invalid names, : is a forbidden character")
		g.JSON(http.StatusBadRequest, errors.New("invalid experiment or branch name, "+
			"it contains invalid characters").Error())
		return
	}

	g.JSON(http.StatusOK, e.DAO.AddExperiment(experiment))
}

// AddOrUpdateParameter godoc
// @Summary add a parameter to an experiment/arm or update the parameter value if already exists.
// @Schemes
// @Description add or create experiment/arm parameter
// @Tags experiments-params
// @Accept json
// @Param experiment body model.RewardDataParameterUpsert true "Reward data parameter"
// @Produce json
// @Success 200 {string} error
// @Router /experiment/parameter [post]
func (e *ExperimentParamsController) AddOrUpdateParameter(g *gin.Context) {
	var rewardDataParameter model.RewardDataParameterUpsert
	if err := g.BindJSON(&rewardDataParameter); err != nil {
		log.Print("error occurred", err)
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := e.DAO.AddOrUpdateRewardParameter(rewardDataParameter)

	if result == nil {
		g.JSON(http.StatusOK, "updated ok")
	}

	g.JSON(http.StatusInternalServerError, result.Error())

}
