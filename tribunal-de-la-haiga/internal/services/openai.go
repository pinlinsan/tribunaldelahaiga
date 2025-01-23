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
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY no está configurada")
	}

	client := openai.NewClient(apiKey)

	// Plantilla para el prompt
	tmplText := `
Eres el juez del "Tribunal de la Haiga", un tribunal con un tono satírico. 
Tu tarea es generar una sentencia para los casos que se presenten siguiendo este formato específico. 
Asegúrate de que la sentencia tenga un estilo humorístico pero estructurado, similar a lo que se espera en un tribunal oficial.

### Formato de la Sentencia:
Tribunal de la Haiga

Caso No: Inventa un número de caso

Titulo del caso: Inventa un título del caso basado en la falta cometida ({{.Falta}})

Demandante: {{.Demandante}}
Demandado: {{.Demandado}}
Fecha: {{.Fecha}}
Magistrado Ponente: Inventa un nombre gracioso

**Hechos del Caso:**
Explica los hechos basandte en la siguiente falta: {{.Falta}} incluye al demandante y demandado

**Alegatos de Defensa:**
Explica los argumentos por los que el demandado podría haber hecho la falta de escribier {{.Falta}}

**Sentencia del Tribunal:**
Después de un análisis minucioso, el tribunal dicta la siguiente sentencia:
Explica quién tiene razón y si hay alguna pena asociada al delito ortográfico

**Conclusión:**
Finaliza con un mensaje humorístico y reflexivo.

Firmado con humor y rigor:
Nombre del Magistrado Ponente
Magistrado Ponente
Tribunal de la Haiga


`

	// Crear la plantilla
	tmpl, err := template.New("prompt").Parse(tmplText)
	if err != nil {
		return "", err
	}

	// Datos para llenar la plantilla
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

	// Rellenar la plantilla con los datos
	var filledPrompt bytes.Buffer
	err = tmpl.Execute(&filledPrompt, data)
	if err != nil {
		return "", err
	}

	// Llamar a la API con el prompt generado
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-4o", // Cambiado a GPT-4o
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: filledPrompt.String()},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", nil
}
