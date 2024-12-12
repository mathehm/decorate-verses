# Use a imagem oficial do Go como base (versão 1.23 ou superior)
FROM golang:1.23 AS builder

# Configure o diretório de trabalho
WORKDIR /app

# Copie os arquivos do projeto
COPY . .

# Baixe as dependências e compile o binário
RUN go mod tidy && go build -o main .

# Use a mesma imagem para execução
FROM golang:1.23

# Copie o binário da etapa de build
COPY --from=builder /app/main /

# Exponha a porta 8080
EXPOSE 8080

# Comando para executar a aplicação
CMD ["/main"]
