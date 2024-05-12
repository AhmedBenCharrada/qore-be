package person

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"qore-be/internal/domain/dto"
	"qore-be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Service the person services.
//
//go:generate mockery --name=Service --structname=PersonService  --case underscore --output=../mocks/ --filename=person_service.go
type Service interface {
	Create(context.Context, dto.PersonDTO) (dto.PersonDTO, error)
	GetByID(context.Context, int) (*dto.PersonDTO, error)
	GetAll(context.Context, int, int) ([]dto.PersonDTO, error)
}

// Controller represents the person controller.
type Controller struct {
	svc Service
	log *slog.Logger
}

// ControllerOption ..
type ControllerOption = func(*Controller) error

var (
	errMissingPersonService = fmt.Errorf("nil person service")
)

// NewController creates a new person controller.
func NewController(opts ...ControllerOption) (*Controller, error) {
	ctrl := &Controller{}

	for _, opt := range opts {
		if err := opt(ctrl); err != nil {
			return nil, err
		}
	}

	if ctrl.svc == nil {
		return nil, errMissingPersonService
	}

	if ctrl.log == nil {
		ctrl.log = slog.Default()
	}

	return ctrl, nil
}

// WithService initialize the person-controller wit a person-service.
func WithService(svc Service) ControllerOption {
	return func(ctrl *Controller) error {
		if svc == nil {
			return errMissingPersonService
		}
		ctrl.svc = svc
		return nil
	}
}

// Create represents the create a new person endpoint handler.
func (c *Controller) Create(ctx *gin.Context) {
	var req dto.PersonDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := c.svc.Create(ctx, req)
	if err != nil {
		c.log.Error("failed to create person", "error", err.Error(), "person", p)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.log.Info("person successfully created", "person", p)
	//ctx.JSON(http.StatusCreated, p)
	ctx.JSON(http.StatusOK, p)
}

func (c *Controller) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.log.Error("invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := c.svc.GetByID(ctx, id)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, p)
	case errors.Is(err, ErrRecordNotFound):
		c.log.Error("failed to get person data", "error", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.log.Error("failed to get person data", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// TODO.
func (c *Controller) GetAll(ctx *gin.Context) {
	page := utils.StringToInt(ctx.Query("page"), 0)
	limit := utils.StringToInt(ctx.Query("limit"), 25)
	persons, err := c.svc.GetAll(ctx, page, limit)
	if err != nil {
		c.log.Error("failed to get person data", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"content": persons,
		"page":    page,
		"size":    limit,
	})
}
