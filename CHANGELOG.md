## [1.5.0](https://github.com/nucleo-de-esportes/backend/compare/v1.4.2...v1.5.0) (2025-11-26)


### Features

* CreateAviso function, need to fix some errors ([0708b14](https://github.com/nucleo-de-esportes/backend/commit/0708b147f436d8b04dfc4287171d2aa9353c083e))


### Bug Fixes

* fixing CreateAviso function ([147109e](https://github.com/nucleo-de-esportes/backend/commit/147109e13c34d1f5ae1cb0903382358856593c01))

## [1.4.2](https://github.com/nucleo-de-esportes/backend/compare/v1.4.1...v1.4.2) (2025-11-12)


### Bug Fixes

* tipo de uuid de turmas ([07e2a5d](https://github.com/nucleo-de-esportes/backend/commit/07e2a5d49c4ccfb824c6b35604201ac0a3e0ff1c))

## [1.4.1](https://github.com/nucleo-de-esportes/backend/compare/v1.4.0...v1.4.1) (2025-11-12)


### Bug Fixes

* merge conflict in Turma struct definition ([c7131ad](https://github.com/nucleo-de-esportes/backend/commit/c7131ad26f4b7bb36e3a0c8b376515179388d22e))

## [1.4.0](https://github.com/nucleo-de-esportes/backend/compare/v1.3.1...v1.4.0) (2025-11-04)


### Features

* UpdateUser implementada para admin alterar user_type ([b0dfcc1](https://github.com/nucleo-de-esportes/backend/commit/b0dfcc164dd50bb47d2a5f8dfa215b1d6ef29df6))


### Bug Fixes

* fixing getTurmasByUser response ([8023704](https://github.com/nucleo-de-esportes/backend/commit/8023704454196b6a8be225d1c99d9174e1432416))
* fixing more json structures response ([bdc4b08](https://github.com/nucleo-de-esportes/backend/commit/bdc4b08169d3cb0b8b98c44e78d104dbea11b0b3))
* fixing responses structure ([91d18d1](https://github.com/nucleo-de-esportes/backend/commit/91d18d189b7533ce4ff9a2fd43d3faafab35398d))

## [1.3.1](https://github.com/nucleo-de-esportes/backend/compare/v1.3.0...v1.3.1) (2025-10-30)


### Bug Fixes

* ajustando token de usuario ([d8e7f29](https://github.com/nucleo-de-esportes/backend/commit/d8e7f295cf1709bba69a48dda83c9e4378083ded))

## [1.3.0](https://github.com/nucleo-de-esportes/backend/compare/v1.2.0...v1.3.0) (2025-10-29)


### Features

* implementado endpoint GET /professor/{id_professor}/aulas ([dac2a0f](https://github.com/nucleo-de-esportes/backend/commit/dac2a0fa6ff2a004b0f49fb49bc5e6da32edc860))


### CI/CD

* ajustado path do arquivo main.go no workflow de deploy do swagger ([58f34b9](https://github.com/nucleo-de-esportes/backend/commit/58f34b980a2f762586f2735ab98ace740dcf325c))


### Chores

* ajuste de timezone dos serviços no docker compose ([6efae81](https://github.com/nucleo-de-esportes/backend/commit/6efae81d2806e537445e8ea0d42093f843c1f324))

## [1.2.0](https://github.com/nucleo-de-esportes/backend/compare/v1.1.0...v1.2.0) (2025-10-29)


### Features

* implementado endpoint GET /turma/{id_turma}/alunos ([8fb40d4](https://github.com/nucleo-de-esportes/backend/commit/8fb40d4e0f1e6365a13f693c5b22667f15d3d868))


### Bug Fixes

* corrigido conversão de uuid durante criação do jwt e adicionado variável SECRET_KEY ao docker-compose.yml ([adcba38](https://github.com/nucleo-de-esportes/backend/commit/adcba38651fbcaa877c5bdf76fcdf93cfecabf33))


### Chores

* atualizado documentação das variáveis de ambiente do README.md ([4e23c1b](https://github.com/nucleo-de-esportes/backend/commit/4e23c1bd7dad550a4b146541db5b054f6aa8d245))

## [1.1.0](https://github.com/nucleo-de-esportes/backend/compare/v1.0.3...v1.1.0) (2025-10-23)


### Features

* Instancia de aulas criada, modificacao na estrutura de modalidades com dias da semana ([9d19eb0](https://github.com/nucleo-de-esportes/backend/commit/9d19eb0a9ac1260ed825d3228abe07d8438c57ba))


### Bug Fixes

* altera pacote de uuid para gorm.io/datatypes visando compatibilidade com gorm ([d2d9377](https://github.com/nucleo-de-esportes/backend/commit/d2d9377666fe16e95c1cfa0e0df8a8cc59016b5c))


### Chores

* movido diretório sql para scripts ([bd431bc](https://github.com/nucleo-de-esportes/backend/commit/bd431bc939331a17528d7b8991b7f108653ba48b))
* remove branch de teste da configuração de release ([7c12943](https://github.com/nucleo-de-esportes/backend/commit/7c12943332e3bf5cf1d7b6096b2635392e0f18c1))

## [1.0.3](https://github.com/nucleo-de-esportes/backend/compare/v1.0.2...v1.0.3) (2025-10-19)


### Refactoring

* unifica release e build em um único workflow ([95617eb](https://github.com/nucleo-de-esportes/backend/commit/95617eb0e2e8468644f1becb282e02d3f215f9c9))

## [1.0.2](https://github.com/nucleo-de-esportes/backend/compare/v1.0.1...v1.0.2) (2025-10-19)


### Bug Fixes

* adiciona trigger para pre-releases ([740c190](https://github.com/nucleo-de-esportes/backend/commit/740c190f173c640e7ca3582e1ebaf2855ae05e31))

## [1.0.1](https://github.com/nucleo-de-esportes/backend/compare/v1.0.0...v1.0.1) (2025-10-19)


### Bug Fixes

* adiciona trigger de release para acionar build automaticamente ([4a1eb84](https://github.com/nucleo-de-esportes/backend/commit/4a1eb8431ac2788124dc82becdde86ee7bc196f5))

## 1.0.0 (2025-10-19)


### Features

* adicionando handler e endpoint para confirmação de presenca de aluno ([422813d](https://github.com/nucleo-de-esportes/backend/commit/422813d1b2376ffc073d1c7cee8bc759a07e4bf8))
* CreateAula implemented for GetNextAula method ([9f937ae](https://github.com/nucleo-de-esportes/backend/commit/9f937ae220b2795263217d15451e9013a556aa65))
* Criacao do metodo DELETE para turma ([ff247c5](https://github.com/nucleo-de-esportes/backend/commit/ff247c5a05aa2366cc65f1e2a1aefd16ee1ac5a5))
* criacao dos model e tabelas no banco de dados para suporte a aulas e presencas ([92239eb](https://github.com/nucleo-de-esportes/backend/commit/92239eb4b7b980850c81742636b36ac3d0bb0a41))
* criado factory para banco de dados ([c03b5e7](https://github.com/nucleo-de-esportes/backend/commit/c03b5e7a49f4665675d841df32315379d67f0c56))
* criado func para carregar config do .env ([6441898](https://github.com/nucleo-de-esportes/backend/commit/64418980dadc4be8d78ca440e5bd68585a1e4369))
* criado package de config ([bcd8f3d](https://github.com/nucleo-de-esportes/backend/commit/bcd8f3d9b5d665153503330459915afff323818c))
* delete user e delete user de turma ([c4a4aca](https://github.com/nucleo-de-esportes/backend/commit/c4a4aca7971bd391d81c5cfa2d7c63b57eaf8d93))
* deleteUserById e delete user de uma turma ([9840722](https://github.com/nucleo-de-esportes/backend/commit/98407223f2a84215c3c58f942e8b331319747855))
* function GetTurmasByUser implemented ([fbd3a50](https://github.com/nucleo-de-esportes/backend/commit/fbd3a5031042ecd2b0c107aa934e4b0be598a4bc))
* get proxima aula de turma implementada ([c5b3ba7](https://github.com/nucleo-de-esportes/backend/commit/c5b3ba7d92058381d0d3770924fc00432ea0dfc7))
* GetAllTurmas method ([e7a158d](https://github.com/nucleo-de-esportes/backend/commit/e7a158dbf24bd6b08fd545583c97feb444055dfd))
* GetNextClassById implemented ([e674f39](https://github.com/nucleo-de-esportes/backend/commit/e674f39e38621f5b857adb9036b8deda39046d04))
* Login function created ([d527824](https://github.com/nucleo-de-esportes/backend/commit/d527824381af96339ef445ee819e70ce1939fa2d))
* method RegisterUser implemented ([d8edb06](https://github.com/nucleo-de-esportes/backend/commit/d8edb06e426abdb74dd9a3742b7f7330ec9fc77c))
* method updateTurma ([00219b9](https://github.com/nucleo-de-esportes/backend/commit/00219b92435bc5a4b0756e2097a855bc2a7a9572))


### Bug Fixes

* adicionado atributos Port e Name a struct de config do banco de dados ([7e18dd4](https://github.com/nucleo-de-esportes/backend/commit/7e18dd43ec3e0cd177c2151417fe84e259f687b5))
* ajusta permissões do workflow seguindo padrão oficial ([3ccf3bf](https://github.com/nucleo-de-esportes/backend/commit/3ccf3bfd9b88cd95ea77fdecd5b6a0b8cdcc3a79))
* atualizado script Dockerfile ([9cfe875](https://github.com/nucleo-de-esportes/backend/commit/9cfe875deb49d67d86c5d12a7c75f577afeb576c))
* correct data flow ([b202e7e](https://github.com/nucleo-de-esportes/backend/commit/b202e7eed5f29e99448152bd6be8a4a9332cbf20))
* corrige autenticação do semantic-release no workflow ([c197602](https://github.com/nucleo-de-esportes/backend/commit/c19760200a58dbb6f7f123c3778a92ad783239a3))
* corrige instalação de dependências no workflow de release ([53d3ff3](https://github.com/nucleo-de-esportes/backend/commit/53d3ff3e548815aa7470d3967a703dcc3947d674))
* estado de presença mais dinamico ([52e435d](https://github.com/nucleo-de-esportes/backend/commit/52e435d30cde41157b70ac84082ac103c0759fa1))
* fixing function ([b1f6e9c](https://github.com/nucleo-de-esportes/backend/commit/b1f6e9c1178a3dea47ed3ed2870989a710ca4056))
* fixing json structure ([2d445f7](https://github.com/nucleo-de-esportes/backend/commit/2d445f7aeab322dac88522f07a6b7e916f5ecb9f))
* fixing response on GetTurmasByUser method ([745185c](https://github.com/nucleo-de-esportes/backend/commit/745185cd445ca8fe42a2f8b4e0efd74fbf5fb61a))
* fixing syntax error ([600ca27](https://github.com/nucleo-de-esportes/backend/commit/600ca273f12d388ac342a2245012d15192889d2c))
* Funcoes GET agora retornam id da turma ([37062b5](https://github.com/nucleo-de-esportes/backend/commit/37062b52a67ca44ac939e96359184e5e8557f380))
* implementing time ([63ec2d1](https://github.com/nucleo-de-esportes/backend/commit/63ec2d1ebc6bec11530403a2e2eb4c7122d48ad2))
* removing field user_type from register request's body ([f36c215](https://github.com/nucleo-de-esportes/backend/commit/f36c21533b493fbf117a7d12a872774432a58a0b))
* removing field user_type from register request's body ([73be686](https://github.com/nucleo-de-esportes/backend/commit/73be6868c98a76a99d4ae9dd7a0223435a5790c5))
* removing user_type field from response body and adding it to the JWT token ([ea25ddb](https://github.com/nucleo-de-esportes/backend/commit/ea25ddba0065fb0c3b0bedc160efb8e41e0b2d62))
* utilizando config para init do banco e server ([c68d666](https://github.com/nucleo-de-esportes/backend/commit/c68d666b46456dff6fb9d81949252c23bdc2464b))


### Refactoring

* alterado forma de carregamento das variáveis de ambiente ([376275d](https://github.com/nucleo-de-esportes/backend/commit/376275de77326d1d2611a866a9dabd355d715154))
* estrutura das postas e alguns codigos modificados ([ed65478](https://github.com/nucleo-de-esportes/backend/commit/ed65478398772a82f06afa30b1fe094aebba21b2))
* func InscreverAluno and GetTurmasByUser ([8900b8b](https://github.com/nucleo-de-esportes/backend/commit/8900b8bd9ab229a1a6313d54c877350743ef830f))
* func RegisterUser using GORM instead of supabase ([b3641ba](https://github.com/nucleo-de-esportes/backend/commit/b3641baa86f8584c5abc139ac1f00a4b3f2e5a6a))
* userController ([36bd494](https://github.com/nucleo-de-esportes/backend/commit/36bd494db813251c74de873dc3a5c3aae665176c))


### CI/CD

* adiciona versionamento automático com semantic-release ([8056c64](https://github.com/nucleo-de-esportes/backend/commit/8056c64c9bd85753c5892c267e806c509de53dd2))
* ajuste de branch no yml do github ([2551c9b](https://github.com/nucleo-de-esportes/backend/commit/2551c9bc23720bcf120b6f8f4f4c34388cafb637))
* alterado branchs de execução do job de build ([ff30599](https://github.com/nucleo-de-esportes/backend/commit/ff3059983d0e7c9a877e2e459e29d3e98be7dbb5))
* alterado pipeline para rodar apenas na branch main ([f7d377e](https://github.com/nucleo-de-esportes/backend/commit/f7d377ee647d92d255d2d7da35335c6a8504bd1e))
* alterado pipeline para rodar apenas na branch main ([1edf84d](https://github.com/nucleo-de-esportes/backend/commit/1edf84dff286d4489a234ea122be8b1be6c6c231))
* alterado variável de ci IMAGE_NAME ([07d93af](https://github.com/nucleo-de-esportes/backend/commit/07d93afcb88a29cf6296328d22020b8754cf92f1))
* alterado variável de ci IMAGE_NAME ([5db8e15](https://github.com/nucleo-de-esportes/backend/commit/5db8e1593a316ba35e49c6cf77f861fbd236384e))
* alterado variável de ci IMAGE_NAME ([e341d79](https://github.com/nucleo-de-esportes/backend/commit/e341d79232fabf78c01b4df7282516f955377d8a))
* configurações iniciais para deploy ([9f96f71](https://github.com/nucleo-de-esportes/backend/commit/9f96f715a659f1aa07cf9bb1bc080ada9532dad1))
* criado job para build de imagem docker ([d3789a8](https://github.com/nucleo-de-esportes/backend/commit/d3789a8ecbca6127181e4b4e55dc18f8fafcad91))
* removido tags desnecessárias da imagem docker ([4eb03c7](https://github.com/nucleo-de-esportes/backend/commit/4eb03c7d7c41fea5e880ebd639723a854c752acd))


### Chores

* adicion documentação de variáveis de ambiente ao README ([81cdc8a](https://github.com/nucleo-de-esportes/backend/commit/81cdc8a3d1958de236d972458f1c5a0e7f5617ab))
* ajuste de nome do diretório 'scripts' ([8008235](https://github.com/nucleo-de-esportes/backend/commit/8008235cd77b9e604056e18b554339e9d59eba97))
* atualizado README para documentar docker compose ([1bf0daa](https://github.com/nucleo-de-esportes/backend/commit/1bf0daa0e2c4d7029a75c371eaa82b6984bb5490))
* criado arquivo docker-compose.yml ([6dde41e](https://github.com/nucleo-de-esportes/backend/commit/6dde41e467360b9901a5f8fffc9d809872056d7d))
* removido config fly antiga ([b769039](https://github.com/nucleo-de-esportes/backend/commit/b7690398a87d5116c9d972476976bf3e2683ddfa))

# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Adicionado
- Versionamento automático com semantic-release
- Geração automática de CHANGELOG
- Build automático de imagem Docker ao criar tags
