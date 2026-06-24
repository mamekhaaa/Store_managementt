package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"project-budget-service/internal/domain"
)

type Handler struct {
	service ProjectUseCase
}

func NewHandler(service ProjectUseCase) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/v1/", h.handle)
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticate(r)
	if err != nil {
		writeError(w, err)
		return
	}
	parts := splitPath(r.URL.Path)
	if len(parts) == 0 {
		writeError(w, domain.NewNotFoundError("endpoint not found"))
		return
	}
	if len(parts) == 1 && parts[0] == "resources" && r.Method == http.MethodGet {
		resources, err := h.service.ListResources(r.Context())
		respond(w, http.StatusOK, resources, err)
		return
	}
	if parts[0] != "projects" {
		writeError(w, domain.NewNotFoundError("endpoint not found"))
		return
	}
	h.handleProjects(w, r, userID, parts)
}

func (h *Handler) handleProjects(w http.ResponseWriter, r *http.Request, userID domain.UserID, parts []string) {
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			projects, err := h.service.ListProjects(r.Context(), userID)
			respond(w, http.StatusOK, projects, err)
		case http.MethodPost:
			var req createProjectRequest
			if err := decodeJSON(r, &req); err != nil {
				writeError(w, err)
				return
			}
			project, err := h.service.CreateProject(r.Context(), userID, domain.CreateProject{Name: req.Name, Description: req.Description})
			respond(w, http.StatusCreated, project, err)
		default:
			writeError(w, domain.NewNotFoundError("endpoint not found"))
		}
		return
	}

	projectID, err := domain.ParseProjectID(parts[1])
	if err != nil {
		writeError(w, err)
		return
	}

	if len(parts) == 2 {
		h.handleProject(w, r, userID, projectID)
		return
	}
	switch parts[2] {
	case "stages":
		h.handleStages(w, r, userID, projectID, parts)
	case "resources":
		h.handleProjectResources(w, r, userID, projectID, parts)
	case "access":
		h.handleAccess(w, r, userID, projectID, parts)
	default:
		writeError(w, domain.NewNotFoundError("endpoint not found"))
	}
}

func (h *Handler) handleProject(w http.ResponseWriter, r *http.Request, userID domain.UserID, projectID domain.ProjectID) {
	switch r.Method {
	case http.MethodGet:
		project, err := h.service.GetProject(r.Context(), userID, projectID)
		respond(w, http.StatusOK, project, err)
	case http.MethodPut:
		var req updateProjectRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		project, err := h.service.UpdateProject(r.Context(), userID, projectID, req.toDomain())
		respond(w, http.StatusOK, project, err)
	case http.MethodDelete:
		err := h.service.DeleteProject(r.Context(), userID, projectID)
		respondNoContent(w, err)
	default:
		writeError(w, domain.NewNotFoundError("endpoint not found"))
	}
}

func (h *Handler) handleStages(w http.ResponseWriter, r *http.Request, userID domain.UserID, projectID domain.ProjectID, parts []string) {
	if len(parts) == 3 && r.Method == http.MethodPost {
		var req stageRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		if req.StartDate.Time().IsZero() || req.EndDate.Time().IsZero() {
			writeError(w, domain.NewValidationError("start_date and end_date are required"))
			return
		}
		stage, err := h.service.CreateStage(r.Context(), userID, domain.CreateStage{ProjectID: projectID, Name: req.Name, StartDate: req.StartDate, EndDate: req.EndDate})
		respond(w, http.StatusCreated, stage, err)
		return
	}
	if len(parts) != 4 {
		writeError(w, domain.NewNotFoundError("endpoint not found"))
		return
	}
	stageID, err := domain.ParseStageID(parts[3])
	if err != nil {
		writeError(w, err)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var req stageRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		if req.StartDate.Time().IsZero() || req.EndDate.Time().IsZero() {
			writeError(w, domain.NewValidationError("start_date and end_date are required"))
			return
		}
		stage, err := h.service.UpdateStage(r.Context(), userID, domain.UpdateStage{ProjectID: projectID, StageID: stageID, StartDate: req.StartDate, EndDate: req.EndDate})
		respond(w, http.StatusOK, stage, err)
	case http.MethodDelete:
		respondNoContent(w, h.service.DeleteStage(r.Context(), userID, projectID, stageID))
	default:
		writeError(w, domain.NewNotFoundError("endpoint not found"))
	}
}

