package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/bwmarrin/discordgo"
)

const templatesPath = "cmd/api/discord_bot/response/templates"

type RespondData struct {
	Content    string
	Components []discordgo.MessageComponent
	Embeds     []*discordgo.MessageEmbed
}

func GetParserResponse(cmd string, data any) (*discordgo.InteractionResponseData, error) {
	jsonMsg, err := GetJsonMessage(cmd, data)
	if err != nil {
		return nil, err
	}

	resp, err := ParseMessage(jsonMsg)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ParseMessage(jsonMsg []byte) (*discordgo.InteractionResponseData, error) {

	jsData := struct {
		Content    string                    `json:"content"`
		Components []json.RawMessage         `json:"components"`
		Embeds     []*discordgo.MessageEmbed `json:"embeds"`
	}{}

	err := json.Unmarshal(jsonMsg, &jsData)
	if err != nil {
		return nil, err
	}

	components := make([]discordgo.MessageComponent, len(jsData.Components))
	for i, raw := range jsData.Components {
		cmp, err := discordgo.MessageComponentFromJSON(raw)
		if err != nil {
			return nil, err
		}
		components[i] = cmp
	}

	response := &discordgo.InteractionResponseData{
		Content:    jsData.Content,
		Components: components,
		Embeds:     jsData.Embeds,
	}

	return response, nil
}

func GetJsonMessage(filename string, data any) ([]byte, error) {
	var jsData []byte
	var err error

	if data == nil {
		jsData, err = os.ReadFile(fmt.Sprintf("%s/%s.json.tmpl", templatesPath, filename))
		if err != nil {
			return nil, err
		}
	} else {
		tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s.json.tmpl", templatesPath, filename))
		if err != nil {
			return nil, err
		}

		buffer := &bytes.Buffer{}
		err = tmpl.Execute(buffer, data)
		if err != nil {
			return nil, err
		}

		jsData = buffer.Bytes()
	}

	return jsData, nil
}
