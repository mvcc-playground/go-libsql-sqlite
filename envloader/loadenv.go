package envloader

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// LoadOptions configurações opcionais para carregar variáveis
type LoadOptions struct {
	EnvFile           string // Caminho para o arquivo .env (opcional)
	RequiredByDefault bool   // Se campos sem tag validate devem ser required (padrão: true)
}

// Load carrega e valida as variáveis de ambiente na struct
func Load(target interface{}, opts ...LoadOptions) error {
	var options LoadOptions
	if len(opts) > 0 {
		options = opts[0]
	} else {
		options = LoadOptions{
			EnvFile:           ".env",
			RequiredByDefault: true,
		}
	}

	// Carregar arquivo .env se especificado
	if options.EnvFile != "" {
		if _, err := os.Stat(options.EnvFile); err == nil {
			if err := godotenv.Load(options.EnvFile); err != nil {
				return fmt.Errorf("erro ao carregar arquivo %s: %w", options.EnvFile, err)
			}
		}
	}

	// Processar a struct usando reflection
	if err := processStruct(target, options); err != nil {
		return err
	}

	// Validar usando validator
	if err := validate.Struct(target); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// processStruct processa os campos da struct e popula com variáveis de ambiente
func processStruct(target interface{}, options LoadOptions) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target deve ser um ponteiro para struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	// Mapa para guardar valores das variáveis para interpolação
	envValues := make(map[string]string)

	// Primeira passada: carregar todas as variáveis simples
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Pular campos não exportados
		if !field.CanSet() {
			continue
		}

		// Obter nome da variável de ambiente
		envName := getEnvName(fieldType)

		// Obter valor da variável de ambiente
		envValue := os.Getenv(envName)

		// Guardar no mapa para interpolação
		if envValue != "" {
			envValues[envName] = envValue
			envValues[fieldType.Name] = envValue // Também usar nome do campo
		}
	}

	// Segunda passada: processar e definir valores com interpolação
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Pular campos não exportados
		if !field.CanSet() {
			continue
		}

		// Obter nome da variável de ambiente
		envName := getEnvName(fieldType)

		// Obter valor da variável de ambiente
		envValue := os.Getenv(envName)

		// Verificar se é required
		isRequired := isFieldRequired(fieldType, options.RequiredByDefault)

		if envValue == "" && isRequired {
			return fmt.Errorf("variável de ambiente obrigatória não encontrada: %s", envName)
		}

		// Processar interpolação se necessário
		if envValue != "" {
			processedValue, err := interpolateValue(envValue, envValues)
			if err != nil {
				return fmt.Errorf("erro ao processar interpolação para %s: %w", envName, err)
			}

			if err := setFieldValue(field, processedValue); err != nil {
				return fmt.Errorf("erro ao definir valor para %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

// getEnvName retorna o nome da variável de ambiente para o campo
func getEnvName(field reflect.StructField) string {
	envTag := field.Tag.Get("env")
	if envTag != "" {
		return envTag
	}
	return field.Name
}

// isFieldRequired verifica se o campo é obrigatório
func isFieldRequired(field reflect.StructField, defaultRequired bool) bool {
	validateTag := field.Tag.Get("validate")

	if validateTag == "" {
		return defaultRequired
	}

	// Verificar se contém "required"
	if strings.Contains(validateTag, "required") {
		return true
	}

	// Se contém "omitempty", não é required
	if strings.Contains(validateTag, "omitempty") {
		return false
	}

	return defaultRequired
}

// interpolateValue processa interpolação de variáveis em valores que estão entre []
func interpolateValue(value string, envValues map[string]string) (string, error) {
	// Verificar se está entre []
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		// Remover os colchetes
		template := value[1 : len(value)-1]

		// Regex para encontrar {VARIAVEL}
		re := regexp.MustCompile(`\{([^}]+)\}`)

		// Substituir todas as ocorrências
		result := re.ReplaceAllStringFunc(template, func(match string) string {
			// Extrair nome da variável (remover { e })
			varName := match[1 : len(match)-1]

			// Buscar valor da variável
			if val, exists := envValues[varName]; exists {
				return val
			}

			// Se não encontrar, manter o original
			return match
		})

		// Verificar se ainda há variáveis não resolvidas
		if re.MatchString(result) {
			unresolvedVars := re.FindAllStringSubmatch(result, -1)
			var missingVars []string
			for _, match := range unresolvedVars {
				missingVars = append(missingVars, match[1])
			}
			return "", fmt.Errorf("variáveis não encontradas para interpolação: %s", strings.Join(missingVars, ", "))
		}

		return result, nil
	}

	// Se não está entre [], retornar como está
	return value, nil
}
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Para int, você pode expandir isso se precisar
		return fmt.Errorf("tipo int não implementado ainda")
	case reflect.Bool:
		// Para bool, você pode expandir isso se precisar
		return fmt.Errorf("tipo bool não implementado ainda")
	default:
		return fmt.Errorf("tipo não suportado: %s", field.Kind())
	}
	return nil
}

// formatValidationError formata erros de validação de forma mais legível
func formatValidationError(err error) error {
	var errorMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		var message string

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("campo '%s' é obrigatório", err.Field())
		case "url":
			message = fmt.Sprintf("campo '%s' deve ser uma URL válida", err.Field())
		case "email":
			message = fmt.Sprintf("campo '%s' deve ser um email válido", err.Field())
		default:
			message = fmt.Sprintf("campo '%s' falhou na validação '%s'", err.Field(), err.Tag())
		}

		errorMessages = append(errorMessages, message)
	}

	return fmt.Errorf("erros de validação:\n- %s", strings.Join(errorMessages, "\n- "))
}

// MustLoad como Load, mas entra em pânico se houver erro
func MustLoad(target interface{}, opts ...LoadOptions) {
	if err := Load(target, opts...); err != nil {
		panic(err)
	}
}

// ex
