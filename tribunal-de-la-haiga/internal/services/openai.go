package services

import (
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

	// Template del prompt con estructura HTML
	promptTemplate := `
<h1>Sentencia del Tribunal de la Haiga</h1>
<p><strong>Tribunal de la Haiga Caso No:</strong> Inventa un número de caso</p>
<p><strong>Título del caso:</strong> Inventa un título del caso</p>
<p><strong>Demandante:</strong> {{.Demandante}}</p>
<p><strong>Demandado:</strong> {{.Demandado}}</p>
<p><strong>Fecha:</strong> {{.Fecha}}</p>
<p><strong>Magistrado Ponente:</strong> Inventa un nombre gracioso</p>
<h2>Hechos del Caso:</h2>
<p>{{.Falta}}</p>
<h2>Alegatos de Defensa:</h2>
<p>Explica los argumentos por los que el demandado podría haber hecho la falta</p>
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

	// Variables para el template
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
	var renderedPrompt string
	err = tmpl.Execute(&renderedPrompt, data)
	if err != nil {
		return "", fmt.Errorf("error ejecutando el template: %v", err)
	}

	// Llamada a la API
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: renderedPrompt},
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
