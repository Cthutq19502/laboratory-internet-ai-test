package gigachat

type Request struct {
	Messages []string
}

// ------------------------------------------------------

type AccessTokenResponse struct {
	Token     string `json:"access_token"`
	ExpiresAt uint64 `json:"expires_at"`
}

//-------------------------------------------------------

type RequestMessage struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	N                 int       `json:"n"`
	Temperature       float64   `json:"temperature"`
	Stream            bool      `json:"stream"`
	MaxTokens         int       `json:"max_tokens"`
	RepetitionPenalty int       `json:"repetition_penalty"`
	UpdateInterval    int       `json:"update_interval"`
	*FunctionGigachat
}

type FunctionGigachat struct {
	FunctionCall    functionCall `json:"function_call"`
	Functions       []Functions  `json:"functions"`
	FewShotExamples interface{}  `json:"few_shot_examples"`
}
type functionCall struct {
	Name string `json:"name"`
}

type Functions struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  parameters `json:"parameters"`
	Required    []string   `json:"required"`
}

type parameters struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

//-------------------------------------------------------

type Response struct {
	Choices []Choices `json:"choices"`
}

type Choices struct {
	MessageResponse MessagesResponse `json:"message"`
	FinishReason    string           `json:"finish_reason"`
}

type MessagesResponse struct {
	Role         string       `json:"role"`
	Content      string       `json:"content"`
	FunctionCall FunctionCall `json:"function_call"`
}

type FunctionCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

//------------------------------------------------------------

type RequestEmbedding struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ResponseEmbedding struct {
	Object string          `json:"object"`
	Model  string          `json:"model"`
	Data   []dataEmbedding `json:"data"`
}

type dataEmbedding struct {
	Embedding []float32 `json:"embedding"`
}

//------------------------------------------------------------
