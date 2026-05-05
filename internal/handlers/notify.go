package handlers

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"

	"github.com/shridarpatil/whatomate/internal/websocket"
)

type NotifyRequest struct {
	OrgID   string         `json:"org_id"`
	UserID  string         `json:"user_id"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (a *App) Notify(r *fastglue.Request) error {
	var req NotifyRequest
	// 1. Decode request con error visible
	if err := json.Unmarshal(r.RequestCtx.PostBody(), &req); err != nil {
		return r.SendErrorEnvelope(
			fasthttp.StatusBadRequest,
			"invalid JSON body",
			nil,
			"",
		)
	}

	// 2. Validaciones obligatorias
	if req.OrgID == "" {
		return r.SendErrorEnvelope(
			fasthttp.StatusBadRequest,
			"org_id is required",
			nil,
			"",
		)
	}

	if req.Type == "" {
		return r.SendErrorEnvelope(
			fasthttp.StatusBadRequest,
			"type is required",
			nil,
			"",
		)
	}

	// 3. Parse org_id
	orgID, err := uuid.Parse(req.OrgID)
	if err != nil {
		return r.SendErrorEnvelope(
			fasthttp.StatusBadRequest,
			"invalid org_id",
			nil,
			"",
		)
	}

	// 4. Construir mensaje websocket
	msg := websocket.WSMessage{
		Type:    req.Type,
		Payload: req.Payload,
	}

	// 5. Enviar a usuario específico o a toda la org
	if req.UserID != "" {
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return r.SendErrorEnvelope(
				fasthttp.StatusBadRequest,
				"invalid user_id",
				nil,
				"",
			)
		}

		a.WSHub.BroadcastToUser(orgID, userID, msg)
	} else {
		a.WSHub.BroadcastToOrg(orgID, msg)
	}

	// 6. Respuesta OK
	return r.SendEnvelope(map[string]any{
		"ok": true,
	})
}