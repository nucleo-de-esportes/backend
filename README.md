
# Núcleo de Esportes (Backend)

# 📄 Documentação


- **Golang** com framework **GIN** para controle dos endpoints
- **Pacote golang-jwt** para autenticação de usuários com token jwt e controle de roles(Admin, Professor, Aluno)  
  - token armazenado nos cookies
- **Banco de dados postgreSQL**
- **GORM** para manipulação do banco de dados com mais facilidade
- **Aplicação divida em arquitetura em camadas**
- **Estrutura das pastas** seguindo os padrões encontrados aqui →  
  https://github.com/golang-standards/project-layout

  ## 🗂 Estrutura de Pastas do Projeto
  ├── cmd/
│ └── nucleo_esportes/
│ └── main.go # Ponto de entrada da aplicação (inicializa servidor)
│
├── docs/ # Documentação 
│
├── internal/ # Código interno da aplicação 
│ ├── config/ # Configurações gerais e inicialização 
│ ├── handlers/ # Recebem as requisições HTTP
│ ├── middleware/ # Middlewares (autenticação, CORS, etc.)
│ ├── model/ # Modelos, entidades 
│ ├── repository/ # Acesso ao banco de dados 
│ └── services/ # Regras de negócio 
│
└── scripts/
└── scripts.sql # Script SQL para estrutura do banco


- **Documentação dos endpoints** →  
  https://nucleo-de-esportes.github.io/backend/swagger/index.html



  


------------------------------------------------------------------------

## 🚀 Instalação

Clone o repositório:

``` sh
git clone https://github.com/nucleo-de-esportes/backend.git
cd backend
```

------------------------------------------------------------------------

## ⚙️ Configuração das variáveis de ambiente

O backend utiliza variáveis de ambiente para configurar o servidor e o
banco de dados.\
Essas variáveis podem ser definidas em um arquivo `.env` na raiz do
projeto ou exportadas diretamente no ambiente do sistema.

### Exemplo de `.env`

``` env
# Configurações do Banco de Dados
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT="5432"
DB_NAME=nucleo_esportes
DB_USER=postgres
DB_PASSWORD=senha_segura

# Configuração do Servidor
SERVER_PORT="8080"

# Autenticação JWT
SECRET_KEY=sua-chave-secreta-aqui
```

### Descrição das variáveis

| Variável        | Obrigatória  | Padrão      | Descrição                                                                        |
| --------------- | ------------ | ----------- | -------------------------------------------------------------------------------- |
| **DB_DRIVER**   | ✅ Sim       | —           | Define o driver do banco de dados (ex: `postgres`).                              |
| **DB_HOST**     | ❌ Não       | `localhost` | Endereço do host do banco de dados.                                              |
| **DB_PORT**     | ❌ Não       | `5432`      | Porta de conexão com o banco de dados.                                           |
| **DB_NAME**     | ✅ Sim       | —           | Nome do banco de dados.                                                          |
| **DB_USER**     | ❌ Não       | `postgres`  | Usuário do banco de dados.                                                       |
| **DB_PASSWORD** | ✅ Sim       | —           | Senha do banco de dados.                                                         |
| **SERVER_PORT** | ❌ Não       | `8000`      | Porta em que o servidor irá escutar as requisições.                              |
| **SECRET_KEY**  | ✅ Sim       | —           | Chave secreta para assinatura e validação de tokens JWT. Use uma string segura.  |


> ⚠️ Caso alguma variável obrigatória não seja definida, a aplicação
> será encerrada com erro no carregamento da configuração.

------------------------------------------------------------------------

## 🧩 Executando o projeto

### Compilar o código

``` sh
go build ./...
```

### Executar em modo de desenvolvimento

``` sh
go run ./...
```

## 🐳 Executando com Docker Compose

Para subir a API e o banco de dados de forma simplificada, utilize Docker Compose.

### Passos

1. Crie o arquivo .env na raiz do projeto (conforme exemplo acima).

2. Suba os containers:

``` sh
docker compose up --build
```

3. A API ficará acessível em:

http://localhost:8080


4. O PostgreSQL estará disponível em:

``` yml
Host: localhost
Port: 5432
Database: nucleo_esportes
User: postgres
Password: senha_segura
```

------------------------------------------------------------------------

## 🧠 Estrutura de Configuração (Go)

O carregamento das variáveis de ambiente é feito no pacote `config`
usando a biblioteca [caarlos0/env](https://github.com/caarlos0/env).

``` go
func LoadConfig() *Config {
    var cfg Config

    if err := env.Parse(&cfg); err != nil {
        log.Fatalf("erro ao carregar variáveis de ambiente: %v", err)
    }

    return &cfg
}
```

