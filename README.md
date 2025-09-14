# Defender Vehicle - Protótipo de Detecção Gamificada de Anomalias

Este projeto é um protótipo funcional para simular a detecção de anomalias em dados de veículos autônomos, usando uma abordagem interativa e gamificada.

## Entendendo o Protótipo (Para Leigos)

### O Cenário: O Desafio da Cibersegurança Veicular

Pense em um carro autônomo como um computador sobre rodas. Ele depende de dezenas de sensores (GPS, velocidade, câmeras) para "sentir" o mundo ao seu redor e tomar decisões seguras. Mas o que acontece se um hacker consegue enganar esses sentidos? O resultado pode ser perigoso. Este protótipo simula exatamente esse cenário: você é o analista de segurança responsável por proteger o veículo.

### Sua Missão como Analista

Sua missão é atuar como um "firewall humano". O dashboard à sua frente mostra os dados vitais do veículo em tempo real. Seu objetivo é usar sua intuição e atenção para identificar comportamentos estranhos nos dados que possam indicar um ataque em andamento.

### Os Ataques Simulados: Entendendo as Ameaças

Nesta simulação, você enfrentará três tipos de ataques:

1.  **📍 `GPS Spoofing` (GPS Falsificado):**
    -   **O que é:** O hacker envia um sinal de GPS falso para o carro, fazendo-o pensar que está em outro lugar.
    -   **Como aparece no gráfico:** Você verá um "salto" repentino e ilógico na rota do veículo no mapa. Um carro não pode se teletransportar.

2.  **⚡ `Speed Spike` (Pico de Velocidade Falso):**
    -   **O que é:** O sensor de velocidade é manipulado para reportar um valor absurdamente alto por um instante.
    -   **O que causa:** O carro pode pensar que está a 150 km/h quando na verdade está a 60 km/h, fazendo-o frear bruscamente ou tomar outras decisões perigosas.

3.  **⚫ `Sensor Dropout` (Apagão do Sensor):**
    -   **O que é:** O hacker consegue "desligar" um sensor. O carro para de receber informações vitais.
    -   **Como aparece no gráfico:** Os dados do sensor (velocidade, posição) caem para zero subitamente. É como se o carro ficasse "cego" e "surdo".

### Como Defender: A Gamificação e a Pontuação

Para tornar a tarefa mais interativa, usamos a gamificação.

-   **Como Jogar:** É simples. Observe os gráficos. Quando vir um ponto de dados que parece suspeito no gráfico de GPS, **clique sobre ele**.
-   **A Pontuação (Seu Desempenho):**
    -   `Detecção Correta (+10 pontos)`: Você clicou em um ataque real. Excelente! Você protegeu o veículo.
    -   `Alarme Falso (-5 pontos)`: Você clicou em um ponto normal. Em um cenário real, alarmes falsos custam tempo e recursos, por isso há uma pequena penalidade.
-   **Badges (Conquistas):** Conforme você identifica ameaças, você ganha "medalhas" por seus feitos, como "Primeira Captura" ou "Caçador de Ameaças".

### As Respostas (Mitigações): O Que o Carro Faria?

Após você identificar uma ameaça real, o sistema mostra um "Plano de Mitigação". Isso simula as ações que o sistema de defesa do carro tomaria automaticamente. Por exemplo:

-   Se o GPS for atacado, o carro pode decidir **ignorar o GPS temporariamente e confiar em outros sensores** para se manter na rota.
-   Em caso de falha crítica, ele pode entrar em **"modo de segurança"**, reduzindo a velocidade e alertando o motorista para assumir o controle.

## Features

- **Gerador de Telemetria:** Backend em Go que gera dados simulados de veículos.
- **Injeção de Anomalias:** Simula os 3 ataques descritos acima.
- **API RESTful:** API em Go (usando Gin) para expor os dados e receber interações.
- **Dashboard Interativo:** Frontend em HTML/JS/Chart.js para visualização e interação.
- **Gamificação:** Sistema de pontuação e badges para engajar o analista.
- **Exportação de Logs:** Funcionalidade para exportar as ações do usuário em formato CSV.

## Tech Stack

- **Backend:** Go 1.23+, Gin Web Framework
- **Frontend:** HTML5, CSS3, JavaScript (Vanilla), Chart.js
- **Containerização:** Docker

## Prerequisites

- Go 1.23 ou superior instalado.
- Docker instalado.

## Como Executar

### Opção 1: Localmente (sem Docker)

1.  **Clone o repositório:** `git clone https://github.com/lucasdamasceno96/defender-vehicle.git`
2.  **Entre na pasta:** `cd defender-vehicle`
3.  **Instale as dependências:** `go mod tidy`
4.  **Execute o servidor:** `go run ./cmd/server/`
5.  **Acesse:** `http://localhost:8080`

### Opção 2: Usando Docker

1.  **Clone o repositório e entre na pasta.**
2.  **Construa a imagem:** `docker build -t defender-vehicle:latest .`
3.  **Execute o container:** `docker run --rm -p 8080:8080 defender-vehicle:latest`
4.  **Acesse:** `http://localhost:8080`