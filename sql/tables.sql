-- Criando tabela USUARIO
CREATE TABLE usuario (
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    user_type TEXT NOT NULL,
    nome TEXT NOT NULL
);

-- Criando tabela LOCAL
CREATE TABLE local (
    local_id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    campus VARCHAR(100) NOT NULL
);

-- Criando tabela TURMA
CREATE TABLE turma (
    turma_id SERIAL PRIMARY KEY,
    horario_inicio TEXT NOT NULL,
    horario_fim TEXT NOT NULL,
    limite_inscritos INT NOT NULL,
    dia_semana VARCHAR(20) NOT NULL,
    sigla VARCHAR(20) NOT NULL,
    local_id INT NOT NULL,
    modalidade_id INT NOT NULL,
    CONSTRAINT fk_local FOREIGN KEY (local_id) REFERENCES local(local_id) ON DELETE CASCADE,
    CONSTRAINT fk_modalidade FOREIGN KEY (modalidade_id) REFERENCES modalidade(modalidade_id) ON DELETE CASCADE
);

-- Criando tabela MATRICULA
CREATE TABLE matricula (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    turma_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_usuario FOREIGN KEY (user_id) REFERENCES usuario(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_turma FOREIGN KEY (turma_id) REFERENCES turma(turma_id) ON DELETE CASCADE,
    CONSTRAINT uq_matricula UNIQUE (user_id, turma_id) -- evita matrícula duplicada
);


CREATE TABLE modalidade (
    modalidade_id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    valor_aluno NUMERIC(10,2) NOT NULL,
    valor_professor NUMERIC(10,2) NOT NULL,
    valor_externo NUMERIC(10,2) NOT NULL
);

ALTER TABLE turma
ADD COLUMN professor_id UUID,
ADD CONSTRAINT fk_professor FOREIGN KEY (professor_id) REFERENCES usuario(user_id) ON DELETE SET NULL;

ALTER TABLE usuario
ADD COLUMN password TEXT NOT NULL;