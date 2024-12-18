package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"activator/internal"
	"activator/internal/app/models"
	"activator/internal/config"
	"activator/internal/mailer"
	m "activator/internal/rest/models"
)

type UserService interface {
	Create(params m.CreateParams) (models.User, error)
	Delete(id int32) error
	Find(id int32) (models.User, error)
	Activate(tokenPlaintext string) (int32, error)
}

type TokenService interface {
	Create(user models.User) (models.Token, error)
	DeleteAll(id int32) error
}

type UserHandler struct {
	cfg      config.Config
	logger   *slog.Logger
	svcUser  UserService
	svcToken TokenService
	mailer   mailer.Mailer
}

func NewUserHandler(
	cfg config.Config,
	logger *slog.Logger,
	svcUser UserService,
	svcToken TokenService,
	mailer mailer.Mailer,
) *UserHandler {
	return &UserHandler{
		cfg:      cfg,
		logger:   logger,
		svcUser:  svcUser,
		svcToken: svcToken,
		mailer:   mailer,
	}
}

func (h *UserHandler) Register(r *mux.Router) {
	r.HandleFunc("/users", h.create).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.delete).Methods(http.MethodDelete)
	r.HandleFunc("/users/{id}", h.find).Methods(http.MethodGet)
	r.HandleFunc("/activate", h.activate).Methods(http.MethodGet)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var createParams m.CreateParams
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid json user params")

		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	defer r.Body.Close()

	if err := createParams.Validate(); err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid validation user params")

		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	createdUser, err := h.svcUser.Create(createParams)
	if err != nil {
		msg := fmt.Errorf("create failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Info("POST request success, record created", "id", createdUser.ID)

	token, err := h.svcToken.Create(createdUser)
	if err != nil {
		msg := fmt.Errorf("token create failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Debug("token created for user", "token", token.Plaintext)

	if err := h.mailer.Send(createdUser, token.Plaintext); err != nil {
		msg := fmt.Errorf("failed sending activation: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Info("activation letter sent", "user id", createdUser.ID)

	renderResponse(w, createdUser, http.StatusCreated)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid id")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	if err := h.svcUser.Delete(int32(id)); err != nil {
		msg := fmt.Errorf("delete failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	h.logger.Info("DELETE request success, record deleted", "id", id)
	renderResponse(w, map[string]string{"result": "success"}, http.StatusOK)
}

func (h *UserHandler) activate(w http.ResponseWriter, r *http.Request) {
	tokenPlaintext := r.URL.Query().Get("token")
	id, err := h.svcUser.Activate(tokenPlaintext)
	if err != nil {
		msg := fmt.Errorf("activation failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	if err := h.svcToken.DeleteAll(id); err != nil {
		msg := fmt.Errorf("activation failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	h.logger.Info("Activation request success")
	renderResponse(w, map[string]string{"result": "success"}, http.StatusOK)
}

func (h *UserHandler) find(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		msg := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "invalid id")
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	user, err := h.svcUser.Find(int32(id))
	if err != nil {
		msg := fmt.Errorf("find failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}
	h.logger.Info("Find request success, record found", "id", id)

	renderResponse(w, user, http.StatusOK)
}
