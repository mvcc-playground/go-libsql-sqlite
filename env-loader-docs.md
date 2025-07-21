# 📚 ENV Loader - Documentação

Uma biblioteca Go minimalista para carregar e validar variáveis de ambiente automaticamente com suporte a interpolação.

## 🚀 Instalação

```bash
go mod init meu-projeto
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
```

## ⚡ Uso Básico

### 1. Definir Struct

```go
type Config struct {
    DATABASE_URL string `env:"DB_URL" validate:"required,url"`
    API_KEY      string `validate:"required,min=10"`
    DEBUG_MODE   string `validate:"omitempty"`
}
```

### 2. Carregar Variáveis

```go
var config Config
err := envloader.Load(&config)
if err != nil {
    log.Fatal(err)
}
```

### 3. Arquivo .env

```env
DB_URL=https://meudb.turso.io/database
API_KEY=minha-chave-secreta
DEBUG_MODE=true
```

## 🎯 Funcionalidades

### Tags Suportadas

| Tag        | Descrição                    | Exemplo                   |
| ---------- | ---------------------------- | ------------------------- |
| `env`      | Nome customizado da variável | `env:"DB_URL"`            |
| `validate` | Regras de validação          | `validate:"required,url"` |

### Validações Comuns

| Validação        | Descrição         | Exemplo                     |
| ---------------- | ----------------- | --------------------------- |
| `required`       | Campo obrigatório | `validate:"required"`       |
| `omitempty`      | Campo opcional    | `validate:"omitempty"`      |
| `url`            | URL válida        | `validate:"required,url"`   |
| `email`          | Email válido      | `validate:"email"`          |
| `min=10`         | Tamanho mínimo    | `validate:"min=10"`         |
| `oneof=dev prod` | Valor específico  | `validate:"oneof=dev prod"` |

## ✨ Interpolação de Variáveis

### Sintaxe

- **`[...]`** = Indica interpolação
- **`{VAR_NAME}`** = Referencia variável

### Exemplo

```env
# Variáveis base
DB_NAME=meu-app
DB_TOKEN=abc123

# Interpolação
DB_URL=[libsql://{DB_NAME}.turso.io?authToken={DB_TOKEN}]
API_URL=[https://{DB_NAME}.api.com/v1]
```

```go
type Config struct {
    DB_NAME  string `validate:"required"`
    DB_TOKEN string `validate:"required"`
    DB_URL   string `validate:"required,url"`   // Será interpolada
    API_URL  string `validate:"required,url"`   // Será interpolada
}
```

**Resultado:**

- `DB_URL` = `libsql://meu-app.turso.io?authToken=abc123`
- `API_URL` = `https://meu-app.api.com/v1`

## 🔧 Configurações Avançadas

### LoadOptions

```go
type LoadOptions struct {
    EnvFile           string // Caminho do arquivo .env
    RequiredByDefault bool   // Campos sem validate são required?
}
```

### Uso com Opções

```go
err := envloader.Load(&config, envloader.LoadOptions{
    EnvFile:           "config/production.env",
    RequiredByDefault: false,
})
```

## 📋 Exemplos Práticos

### Configuração de Banco

```go
type DatabaseConfig struct {
    HOST     string `validate:"required"`
    PORT     string `validate:"required"`
    NAME     string `validate:"required"`
    USER     string `validate:"required"`
    PASSWORD string `validate:"required"`

    // URL construída automaticamente
    URL string `validate:"required,url"`
}
```

```env
HOST=localhost
PORT=5432
NAME=meuapp
USER=admin
PASSWORD=secret123

URL=[postgres://{USER}:{PASSWORD}@{HOST}:{PORT}/{NAME}]
```

### Configuração de API

```go
type APIConfig struct {
    APP_NAME    string `validate:"required"`
    ENVIRONMENT string `validate:"required,oneof=dev staging prod"`
    VERSION     string `validate:"required"`

    BASE_URL    string `validate:"required,url"`
    HEALTH_URL  string `validate:"required,url"`
}
```

```env
APP_NAME=minha-api
ENVIRONMENT=prod
VERSION=v1

BASE_URL=[https://{APP_NAME}.{ENVIRONMENT}.com/{VERSION}]
HEALTH_URL=[{BASE_URL}/health]
```

### Configuração Turso/LibSQL

```go
type TursoConfig struct {
    DB_NAME     string `validate:"required"`
    DB_TOKEN    string `validate:"required"`
    REGION      string `validate:"required"`
    REPLICA_URL string `validate:"omitempty,url"`

    // URLs construídas
    PRIMARY_URL string `validate:"required,url"`
}
```

```env
DB_NAME=meu-banco
DB_TOKEN=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9...
REGION=aws-us-east-1

PRIMARY_URL=[libsql://{DB_NAME}-usuario.{REGION}.turso.io?authToken={DB_TOKEN}]
REPLICA_URL=[file:./local-{DB_NAME}.db]
```

## 🛠️ API Reference

### Funções Principais

#### `Load(target interface{}, opts ...LoadOptions) error`

Carrega e valida variáveis de ambiente na struct.

**Parâmetros:**

- `target`: Ponteiro para struct
- `opts`: Opções de configuração (opcional)

**Retorna:** `error` se houver problema

#### `MustLoad(target interface{}, opts ...LoadOptions)`

Como `Load()`, mas entra em pânico se houver erro.

**Exemplo:**

```go
// Desenvolvimento - falha rápida
envloader.MustLoad(&config)

// Produção - tratamento de erro
if err := envloader.Load(&config); err != nil {
    log.Fatal("Erro de configuração:", err)
}
```

### Tipos

#### `LoadOptions`

