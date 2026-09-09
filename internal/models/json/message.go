package json

import "github.com/lbenedar/fresta/internal/models/db"

type Message struct {
	Content   string   `json:"content"`
	Style     int      `json:"style"`
	Author    string   `json:"author"`
	Id        string   `json:"_id"`
	Type      string   `json:"type"`
	Timestamp int64    `json:"timestamp"`
	Flavor    string   `json:"flavor"`
	Speaker   Speaker  `json:"speaker"`
	Whisper   []string `json:"whisper"`
	Blind     bool     `json:"blind"`
	Rolls     []string `json:"rolls"`
	Sound     string   `json:"sound"`
	Emote     bool     `json:"emote"`
	Stats     Stats    `json:"_stats"`
	// System    any      `json:"system"`
	// Flags     any `json:"flags"`
}

func (m *Message) ToDB(dest **db.Message) bool {
	if dest == nil {
		return false
	}

	message := &db.Message{
		Content:   m.Content,
		Style:     m.Style,
		Author:    m.Author,
		ID:        m.Id,
		Type:      m.Type,
		Timestamp: m.Timestamp,
		Flavor:    m.Flavor,
		Blind:     m.Blind,
		Sound:     m.Sound,
		Emote:     m.Emote,
	}

	m.Speaker.ToDB(&message.Speaker)
	m.Stats.ToDB(&message.Stats)

	message.Whisper = make([]string, len(m.Whisper))
	copy(message.Whisper, m.Whisper)
	message.Rolls = make([]string, len(m.Rolls))
	copy(message.Rolls, m.Rolls)

	*dest = message

	return true
}

type Speaker struct {
	Scene string `json:"scene,omitempty"`
	Actor string `json:"actor,omitempty"`
	Token string `json:"token,omitempty"`
	Alias string `json:"alias,omitempty"`
}

func (s *Speaker) ToDB(dest *db.Speaker) bool {
	if dest == nil {
		return false
	}

	dest.Scene = s.Scene
	dest.Actor = s.Actor
	dest.Token = s.Token
	dest.Alias = s.Alias

	return true
}
