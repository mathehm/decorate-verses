# Use a imagem oficial do Go como base (versão 1.23 ou superior)
FROM golang:1.23 AS builder

# Configure o diretório de trabalho
WORKDIR /app

# Copie apenas os arquivos necessários para o Go (go.mod e go.sum)
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod tidy

# Copie os arquivos restantes do projeto
COPY . .

# Compile o binário
RUN go build -o main .

# Use a mesma imagem para execução
FROM gcr.io/distroless/base

# Copie o binário da etapa de build
COPY --from=builder /app/main /app/

# Exponha a porta 8080 (Cloud Run requer que o container escute em uma porta configurada)
EXPOSE 8080

# Comando para executar a aplicação
CMD ["/app/main"]
