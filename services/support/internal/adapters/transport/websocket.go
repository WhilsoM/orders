package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Hub struct {
	Chats    map[int]*Chat
	Supports map[int]*Client
}

type Handler struct {
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	idString := r.URL.Query().Get("id")

	if role != "user" && role != "support" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}

	slog.Info("websocket connected")
	defer conn.CloseNow()
	defer slog.Info("websocket disconnected")

	client := &Client{
		ID:   id,
		Role: role,
		Conn: conn,
	}

	slog.Info(
		"websocket connected",
		"id", client.ID,
		"role", client.Role,
	)

	for {
		var msg Message
		err := wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			return
		}

		slog.Info("message received", "message", msg)

		err = wsjson.Write(r.Context(), conn, msg)
		if err != nil {
			return
		}
	}
}
