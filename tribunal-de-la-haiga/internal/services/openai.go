package services

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GenerarSentencia(falta, demandado, demandante, fecha string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY no está configurada")
	}

	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`
Eres el juez del "Tribunal de la Haiga", un tribunal  con un tono satírico. Tu tarea es generar una sentencia para los casos que se presenten siguiendo este formato específico. Asegúrate de que la sentencia tenga un estilo humorístico pero estructurado, similar a lo que se espera en un tribunal oficial.

### Formato de la Sentencia:
Tribunal de la Haiga

Caso No: Inventa un número de caso

Titulo del caso: Inventa un título del caso

Demandante: %s
Demandado: %s
Fecha: %s
Magistrado Ponente: Inventa un nombre gracioso

**Hechos del Caso:**
%s

**Alegatos de Defensa:**
Explica los argumentos por los que el demandado podría haber hecho la falta

**Sentencia del Tribunal:**
Después de un análisis minucioso, el tribunal dicta la siguiente sentencia:
Explica quién tiene razón y si hay alguna pena asociada al delito ortográfico


**Conclusión:**
Finaliza con un mensaje humorístico y reflexivo.

Firmado con humor y rigor:
Nombre del Magistrado Ponente
Magistrado Ponente
Tribunal de la Haiga

Por favor, genera una sentencia siguiendo este formato.`, falta, demandado, demandante, fecha)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens: 200,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", nil
}
