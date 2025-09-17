# Fase 1: Setup
FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias
RUN apk --no-cache add git
RUN apk --no-cache add ca-certificates

RUN update-ca-certificates

# Diretório de trabalho
WORKDIR /app

# Copiar os arquivos de mod para baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código da aplicação
COPY . .

# Criar o diretório bin para armazenar o binário
RUN mkdir -p /app/bin

# Fase 3: Application Build
FROM builder AS build

# Compilar a aplicação Go para um binário pequeno e otimizado
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -o /app/bin/api ./main

# Fase 3: Imagem final
FROM scratch AS final

# Copiar o binário da fase de build
COPY --from=build /app/bin/api /app/api
COPY --from=build /etc/ssl/certs /etc/ssl/certs

# Configurar permissões e diretório de trabalho
WORKDIR /app

# Expor a porta da API
EXPOSE 8080

# Definir o entrypoint
ENTRYPOINT ["./api", "-vars=exported"]
