package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	openai "github.com/sashabaranov/go-openai"
)

func GenerarSentencia(falta, demandado, demandante, fecha string) (string, error) {
	// Verificar que la API key esté configurada
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY no está configurada")
	}

	client := openai.NewClient(apiKey)

	// Template del prompt
	promptTemplate := `
Eres el magistrado del tribunal de la Haiga, un tribunal satírico. 
Encargado de detectar si en la frase - {{.Falta}} - hay una falta de ortografía, tal y como argumenta {{.Demandante}}. O si no la hay.
Si hay una falta hay que dictar sentencia, explicar la falta y proponer una pena.
Centra toda la sentencia en la ortografía y gramática de la frase:  {{.Falta}}

<p><strong>Tribunal de la Haiga Caso No:</strong> Inventa un número de caso</p>
<p><strong>Título del caso:</strong> Inventa un título del caso</p>
<p><strong>Demandante:</strong> {{.Demandante}}</p>
<p><strong>Demandado:</strong> {{.Demandado}}</p>
<p><strong>Fecha:</strong> {{.Fecha}}</p>
<p><strong>Magistrado Ponente:</strong> Inventa un nombre gracioso</p>
<h2>Hechos del Caso:</h2>
<p>Expone unos hechos básandote en la siguiente falta: {{.Falta}}</p>
<h2>Alegatos de Defensa:</h2>
<p>Explica los argumentos por los que el demandado podría haber hecho la falta de ortografía</p>
<h2>Sentencia del Tribunal:</h2>
<p>Después de un análisis minucioso, el tribunal dicta la siguiente sentencia:
Explica quién tiene razón y si hay alguna pena asociada al delito ortográfico</p>
<h2>Conclusión:</h2>
<p>Finaliza con un mensaje humorístico y reflexivo</p>
<p>Firmado con humor y rigor: Nombre del Magistrado Ponente</p>
`

	// Compilar el template
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("error creando el template: %v", err)
	}

	// Datos para el template
	data := struct {
		Falta      string
		Demandado  string
		Demandante string
		Fecha      string
	}{
		Falta:      falta,
		Demandado:  demandado,
		Demandante: demandante,
		Fecha:      fecha,
	}

	// Renderizar el prompt
	var filledPrompt bytes.Buffer
	err = tmpl.Execute(&filledPrompt, data)
	if err != nil {
		return "", fmt.Errorf("error ejecutando el template: %v", err)
	}

	// Llamada a la API de OpenAI
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-4o", // Modelo actualizado a GPT-4o
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: filledPrompt.String()},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return "", err
	}

	// Verificar la respuesta
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", nil
}
