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
		Content: "你是一位专精于技术分析和市场趋势预测的金融专家，利用您的專業知識和可用的市場分析工具來掃描並根据用户提供的当前市场数据，" +
			"確保所有資訊都是最新的，並與當前的市場狀況相關，並依造當前的市場情緒判斷出要做多還是做空，并解释其背后的原理。",
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
		//Content: "用繁體中文回答，根據全部資料預測出 未來1hr內,未來3-6hr, 未來24-72hr, 未來72hr的價格區間，並給出區間內 做多或做空 的入場價格與目標價格和止損價格，" +
		//	"並解释其背后的原理，如果需要其他資料再提供需要的資料",
		Content: "用繁體中文回答，根據我提供的資料及你在網路上蒐集的資料預測出短中長期的走勢，並給出一個完整的交易策略以及入場跟止損和目標位置",
	})

	res, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4oMini,
		Messages:    request,
		Temperature: 0.2,
		TopP:        0.5,
	})
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion err: %w", err)
	}

	return res.Choices[0].Message.Content, nil
}
