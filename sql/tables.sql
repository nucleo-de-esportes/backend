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
CREATE TABLE modalidade (
    modalidade_id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    valor_aluno NUMERIC(10,2) NOT NULL,
    valor_professor NUMERIC(10,2) NOT NULL,
    valor_externo NUMERIC(10,2) NOT NULL
);

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

CREATE TABLE matricula (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    turma_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_usuario FOREIGN KEY (user_id) REFERENCES usuario(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_turma FOREIGN KEY (turma_id) REFERENCES turma(turma_id) ON DELETE CASCADE,
    CONSTRAINT uq_matricula UNIQUE (user_id, turma_id) -- evita matrícula duplicada
);

ALTER TABLE turma
ADD COLUMN professor_id UUID,
ADD CONSTRAINT fk_professor FOREIGN KEY (professor_id) REFERENCES usuario(user_id) ON DELETE SET NULL;

ALTER TABLE usuario
ADD COLUMN password TEXT NOT NULL;


CREATE TABLE aula (
    id SERIAL PRIMARY KEY,
    turma_id BIGINT NOT NULL,
    data_hora TIMESTAMP NOT NULL,
    criado_em TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_aula_turma FOREIGN KEY (turma_id) REFERENCES turma(id) ON DELETE CASCADE
);


CREATE TABLE presenca (
    id SERIAL PRIMARY KEY,
    aula_id BIGINT NOT NULL,
    user_id UUID NOT NULL,
    presente BOOLEAN NOT NULL DEFAULT FALSE,
    criado_em TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_presenca_aula FOREIGN KEY (aula_id) REFERENCES aula(id) ON DELETE CASCADE,
    CONSTRAINT fk_presenca_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_presenca UNIQUE (aula_id, user_id)  -- 1 aluno só tenha 1 registro por aula
);

/**
AJUSTANDO O INCREMENT DOS IDS DAS TURMAS, POIS ESTAVA COMEÇANDO COM TURMA DE ID 0
**/
ALTER SEQUENCE turma_turma_id_seq RESTART WITH 1;




























