# Defender Vehicle - Gamified Anomaly Detection Prototype

Este projeto é um protótipo funcional construído em um dia para simular a detecção de anomalias em dados de telemetria de veículos autônomos, usando uma abordagem gamificada.

## Features

- **Gerador de Telemetria:** Backend em Go que gera dados de veículos (posição, velocidade) em tempo real.
- **Injeção de Anomalias:** Simula 3 tipos de ataques: `gps_spoof`, `speed_spike`, e `sensor_dropout`.
- **API RESTful:** API em Go (usando Gin) para expor os dados e receber interações do usuário.
- **Dashboard Interativo:** Frontend em HTML/JS/Chart.js que exibe os dados em gráficos e permite que o usuário identifique anomalias.
- **Gamificação:** Sistema de pontuação e badges para engajar o usuário na tarefa de detecção.
- **Exportação de Logs:** Funcionalidade para exportar as ações do usuário em formato CSV para análise.
- **Simulação de Mitigação:** Sugere ações de resposta automatizadas para anomalias confirmadas.

## Tech Stack

- **Backend:** Go 1.23+, Gin Web Framework
- **Frontend:** HTML5, CSS3, JavaScript (Vanilla), Chart.js
- **Arquitetura:** Camadas (Handlers, Services, Models)

## Prerequisites

- Go 1.23 ou superior instalado.

## Local Setup & Execution

1.  **Clone o repositório:**
    ```bash
    git clone [https://github.com/lucasdamasceno96/defender-vehicle.git](https://github.com/lucasdamasceno96/defender-vehicle.git)
    cd defender-vehicle
    ```

2.  **Instale as dependências:**
    ```bash
    go mod tidy
    ```

3.  **Execute o servidor:**
    ```bash
    go run ./cmd/server/
    ```

4.  **Acesse a aplicação:**
    Abra seu navegador e acesse `http://localhost:8080`.

## API Endpoints

| Método | Path                  | Descrição                                         |
| :----- | :-------------------- | :------------------------------------------------ |
| `GET`  | `/api/health`         | Verifica o status da aplicação.                   |
| `GET`  | `/api/telemetry`      | Retorna dados de telemetria (suporta `?limit=`).  |
| `GET`  | `/api/gamestate`      | Retorna a pontuação atual e os badges do jogador. |
| `POST` | `/api/detect`         | Recebe a detecção de anomalia do usuário.         |
| `POST` | `/api/mitigate/:id`   | Simula a ativação de mitigação para um ponto.     |
| `GET`  | `/api/logs`           | Faz o download dos logs de detecção em CSV.       |