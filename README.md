# Triador AIIA

Sistema para análise de aderência entre currículo e vaga utilizando IA.

O projeto recebe o texto de um currículo e uma descrição de vaga, envia os dados para uma LLM, gera uma análise estruturada e salva o histórico no banco de dados.

---

## Funcionalidades

- Análise automática entre currículo e vaga
- Extração de skills técnicas
- Cálculo de score de aderência entre 0 e 100
- Estimativa de experiência profissional
- Geração de resumo explicando a pontuação
- Persistência das análises em banco relacional
- Consulta de histórico
- Interface web funcional com loading visual
- Deploy funcional

---

## Tecnologias

### Backend

- Go
- Gin
- PostgreSQL
- GORM
- OpenAI API
- Docker Compose

### Frontend

- Next.js
- TypeScript
- Tailwind CSS

### Infraestrutura

- Railway
- Vercel

---

## Arquitetura

Backend organizado em camadas:

```txt
backend
├── cmd/api
├── internal
│   ├── database
│   ├── dto
│   ├── handler
│   ├── llm
│   ├── model
│   ├── repository
│   └── service
```

Fluxo principal:

```txt
Frontend
    ↓
POST /analyses
    ↓
Handler
    ↓
Service
    ↓
Prompt
    ↓
OpenAI
    ↓
Sanitização
    ↓
Validação
    ↓
Repository
    ↓
PostgreSQL
    ↓
Response
```

---

## Decisões técnicas

### Escolha do Go

Escolhi Go para o backend por ser uma linguagem performática, simples para criação de APIs HTTP e adequada para serviços pequenos e objetivos. Como o desafio valoriza separação de responsabilidades e clareza técnica, Go permitiu construir uma API enxuta, com camadas bem definidas e sem excesso de abstrações.

### Uso do PostgreSQL

Foi escolhido PostgreSQL por ser um banco relacional robusto e amplamente usado em aplicações reais. Apesar de SQLite também atender ao desafio, PostgreSQL deixa o ambiente mais próximo de uma aplicação backend em produção.

### Uso do GORM

O GORM foi utilizado para reduzir boilerplate de persistência, facilitar a criação da entidade de análise e permitir automigration no ambiente local. Para um ambiente de produção, migrations versionadas poderiam ser consideradas.

### Modelagem de skills

As skills são armazenadas como JSON serializado em uma coluna de texto. Essa decisão reduz complexidade relacional para o escopo atual, já que o sistema não exige filtros ou relatórios por skill.

Trade-off:

**Vantagem**
- Menor complexidade

**Desvantagem**
- Dificulta consultas futuras por skill específica

Em uma evolução do projeto as skills poderiam ser normalizadas em tabela própria.

### Integração com LLM

A integração com OpenAI foi isolada na camada `llm`, evitando acoplamento direto do service com detalhes do SDK externo. O backend não confia cegamente na resposta do modelo: a saída é convertida para struct, sanitizada e validada antes da persistência.

Fluxo aplicado:

- Prompt estruturado
- Sanitização
- Conversão para struct
- Validação
- Persistência

---

## Pré-requisitos

Antes de executar:

- Go
- Node.js
- npm
- Docker Desktop
- Chave OpenAI

---

# Como executar localmente

## Clonar repositório

```bash
git clone https://github.com/BryanPinheiro77/triador-aiia.git

cd triador-aiia
```

---

# Backend

Entrar na pasta:

```bash
cd backend
```

Subir banco:

```bash
docker compose up -d
```

Criar `.env`

```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=triador port=5432 sslmode=disable
OPENAI_API_KEY=sua_chave_openai
OPENAI_MODEL=gpt-4o-mini
FRONTEND_URL=http://localhost:3000
```

Instalar dependências:

```bash
go mod tidy
```

Executar:

```bash
go run cmd/api/main.go
```

API:

```txt
http://localhost:8080
```

Health check:

```txt
GET /health
```

---

# Frontend

Abrir novo terminal:

```bash
cd frontend
```

Criar:

`.env.local`

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Instalar:

```bash
npm install
```

Executar:

```bash
npm run dev
```

Aplicação:

```txt
http://localhost:3000
```

---

## Deploy

Projeto publicado utilizando:

- Frontend: Vercel
- Backend: Railway
- Banco: PostgreSQL no Railway

Demonstração pública:

https://triador-aiia.vercel.app

---

## Endpoints

### POST /analyses

Request:

```json
{
  "resume":"Texto currículo",
  "job_description":"Texto vaga"
}
```

Response:

```json
{
  "id":1,
  "candidate_name":"Nome",
  "skills":["Java","Spring Boot"],
  "years_experience":2,
  "fit_score":85,
  "summary":"Resumo justificando a pontuação."
}
```

---

### GET /analyses

Lista histórico.

---

## Validações

Backend implementa:

- Validação de request
- Sanitização da resposta da LLM
- Conversão para struct interna
- Validação de campos obrigatórios
- Validação do score entre 0–100
- Validação de experiência
- Tratamento de erro do provedor
- Tratamento de JSON inválido
- CORS configurado por ambiente

---

## Testes

### Backend

```bash
go fmt ./...

go test ./...
```

Testes implementados:

- Sanitização de respostas da LLM
- Remoção de blocos markdown
- Preservação de JSON válido
- Validação do comportamento esperado do núcleo de tratamento da IA

### Frontend

```bash
npm run lint

npm run build
```

---

## Limitações atuais

- Não há autenticação
- Não há upload de PDF
- Não há streaming da resposta da IA
- Skills não são normalizadas
- Não há paginação no histórico'

---

## Melhorias futuras

- Upload de PDF
- Retry com backoff para rate limit
- Normalização de skills
- Paginação
- Comparação múltipla de vagas
- Melhorias visuais adicionais
- Separação entre skills encontradas e ausentes

---

## Autor

Bryan Mendes Pinheiro da Silva