func (h *Handler) handleProjectResources(w http.ResponseWriter, r *http.Request, userID domain.UserID, projectID domain.ProjectID, parts []string) {
	if len(parts) == 3 {
		switch r.Method {
		case http.MethodGet:
			resources, err := h.service.ListProjectResources(r.Context(), userID, projectID)
			respond(w, http.StatusOK, resources, err)
		case http.MethodPost:
			var req projectResourceRequest
			if err := decodeJSON(r, &req); err != nil {
				writeError(w, err)
				return
			}
			resourceID, err := domain.ParseResourceID(req.ResourceID)
			if err != nil {
				writeError(w, err)
				return
			}
			res, err := h.service.AddProjectResource(r.Context(), userID, domain.AddProjectResource{ProjectID: projectID, ResourceID: resourceID, WorkloadPercentage: req.WorkloadPercentage})
			respond(w, http.StatusCreated, res, err)
		default:
			writeError(w, domain.NewNotFoundError("endpoint not found"))
		}
		return
	}
	if len(parts) != 4 {
		writeError(w, domain.NewNotFoundError("endpoint not found"))
		return
	}
	resourceID, err := domain.ParseResourceID(parts[3])
	if err != nil {
		writeError(w, err)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var req workloadRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		res, err := h.service.UpdateProjectResourceWorkload(r.Context(), userID, projectID, resourceID, req.WorkloadPercentage)
		respond(w, http.StatusOK, res, err)
	case http.MethodDelete:
		respondNoContent(w, h.service.DeleteProjectResource(r.Context(), userID, projectID, resourceID))
	default:
		writeError(w, domain.NewNotFoundError("endpoint not found"))
	}
}

func (h *Handler) handleAccess(w http.ResponseWriter, r *http.Request, userID domain.UserID, projectID domain.ProjectID, parts []string) {
	if len(parts) != 3 {
		writeError(w, domain.NewNotFoundError("endpoint not found"))
		return
	}
	switch r.Method {
	case http.MethodPost:
		var req accessRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		target, err := domain.ParseUserID(req.UserID)
		if err != nil {
			writeError(w, err)
			return
		}
		err = h.service.GrantAccess(r.Context(), userID, projectID, domain.ProjectAccess{UserID: target, Role: domain.AccessRole(req.Role)})
		respond(w, http.StatusCreated, map[string]string{"status": "granted"}, err)
	case http.MethodDelete:
		var req revokeAccessRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, err)
			return
		}
		target, err := domain.ParseUserID(req.UserID)
		if err != nil {
			writeError(w, err)
			return
		}
		respondNoContent(w, h.service.RevokeAccess(r.Context(), userID, projectID, target))
	default:
		writeError(w, domain.NewNotFoundError("endpoint not found"))
	}
}

func authenticate(r *http.Request) (domain.UserID, error) {
	header := r.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" || token == header {
		return "", domain.NewUnauthorizedError("bearer token is required")
	}
	return domain.ParseUserID(token)
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return domain.NewValidationError("invalid json body")
	}
	return nil
}

func splitPath(path string) []string {
	clean := strings.TrimPrefix(path, "/api/v1/")
	clean = strings.Trim(clean, "/")
	if clean == "" {
		return nil
	}
	return strings.Split(clean, "/")
}

func respond(w http.ResponseWriter, status int, data any, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, status, data)
}

func respondNoContent(w http.ResponseWriter, err error) {
	if err != nil {
		writeError(w, err)
		return
	}
	writeNoContent(w)
}