```go
type LoadOptions struct {
    EnvFile           string // Arquivo .env customizado
    RequiredByDefault bool   // Campos sem validate são required
}
```

## ⚠️ Tratamento de Erros

### Tipos de Erro

#### Variável Não Encontrada

```
variável de ambiente obrigatória não encontrada: DB_TOKEN
```

#### Erro de Validação

```
erros de validação:
- campo 'DATABASE_URL' deve ser uma URL válida
- campo 'API_KEY' falhou na validação 'min'
```

#### Erro de Interpolação

```
erro ao processar interpolação para DB_URL: variáveis não encontradas para interpolação: MISSING_VAR
```

### Exemplo de Tratamento

```go
func loadConfig() (*Config, error) {
    var config Config

    if err := envloader.Load(&config); err != nil {
        return nil, fmt.Errorf("falha ao carregar configuração: %w", err)
    }

    return &config, nil
}
```

## 🎨 Casos de Uso

### 1. Aplicação Web

```go
type WebConfig struct {
    PORT         string `validate:"required"`
    HOST         string `validate:"required"`
    SSL_ENABLED  string `validate:"omitempty"`

    SERVER_URL   string `validate:"required,url"`
}
```

```env
PORT=8080
HOST=localhost
SSL_ENABLED=false

SERVER_URL=[http://{HOST}:{PORT}]
```

### 2. Microserviços

```go
type ServiceConfig struct {
    SERVICE_NAME string `validate:"required"`
    NAMESPACE    string `validate:"required"`
    CLUSTER      string `validate:"required"`

    SERVICE_URL  string `validate:"required,url"`
    METRICS_URL  string `validate:"required,url"`
}
```

```env
SERVICE_NAME=user-api
NAMESPACE=production
CLUSTER=k8s-cluster

SERVICE_URL=[https://{SERVICE_NAME}.{NAMESPACE}.{CLUSTER}.com]
METRICS_URL=[{SERVICE_URL}/metrics]
```

### 3. Deploy Multi-Ambiente

```env
# .env.development
ENVIRONMENT=development
DB_HOST=localhost
API_HOST=localhost

# .env.production
ENVIRONMENT=production
DB_HOST=prod-db.company.com
API_HOST=api.company.com

# URLs construídas automaticamente
DB_URL=[postgres://user:pass@{DB_HOST}:5432/app]
API_URL=[https://{API_HOST}/v1]
```

## 📝 Boas Práticas

### 1. Organização de Variáveis

```go
type Config struct {
    // ✅ Agrupar por contexto

    // Database
    DB_NAME  string `validate:"required"`
    DB_TOKEN string `validate:"required"`
    DB_URL   string `validate:"required,url"`

    // API
    API_KEY    string `validate:"required,min=32"`
    API_TIMEOUT string `validate:"omitempty"`

    // App
    DEBUG_MODE string `validate:"omitempty"`
    LOG_LEVEL  string `validate:"omitempty,oneof=debug info warn error"`
}
```

### 2. Validações Específicas

```go
type Config struct {
    // ✅ Usar validações específicas
    EMAIL     string `validate:"required,email"`
    PORT      string `validate:"required,min=1,max=65535"`
    URL       string `validate:"required,url"`
    TIMEOUT   string `validate:"omitempty,min=1"`
    LOG_LEVEL string `validate:"omitempty,oneof=debug info warn error"`
}
```

### 3. Documentação das Variáveis

```go
type Config struct {
    // Database connection URL - formato: libsql://host?authToken=token
    DATABASE_URL string `env:"DB_URL" validate:"required,url"`

    // API key para serviços externos - mínimo 32 caracteres
    API_KEY string `validate:"required,min=32"`

    // Modo debug - valores: true/false
    DEBUG_MODE string `validate:"omitempty"`
}
```

## ⚡ Performance

- **Carregamento único**: Variáveis carregadas uma vez na inicialização
- **Cache interno**: Valores interpolados são resolvidos eficientemente
- **Validação rápida**: Usando `validator/v10` otimizado
- **Memory-safe**: Sem vazamentos de memória

## 🔒 Segurança

### Proteção de Secrets

```go
// ✅ Bom - secret só no environment
type Config struct {
    DB_TOKEN string `validate:"required"`
    DB_URL   string `validate:"required,url"`
}

// ❌ Evitar - secret no código
type Config struct {
    DB_TOKEN string `default:"token_hardcoded"` // Não fazer!
}
```

### Validação de URLs

```go
type Config struct {
    // ✅ URLs são validadas automaticamente
    API_URL string `validate:"required,url"`
    DB_URL  string `validate:"required,url"`
}
```

## 🐛 Debugging

### Verificar Interpolação

```go
// Adicionar logs temporários
fmt.Printf("Valor bruto: %s\n", os.Getenv("DB_URL"))
fmt.Printf("Valor interpolado: %s\n", config.DB_URL)
```

### Validar .env

```bash
# Verificar se arquivo existe
ls -la .env

# Ver conteúdo
cat .env

# Testar variáveis
echo $DB_NAME
```

---

## 📞 Suporte

Esta biblioteca é focada em **simplicidade** e **produtividade**. Para casos de uso específicos, você pode estender facilmente as funcionalidades.

**Funcionalidades principais:**

- ✅ Carregamento automático de .env
- ✅ Validação com tags
- ✅ Interpolação de variáveis
- ✅ Mensagens de erro claras
- ✅ Zero configuração para casos básicos

**Exemplo mínimo:**

```go
type Config struct {
    DATABASE_URL string `validate:"required,url"`
}

var config Config
envloader.MustLoad(&config)
```

**Pronto para usar!** 🚀
