# 🐳 Flyway - Controle de Versões do Banco de Dados

## 📘 Visão Geral

O **Flyway** é uma ferramenta de **migração de banco de dados** que aplica scripts SQL de forma incremental e controlada, garantindo que todos os ambientes (desenvolvimento, homologação e produção) mantenham o mesmo estado de schema.

Ele funciona de forma semelhante ao versionamento de código: cada script SQL tem uma versão (`V1`, `V2`, `V3`, …), e o Flyway mantém o histórico do que já foi aplicado no banco.

No nosso projeto, o Flyway é executado via **Docker**, garantindo consistência entre ambientes e dispensando a necessidade de instalação local.

---

## 🔧 Como usar

1. Dê permissão de execução:
    ```shell
    chmod +x flyway.sh
    ```
1. Rode:
    ```shell
    ./flyway.sh           # padrão: migrate
    ./flyway.sh info      # mostra status das migrações
    ./flyway.sh validate  # valida integridade
    ./flyway.sh repair    # corrige metadados

    ```

## ⚙️ Estrutura de Diretórios

```text
project/
├── db/
│   └── migration/
│       ├── V1__create_tables.sql
│       ├── V2__add_indexes.sql
│       └── V3__insert_seed_data.sql
├── flyway.sh
└── README_FLYWAY.md
````

- db/migration/ → onde ficam os scripts SQL versionados.
- flyway.sh → script utilitário para executar o Flyway via Docker.

## 🚀 Como funciona o processo

1. Você cria scripts SQL versionados na pasta db/migration/.
1. O Flyway executa apenas os scripts novos ainda não aplicados no banco.
1. Ele mantém um controle interno (tabela flyway_schema_history) com o histórico de execuções.
1. Cada ambiente aplica as migrações na mesma ordem, garantindo consistência.

## 🧩 Nomeação dos arquivos de migração

O padrão de nomes é obrigatório e segue esta convenção:

```
V<versão>__<descrição>.sql
```

Exemplos:

```
V1__create_users_table.sql
V2__add_email_to_users.sql
V3__insert_default_roles.sql
```

Regras:

- O prefixo V indica uma versão sequencial.
- Duplo sublinhado __ separa o número da descrição.
- Use snake_case na descrição.
- Nunca altere um arquivo que já foi aplicado — sempre crie uma nova versão (V4__...).

## 🧠 Scripts Repetíveis

Além dos scripts versionados (V1__, V2__), o Flyway suporta scripts repetíveis com o prefixo R__.

Eles são reexecutados sempre que seu conteúdo muda, úteis para:

- Views
- Stored Procedures
- Funções SQL

Exemplo:

```
R__refresh_materialized_views.sql
```

## 🧰 Como rodar o Flyway

### 🐚 Via script flyway.sh

O script flyway.sh encapsula o comando Docker e permite rodar migrações facilmente.

Uso:

```shell
./flyway.sh [comando]
```

Comando padrão: migrate

Exemplos:

```shell
./flyway.sh           # Executa 'migrate'
./flyway.sh info      # Mostra status das migrações
./flyway.sh validate  # Valida integridade
./flyway.sh repair    # Corrige metadados do histórico
./flyway.sh clean     # ⚠️ Limpa o banco (use com cuidado)
```

## ⚙️ O que o script faz internamente

O flyway.sh:

1. Resolve o IP do host Windows automaticamente (para o caso de WSL2).
1. Monta o diretório db/migration dentro do container Docker.
1. Passa a URL de conexão, usuário e senha para o Flyway.
1. Executa o comando especificado (migrate, info, etc.) dentro do container.

Exemplo simplificado do comando gerado:

```shell
docker run --rm \
  -v $(pwd)/db/migration:/flyway/sql \
  flyway/flyway:11.15.0 \
  -locations=filesystem:/flyway/sql \
  -url=jdbc:postgresql://<host_ip>:5432/<banco> \
  -user=<usuario> \
  -password=<senha> \
  migrate
```

## ✅ Boas práticas

1. Scripts imutáveis: nunca altere um arquivo de migração já aplicado.
1. Crie novas versões: para cada alteração de schema, crie Vx__nova_acao.sql.
1. Use transações: sempre que possível, para garantir rollback em caso de erro.
1. Valide antes de aplicar:
    ```shell
    ./flyway.sh validate
    ```
1. Evite clean em produção.
1. Inclua o diretório db/migration no versionamento Git.

## 🔍 Comandos principais do Flyway

| Comando  | Descrição                                                       |  
|----------|-----------------------------------------------------------------|
| migrate  | Aplica novas migrações pendentes                                |
| info     | Lista todas as migrações e status                               |
| validate | Verifica integridade entre scripts e histórico                  |
| repair   | Corrige inconsistências no histórico                            |
| clean    | ⚠️ Remove todos os objetos do schema (use apenas em dev/teste)  |

## 🧩 Exemplo de fluxo de trabalho

```shell
# 1. Criar nova migração
touch db/migration/V4__add_products_table.sql

# 2. Editar o script SQL
vim db/migration/V4__add_products_table.sql

# 3. Rodar migração
./flyway.sh migrate

# 4. Ver status
./flyway.sh info

```

## 🧾 Referências

- [Documentação oficial do Flyway](https://documentation.red-gate.com/fd)

- [Imagem oficial Docker Flyway](https://hub.docker.com/r/flyway/flyway)

## 💬 Dica do engenheiro:
Trate as migrações de banco como código.
Versione, revise e teste antes de aplicar — isso mantém a base de dados tão confiável quanto o seu código-fonte.