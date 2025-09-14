# --- Estágio 1: Builder ---
# Usamos uma imagem completa do Go para compilar nossa aplicação
FROM golang:1.23-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de dependência primeiro para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código fonte
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 cria um binário estático, essencial para imagens mínimas.
# GOOS=linux garante que seja compilado para o kernel do Linux (usado em containers).
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /defender-vehicle ./cmd/server

# --- Estágio 2: Runner ---
# Usamos uma imagem mínima para a versão final, apenas para rodar o app compilado
FROM alpine:latest

WORKDIR /app

# Copia APENAS o binário compilado do estágio 'builder'
COPY --from=builder /defender-vehicle .

# Copia APENAS os arquivos estáticos (HTML/JS) do estágio 'builder'
COPY --from=builder /app/static ./static

# Expõe a porta que nossa aplicação usa dentro do container
EXPOSE 8080

# Comando para iniciar a aplicação quando o container for executado
CMD ["./defender-vehicle"]