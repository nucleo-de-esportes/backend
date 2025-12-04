# Fase 1: Setup
FROM golang:1.24 AS builder

# Variáveis para binário estático
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copia apenas go.mod/go.sum para cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN go build -o api ./cmd/nucleo_esportes/

# Fase 2: Imagem final
FROM debian:bookworm-slim AS final

# Criar diretório e usuário sem privilégios
RUN useradd -m appuser

WORKDIR /app

# Copiar apenas o binário
COPY --from=builder /app/api /app/api

# Permissões
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Definir o entrypoint
ENTRYPOINT ["./api", "-vars=exported"]
