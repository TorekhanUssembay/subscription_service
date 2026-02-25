package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"subscription_service/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

// -------------------- CREATE --------------------

// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body service.CreateSubscriptionDTO true "Subscription DTO"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var dto service.CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.CreateSubscription(c.Request.Context(), dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// -------------------- GET --------------------

// @Summary Get subscription by ID
// @Description Get a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// -------------------- UPDATE --------------------

// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body service.CreateSubscriptionDTO true "Subscription DTO"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var dto service.CreateSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.service.UpdateSubscription(c.Request.Context(), id, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// -------------------- DELETE --------------------

// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// -------------------- LIST --------------------

// @Summary List subscriptions
// @Description List subscriptions filtered by user_id and optionally service_name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service Name"
// @Success 200 {array} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	serviceName := c.Query("service_name")
	var sn *string
	if serviceName != "" {
		sn = &serviceName
	}

	subs, err := h.service.ListSubscriptions(c.Request.Context(), userID, sn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// -------------------- SUM --------------------

// @Summary Sum subscriptions
// @Description Calculate total subscription price for a user in a period, optionally filtered by service_name
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service Name"
// @Param from query string true "From month-year MM-YYYY"
// @Param to query string true "To month-year MM-YYYY"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /subscriptions/sum [get]
func (h *SubscriptionHandler) SumSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	serviceName := c.Query("service_name")
	var sn *string
	if serviceName != "" {
		sn = &serviceName
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}

	sum, err := h.service.SumSubscriptions(c.Request.Context(), userID, sn, from, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sum": sum})
}