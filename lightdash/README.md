# 🧭 Tutorial de Instalação — Lightdash + dbt com ASDF

Este guia ensina a configurar um ambiente com **Node.js**, **Python**, **Lightdash CLI** e **dbt** utilizando o gerenciador de versões **asdf**.

---

## 📦 1. Instalar dependências do sistema

```bash
sudo apt update
sudo apt install -y \
  build-essential make wget curl ca-certificates \
  zlib1g-dev libbz2-dev libreadline-dev libsqlite3-dev \
  libffi-dev liblzma-dev uuid-dev tk-dev \
  libgdbm-dev libgdbm-compat-dev \
  libncurses5-dev libncursesw5-dev \
  xz-utils libssl-dev git
```

Esses pacotes são necessários para compilar o Python e dependências nativas.

---

## 🟢 2. Instalar Node.js via ASDF

```bash
asdf plugin add nodejs
asdf install nodejs 20.8.0
asdf local nodejs 20.8.0
```

Verifique a versão:

```bash
node -v
npm -v
```

---

## ⚡ 3. Instalar o Lightdash CLI

```bash
npm install -g @lightdash/cli@0.1826.3
```

Verifique a instalação:

```bash
lightdash --version
```

---

## 🔐 4. Fazer login no Lightdash

Com **token**:

```bash
lightdash login https://test-lightdash.yec7w0.easypanel.host --token my-super-secret-token
```

Ou com **OAuth**:

```bash
lightdash login https://test-lightdash.yec7w0.easypanel.host --oauth
```

Após o login, você pode testar o acesso:

```bash
lightdash status
```

---

## 🐍 5. Instalar Python via ASDF

```bash
asdf plugin add python
asdf install python 3.12.12
asdf local python 3.12.12
```

Verifique:

```bash
python --version
```

---

## 📈 6. Instalar dbt-core e dbt-postgres

```bash
python -m pip install --upgrade pip
pip install dbt-core dbt-postgres
```

Verifique a instalação:

```bash
dbt --version
```

Saída esperada (exemplo):

```
Core:
  - installed: 1.10.13
  - latest:    1.10.13 - Up to date!

Plugins:
  - postgres: 1.9.1 - Up to date!
```

---

## ✅ 8. Verificação Final

Execute os seguintes comandos para confirmar o ambiente:

```bash
node -v
python -V
dbt --version
lightdash --version
```

Tudo pronto para começar a usar o **dbt** e conectar ao **Lightdash** 🎯
