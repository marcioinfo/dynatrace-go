### Camada de cartão - Payment Layer

### DevSecOps
Ao baixar o projeto pela primeira vez, é necessário baixar o `pre-commit` e instala-lo em sua maquina:
* MacOs:
`brew install pre-commit`
* Linux:
`pip install pre-commit`
Após a instalação, na raiz do projeto rodar o comando: 
`pre-commit install --config .pre-commit.yaml`

### Comandos
É necessário rodar as migrations na primeira execução do projeto para criar o banco de dados em sua maquina

* `docker compose up -d`
  Inicia os serviços app, banco de dados e de migations
* `docker compose down`
  Para os serviços app, banco de dados e de migrations
* `make migration-up`
  Roda todas as migrations
* `make migration-down`
  Desfaz a **última** migration
* `make migration-clear`
  Desfaz **todas** migrations
* `make migration-create name={NOME_DA_MIGRATION}`
  Cria uma nova migration com o nome passado no parâmetro `name`
* `make migration-fix version={VERSAO}`
  Comando utilizado para setar a versão passada no parâmetro `version` sem rodar as migrations
* `make shell-db`
  Comando utilizado para acessar o shell do banco de dados

