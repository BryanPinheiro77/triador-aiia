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
- Interface web funcional com estado de carregamento

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
Validação da resposta
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

As skills são armazenadas como JSON serializado em uma coluna de texto. Essa decisão reduz complexidade relacional para o escopo atual, já que o sistema não exige filtros ou relatórios por skill. Em uma evolução do produto, as skills poderiam ser normalizadas em uma tabela própria.

### Integração com LLM

A integração com OpenAI foi isolada na camada `llm`, evitando acoplamento direto do service com detalhes do SDK externo. O backend não confia cegamente na resposta do modelo: a saída é convertida para struct, sanitizada e validada antes da persistência.

---

## Pré-requisitos

Antes de executar o projeto, é necessário ter instalado:

- Go
- Node.js
- npm
- Docker Desktop
- Conta/chave de API da OpenAI

---

## Como executar localmente

### 1. Clonar o repositório

```bash
git clone https://github.com/BryanPinheiro77/triador-aiia.git
cd triador-aiia
```

---

## Backend

### 2. Entrar na pasta do backend

```bash
cd backend
```

### 3. Subir o PostgreSQL com Docker Compose

Como o arquivo `docker-compose.yml` está dentro da pasta `backend`, execute:

```bash
docker compose up -d
```

Esse comando sobe um container PostgreSQL local para a aplicação.

### 4. Criar o arquivo `.env`

Crie um arquivo chamado `.env` dentro da pasta `backend`.

Exemplo:

```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=triador port=5432 sslmode=disable
OPENAI_API_KEY=sua_chave_openai
OPENAI_MODEL=gpt-4o-mini
```

Também existe um arquivo `.env.example` com as variáveis necessárias, sem valores sensíveis.

### 5. Instalar dependências do backend

```bash
go mod tidy
```

### 6. Executar o backend

```bash
go run cmd/api/main.go
```

A API ficará disponível em:

```txt
http://localhost:8080
```

Endpoint de teste:

```txt
GET http://localhost:8080/health
```

---

## Frontend

Abra outro terminal a partir da raiz do projeto.

### 7. Entrar na pasta do frontend

```bash
cd frontend
```

### 8. Criar o arquivo `.env.local`

Crie um arquivo chamado `.env.local` dentro da pasta `frontend`.

Exemplo:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

Também existe um arquivo `.env.example` com a variável necessária.

### 9. Instalar dependências do frontend

```bash
npm install
```

### 10. Executar o frontend

```bash
npm run dev
```

A aplicação ficará disponível em:

```txt
http://localhost:3000
```

---

## Endpoints principais

### Criar análise

```txt
POST /analyses
```

Body:

```json
{
  "resume": "Texto do currículo",
  "job_description": "Texto da vaga"
}
```

Resposta:

```json
{
  "id": 1,
  "candidate_name": "Nome do candidato",
  "skills": ["Java", "Spring Boot"],
  "years_experience": 2,
  "fit_score": 85,
  "summary": "Resumo justificando a nota."
}
```

### Listar histórico

```txt
GET /analyses
```

---

## Validações e tratamento de erro

O backend implementa:

- Validação de campos obrigatórios no request
- Sanitização da resposta da LLM
- Conversão da resposta do modelo para struct interna
- Validação de campos obrigatórios da resposta da LLM
- Validação de `fit_score` entre 0 e 100
- Validação de experiência não negativa
- Tratamento de erro para falha do provedor LLM
- Tratamento de resposta inválida da LLM
- CORS configurado para o frontend local

---

## Variáveis de ambiente

### Backend

```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=triador port=5432 sslmode=disable
OPENAI_API_KEY=sua_chave_openai
OPENAI_MODEL=gpt-4o-mini
```

### Frontend

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Testes e validação local

Comandos utilizados durante o desenvolvimento:

### Backend

```bash
go fmt ./...
go test ./...
```

### Frontend

```bash
npm run lint
npm run build
```

---

## Limitações atuais

- Não há autenticação
- Não há upload de PDF
- Não há deploy público configurado
- Não há streaming da resposta da IA
- Testes automatizados ainda não foram adicionados
- O prompt pode ser refinado em versões futuras conforme novos cenários de avaliação

---

## Possíveis melhorias futuras

- Upload de PDF com extração de texto
- Testes automatizados para validação da saída da LLM
- Retry com backoff para rate limit do provedor
- Deploy do backend e frontend
- Análise de múltiplas vagas em uma única requisição
- Separação entre skills encontradas e skills ausentes
- Melhorias visuais no histórico e no score

---

## Autor

Bryan Mendes Pinheiro da Silva
