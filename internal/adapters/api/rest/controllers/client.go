package controllers

import (
	"balancer/internal/adapters/api/rest/erros"
	"balancer/internal/domain/models"
	"balancer/internal/domain/usecases"
	"balancer/internal/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateClientRequest struct {
	ID         string `json:"clientId"`
	Capacity   int    `json:"capacity"`
	RatePerSec int    `json:"ratePerSec"`
}

type ClientResponse struct {
	ID         string `json:"clientId"`
	Capacity   int    `json:"capacity"`
	RatePerSec int    `json:"ratePerSec"`
}

type ClientController struct {
	usecase *usecases.ClientUseCase
}

func NewClientController(uc *usecases.ClientUseCase) *ClientController {
	return &ClientController{usecase: uc}
}

func (c *ClientController) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", c.create)
	r.Get("/", c.list)
	r.Get("/{id}", c.get)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)

	return r
}

func (c *ClientController) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "creating client")

	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "invalid request body", "err", err)
		erros.JSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client := &models.Client{
		ID:         req.ID,
		Capacity:   req.Capacity,
		RatePerSec: req.RatePerSec,
	}

	created, err := c.usecase.Create(ctx, client)
	if err != nil {
		logger.Error(ctx, "failed to create client", "err", err)
		erros.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info(ctx, "client created", "id", created.ID)
	_ = json.NewEncoder(w).Encode(toClientResponse(created))
}

func (c *ClientController) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	logger.Info(ctx, "retrieving client", "id", id)

	client, err := c.usecase.Get(ctx, id)
	if err != nil {
		logger.Error(ctx, "client not found", "id", id, "err", err)
		erros.JSON(w, http.StatusNotFound, "Client not found")
		return
	}

	logger.Info(ctx, "client retrieved", "id", id)
	_ = json.NewEncoder(w).Encode(toClientResponse(client))
}

func (c *ClientController) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	logger.Info(ctx, "updating client", "id", id)

	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "invalid update body", "err", err)
		erros.JSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	client := &models.Client{
		ID:         id,
		Capacity:   req.Capacity,
		RatePerSec: req.RatePerSec,
	}

	if err := c.usecase.Update(ctx, client); err != nil {
		logger.Error(ctx, "failed to update client", "id", id, "err", err)
		erros.JSON(w, http.StatusNotFound, err.Error())
		return
	}

	logger.Info(ctx, "client updated", "id", id)
	_ = json.NewEncoder(w).Encode(toClientResponse(client))
}

func (c *ClientController) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	logger.Info(ctx, "deleting client", "id", id)

	if err := c.usecase.Delete(ctx, id); err != nil {
		logger.Error(ctx, "failed to delete client", "id", id, "err", err)
		erros.JSON(w, http.StatusNotFound, err.Error())
		return
	}

	logger.Info(ctx, "client deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (c *ClientController) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info(ctx, "listing all clients")

	clients, err := c.usecase.List(ctx)
	if err != nil {
		logger.Error(ctx, "failed to list clients", "err", err)
		erros.JSON(w, http.StatusInternalServerError, "Failed to fetch client list")
		return
	}

	var result []ClientResponse
	for _, client := range clients {
		result = append(result, toClientResponse(client))
	}

	logger.Info(ctx, "clients listed", "count", len(result))
	_ = json.NewEncoder(w).Encode(result)
}

func toClientResponse(c *models.Client) ClientResponse {
	return ClientResponse{
		ID:         c.ID,
		Capacity:   c.Capacity,
		RatePerSec: c.RatePerSec,
	}
}
