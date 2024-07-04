package main

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"io"
	"log/slog"
	"time"
)

var VertexIns VertexClient

type VertexClient struct {
	ProjectID string
	Location  string
	KeyFile   string
	client    *genai.Client
}

func InitVertexInstance(projectID, location, keyFile string) error {
	slog.Info("NewVertexInstance", slog.Any("projectID", projectID), slog.Any("location", location), slog.Any("keyFile", keyFile))
	var v = &VertexClient{
		ProjectID: projectID,
		Location:  location,
		KeyFile:   keyFile,
	}
	opt := option.WithCredentialsFile(keyFile)
	c, err := genai.NewClient(context.Background(), v.ProjectID, v.Location, opt)
	if err != nil {
		return err
	}
	v.client = c
	VertexIns = *v
	return nil
}

type OpenAIChatCompletion struct {
	Model       string              `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages"`
	Stream      bool                `json:"stream"`
	Temperature float32             `json:"temperature"` // 0-2
}

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"` // string or OpenAIChatContent
}

//type OpenAIChatMessage struct {
//	Role    string      `json:"role"`
//	Content interface{} `json:"content"` // string or OpenAIChatContent
//}

//type OpenAIChatContent struct {
//	Type     string  `json:"type"`
//	Text     *string `json:"text"`
//	ImageUrl *Image  `json:"image_url"`
//}
//
//type Image struct {
//	Url string `json:"url"`
//}

type OpenAIChatChoice struct {
	Index        int                `json:"index"`
	Message      *OpenAIChatMessage `json:"message,omitempty"`
	Delta        *OpenAIDelta       `json:"delta,omitempty"`
	LogProbs     interface{}        `json:"logprobs"`
	FinishReason *string            `json:"finish_reason"`
}

type OpenAIDelta struct {
	Content string `json:"content"`
}

type OpenAIUsage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

type OpenAIChatResponse struct {
	Id                string             `json:"id"`
	Object            string             `json:"object"`
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	SystemFingerprint *string            `json:"system_fingerprint"`
	Choices           []OpenAIChatChoice `json:"choices"`
	Usage             *OpenAIUsage       `json:"usage,omitempty"`
}

func FormatTemperature(o OpenAIChatCompletion) float32 {
	if o.Temperature < 0.0 || o.Temperature > 2 {
		return float32(0.5)
	}
	return o.Temperature
}

func OpenAI2VerTexAI(c *gin.Context, o OpenAIChatCompletion, model *genai.GenerativeModel) error {
	var contents []*genai.Content

	NewMessage := o.Messages[len(o.Messages)-1]
	NewMessageContent := NewMessage.Content
	HistoryMessage := o.Messages[:len(o.Messages)-1]

	model.SetTemperature(FormatTemperature(o))

	session := model.StartChat()
	session.History = nil

	for _, m := range HistoryMessage {
		if m.Content == "" {
			continue
		}
		var content genai.Content
		if m.Role == "system" {
			content.Parts = append(content.Parts, genai.Text(m.Content))
			model.SystemInstruction = &content
			slog.Debug("SystemInstruction", slog.Any("SystemInstruction", model.SystemInstruction))
		} else {
			if m.Role == "user" {
				content.Role = m.Role
			} else {
				content.Role = "model"
			}
			content.Parts = append(content.Parts, genai.Text(m.Content))
			contents = append(contents, &content)
		}

	}
	session.History = contents
	var response OpenAIChatResponse
	response.Id = "chatcmpl-123"
	response.Created = time.Now().Unix()
	response.Model = model.Name()
	response.SystemFingerprint = nil
	if o.Stream {
		i := session.SendMessageStream(c, genai.Text(NewMessageContent))
		return VerTexStreamOutPut(c, i, &response)
	} else {
		msg, err := session.SendMessage(c, genai.Text(NewMessageContent))
		if err != nil {
			return err
		}
		return VerTexOutPut(c, msg, &response)
	}

}

func VerTexOutPut(c *gin.Context, msg *genai.GenerateContentResponse, response *OpenAIChatResponse) error {
	response.Created = time.Now().Unix()
	response.Object = "chat.completion"
	response.Choices = []OpenAIChatChoice{
		{
			Index: 0,
			Message: &OpenAIChatMessage{
				Role:    "assistant",
				Content: string(msg.Candidates[0].Content.Parts[0].(genai.Text)),
			},
			LogProbs:     nil,
			FinishReason: nil,
		},
	}
	if msg.Candidates[0].FinishReason == 1 {
		response.Choices[0].FinishReason = new(string)
		*response.Choices[0].FinishReason = "stop"
	}
	if response.Usage != nil {
		response.Usage = &OpenAIUsage{
			PromptTokens:     msg.UsageMetadata.PromptTokenCount,
			CompletionTokens: msg.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      msg.UsageMetadata.TotalTokenCount,
		}
	}
	c.JSON(200, response)
	return nil
}

func VerTexStreamOutPut(c *gin.Context, i *genai.GenerateContentResponseIterator, response *OpenAIChatResponse) error {
	messageChan := make(chan []byte)
	response.Object = "chat.completion.chunk"
	go func() {
		for {
			resp, err := i.Next()
			if err != nil {
				slog.Error("Stream error", slog.Any("err", err))
				break // Exit the loop if there's an error
			}

			respText := string(resp.Candidates[0].Content.Parts[0].(genai.Text))

			response.Created = time.Now().Unix()
			response.Choices = []OpenAIChatChoice{{
				Index: 0,
				Delta: &OpenAIDelta{
					Content: respText,
				},
				FinishReason: nil,
			}}

			clog, _ := json.Marshal(resp.Candidates[0])
			slog.Info("resp", slog.Any("resp", clog))

			if resp.Candidates[0].FinishReason == 1 {
				response.Choices[0].FinishReason = new(string)
				*response.Choices[0].FinishReason = "stop"
				if resp.UsageMetadata != nil {
					response.Usage = &OpenAIUsage{
						PromptTokens:     resp.UsageMetadata.PromptTokenCount,
						CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
						TotalTokens:      resp.UsageMetadata.TotalTokenCount,
					}
				}
			}

			jsonData, _ := json.Marshal(response) // Simplified error handling
			messageChan <- jsonData
			if resp.Candidates[0].FinishReason == 1 {
				slog.Debug("Stream finish", slog.Any("FinishReason", resp.Candidates[0].FinishReason))
				break // Exit the loop if finish reason is stop
			}
		}
		// 强制发送一个Stop

		close(messageChan) // Close the channel after the loop ends
	}()

	// 从通道读取数据并写入响应
	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-messageChan; ok {
			writeSSE(w, string(msg))
			return true
		}
		//c.SSEvent("message", "[DONE]")
		writeSSE(w, "[DONE]")
		return false
	})
	return nil
}

func writeSSE(w io.Writer, data string) {
	fmt.Fprintf(w, "data: %s\n\n", data)
}
