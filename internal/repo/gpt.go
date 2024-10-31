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
	request = append(request,
		openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleSystem,
			//Content: "你是一位專精於技術分析和市場趨勢預測的金融專家，擅長根據多維度數據和市場情緒進行判斷。" +
			//	"請使用您的專業知識及技術指標（例如 RSI、MACD、ATR 等）以及最新市場資訊，綜合分析用戶提供的數據。" +
			//	"您將計算 RSI 和 MACD，推算市場情緒狀態，以判斷是否進行反向操作來控制風險。" +
			//	"同時，根據 ATR 動態設定止損位，並提供精確的入場和目標價格。" +
			//	"您的目標是通過這些數據為用戶提供一個針對短期（1-6小時）、中期（24-72小時）和長期（72小時以上）的全面動態交易策略。",
			Content: "你是一位專精於技術分析和市場趨勢預測的金融專家，擅長根據多維度數據、聰明錢策略和市場情緒進行判斷。" +
				"請利用您的專業知識、技術指標（如 RSI、MACD、ATR 等）、聰明錢策略（大額資金流、未平倉合約、融資費率、頭寸比率）及其他市場分析工具，綜合分析用戶提供的數據和最新市場資訊。" +
				"您的目標是預測短期（1-6小時）、中期（24-72小時）和長期（72小時以上）的價格走勢。" +
				"根據不同策略的分析結果，提供針對性建議，包括支撐位、阻力位、入場和目標價格，並根據 ATR 計算止損價。" +
				"當市場情緒過高或過低時，通過反向操作策略來控制風險。",
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
		//Content: "請用繁體中文回答，根據我提供的所有資料及您蒐集的最新數據，預測 短期、中期和長期的價格走勢。請包括以下內容：" +
		//	"\n1. **價格區間預測**：根據技術指標和數據，預測每個時段的價格區間，並提供支撐位和阻力位。" +
		//	"\n2. **技術指標解讀**：解釋 RSI 和 MACD 的具體變化（如 RSI 過買過賣，MACD 金叉或死叉等）對預測的影響。" +
		//	"\n3. **交易策略**：根據市場情緒設置動態交易策略，包含每個時段的進場、止損價格（使用 ATR 計算）和目標位置。" +
		//	"\n4. **市場情緒反向操作**：當市場情緒過高或過低時，說明反向操作的建議及其理由，以減少風險。" +
		//	"\n5. **參考資料**：列出您在分析中使用的數據和指標，幫助用戶理解決策依據。" +
		//	"\n若您需要其他數據，請列出需求，以便我提供更準確的預測。",
		Content: "請用繁體中文回答，根據我提供的所有資料及您蒐集的最新數據，預測 短期、中期和長期的價格走勢。請包括以下內容：" +
			"\n1. **價格區間預測**：使用技術指標和聰明錢策略（大額資金流向、未平倉合約、融資費率、多空頭寸比率等），針對不同時段預測價格區間，並提供支撐位和阻力位。" +
			"\n2. **技術指標解讀**：解釋 RSI、MACD 等技術指標的變化及其對價格走勢的影響（如 RSI 超買或超賣信號，MACD 的金叉或死叉）。" +
			"\n3. **交易策略**：結合聰明錢策略和技術指標的信號，根據市場情緒，為短期、中期和長期提供具體的進場點、止損（根據 ATR 計算）和目標位置。" +
			"\n4. **市場情緒反向操作**：當多空比率和資金費率顯示市場情緒過高或過低時，給出反向操作建議並解釋原因。" +
			"\n5. **參考資料**：列出分析中所用的所有數據來源（如技術指標、融資費率、未平倉合約、頭寸比率等），幫助用戶理解決策依據。" +
			"\n若您需要其他數據，請列出需求，以便我提供更準確的預測。",
	})
	res, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4oMini,
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
