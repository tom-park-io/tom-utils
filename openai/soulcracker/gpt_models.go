package openai

// Model constants for use with OpenAI's chat/completion APIs (as of April 2025)
const (
	// GPT-4.1 series (latest as of April 2025)
	ModelGPT41     = "gpt-4.1"      // Latest high-capacity GPT-4.1 model
	ModelGPT41Mini = "gpt-4.1-mini" // Optimized lightweight GPT-4.1 model
	ModelGPT41Nano = "gpt-4.1-nano" // Tiny version for edge devices or minimal inference cost

	// GPT-4-turbo variants
	ModelGPT4Turbo     = "gpt-4-turbo"        // Most capable general GPT-4 model
	ModelGPT4Turbo0125 = "gpt-4-0125-preview" // January 2025 preview
	ModelGPT4Turbo1106 = "gpt-4-1106-preview" // November 2024 preview

	// GPT-4 legacy (may be deprecated or costlier)
	ModelGPT4_32k = "gpt-4-32k" // Legacy GPT-4 with 32k context (may not be available for all users)

	// GPT-3.5-turbo variants (cost-effective alternatives)
	ModelGPT35Turbo     = "gpt-3.5-turbo"      // Default GPT-3.5 model
	ModelGPT35Turbo0125 = "gpt-3.5-turbo-0125" // Latest stable 3.5 version (Jan 2025)
	ModelGPT35Turbo1106 = "gpt-3.5-turbo-1106" // November 2024 variant

	// Embedding models (optional, for vector search etc.)
	ModelTextEmbedding3Small = "text-embedding-3-small" // Lightweight embedding
	ModelTextEmbedding3Large = "text-embedding-3-large" // High quality embedding
)
