package repo

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/zaqxsw800402/binance/internal/model"
	"github.com/zaqxsw800402/binance/pkg/config"
)

type ChatGpt struct {
	client *openai.Client
}

func NewChatGpt(cfg config.OpenAI) *ChatGpt {
	client := openai.NewClient(cfg.ApiKey)
	return &ChatGpt{
		client: client,
	}
}

func (c *ChatGpt) Chat(ctx context.Context, prompts *model.AllPrompt) (string, error) {
	var request []openai.ChatCompletionMessage
	request = append(request, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		Content: "你是一位專精於技術分析和市場趨勢預測的金融專家，擅長根據多維度數據和市場情緒進行判斷。" +
			"利用您的專業知識和市場分析工具來分析用戶提供的所有數據，並結合網路上最新的相關市場資訊，預測出短期、中期和長期的價格走勢。" +
			"在市場情緒過高或過低時，適當做出反向操作，以減少風險。",
	},
	)

	for _, prompt := range prompts.Prompts {
		request = append(request, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt.Question,
		}, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt.Data,
		},
		)
	}

	request = append(request, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: "請用繁體中文回答，根據我提供的所有資料及您蒐集的最新數據，詳細預測出短期（1-6小時）、中期（24-72小時）和長期（72小時以上）價格走勢。" +
			"請提供以下內容：" +
			"\n1. **每個時段的價格區間預測**：預測價格區間，並提供每段時間的支撐位和阻力位。" +
			"\n2. **交易策略**：根據每個時段的市場情緒，決定是做多或做空。包括進場價格、止損價格和目標位置。" +
			"\n3. **市場情緒反向操作**：當市場情緒過高時，說明反向操作的建議及其理由。" +
			"\n4. **參考資料**：列出您在分析中所使用的所有數據來源，幫助用戶理解決策依據。" +
			"\n若需要額外資料，請列出需求。",
	})
	res, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		Messages:    request,
		Temperature: 0.5,
		//TopP:        0.5,
	})
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion err: %w", err)
	}

	return res.Choices[0].Message.Content, nil
}

func (c *ChatGpt) ChatV2(ctx context.Context, prompts *model.AllPrompt) (string, error) {
	var request []openai.ChatCompletionMessage

	// 優化後的系統提示
	request = append(request, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是一位專精於技術分析和市場趨勢預測的金融專家。請根據用戶提供的數據和最新的市場資訊，預測短期、中期和長期的價格走勢，並在市場情緒極端時提出反向操作建議。",
	})

	// 合併用戶的問題和數據
	var combinedData string
	for _, prompt := range prompts.Prompts {
		combinedData += prompt.Question + "\n" + prompt.Data + "\n"
	}

	// 優化後的用戶提示
	userPrompt := "請用繁體中文回答，根據以下資料以及您蒐集的最新數據，詳細預測短期（1-6小時）、中期（24-72小時）和長期（72小時以上）的價格走勢。請提供以下內容：" +
		"\n1. **每個時段的價格區間預測**：預測價格區間，並提供每段時間的支撐位和阻力位。" +
		"\n2. **交易策略**：根據每個時段的市場情緒，決定是做多或做空，並提供進場價格、止損價格和目標位置。" +
		"\n3. **市場情緒反向操作**：當市場情緒過高或過低時，提供反向操作的建議及理由。" +
		"\n4. **參考資料**：列出您在分析中使用的所有數據來源，幫助我理解決策依據。" +
		"\n若需要額外資料，請列出需求。"

	// 將合併的數據和用戶提示添加到請求中
	request = append(request, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: combinedData + "\n" + userPrompt,
	})

	res, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4o,
		Messages:    request,
		Temperature: 0.5,
		//TopP:        0.5,
	})
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion err: %w", err)
	}

	return res.Choices[0].Message.Content, nil
}
