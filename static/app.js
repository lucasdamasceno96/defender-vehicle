// static/app.js

document.addEventListener('DOMContentLoaded', () => {
    // --- Variáveis Globais e Elementos do DOM ---
    let gpsChart, speedChart;
    let telemetryData = [];

    // Seleciona todos os elementos da página que vamos manipular
    const scoreElement = document.getElementById('playerScore');
    const actionStatusElement = document.getElementById('actionStatus');
    const mitigationPanel = document.getElementById('mitigationPanel');
    const threatTypeElement = document.getElementById('threatType');
    const mitigationListElement = document.getElementById('mitigationList');
    const badgeContainer = document.getElementById('badge-container');

    // --- Configurações de Cores para os Gráficos ---
    const colors = {
        normal: 'rgba(54, 162, 235, 0.6)',       // Azul para pontos normais
        anomaly: 'rgba(255, 206, 86, 0.8)',      // Amarelo para anomalias não clicadas
        correctGuess: 'rgba(75, 192, 192, 1)',   // Verde para acertos
        incorrectGuess: 'rgba(255, 99, 132, 1)', // Vermelho para erros
    };

    /**
     * Função principal que inicializa o dashboard.
     * Busca os dados da telemetria, desenha os gráficos e atualiza o estado do jogo.
     */
    const initializeDashboard = async () => {
        try {
            // 1. Busca os dados iniciais da telemetria da nossa API Go.
            const response = await fetch('/api/telemetry?limit=100');
            if (!response.ok) throw new Error('Falha ao buscar dados de telemetria');
            telemetryData = await response.json();

            // 2. Com os dados em mãos, desenha os dois gráficos na tela.
            renderGpsChart();
            renderSpeedChart();

            // 3. Busca o estado inicial do jogo (placar e badges) e atualiza a tela.
            await updateGameState();
        } catch (error) {
            console.error('Falha ao inicializar o dashboard:', error);
            actionStatusElement.textContent = 'Erro ao carregar os dados.';
        }
    };

    /**
     * Busca o placar e os badges mais recentes do backend e atualiza a interface.
     */
    const updateGameState = async () => {
        try {
            const response = await fetch('/api/gamestate');
            if (!response.ok) throw new Error('Falha ao buscar o estado do jogo');
            const state = await response.json();

            // Atualiza o texto do placar
            scoreElement.textContent = state.player_score;

            // Limpa os badges antigos e renderiza os novos
            badgeContainer.innerHTML = '';
            if (state.player_badges && state.player_badges.length > 0) {
                state.player_badges.forEach(badgeText => {
                    const badge = document.createElement('span');
                    badge.textContent = badgeText;
                    badge.className = 'badge'; // Usando classe para o estilo
                    badgeContainer.appendChild(badge);
                });
            } else {
                badgeContainer.innerHTML = '<p id="no-badges">Nenhum ainda.</p>';
            }
        } catch (error) {
            console.error('Falha ao atualizar o estado do jogo:', error);
        }
    };

    /**
     * Renderiza o gráfico de dispersão do GPS (Longitude vs. Latitude).
     */
    const renderGpsChart = () => {
        const ctx = document.getElementById('gpsChart').getContext('2d');
        gpsChart = new Chart(ctx, {
            type: 'scatter',
            data: {
                datasets: [{
                    label: 'Rota do GPS',
                    data: telemetryData.map(p => ({ x: p.lon, y: p.lat })),
                    backgroundColor: telemetryData.map(p => p.is_anomaly ? colors.anomaly : colors.normal),
                    pointRadius: 5,
                    pointHoverRadius: 8,
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { display: false },
                    title: { display: true, text: 'Rota do Veículo (Longitude vs. Latitude)' }
                },
                // A ação mais importante: define o que acontece quando o gráfico é clicado.
                onClick: handleChartClick,
            }
        });
    };

    /**
     * Renderiza o gráfico de linha da Velocidade ao longo do tempo.
     */
    const renderSpeedChart = () => {
        const ctx = document.getElementById('speedChart').getContext('2d');
        speedChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: telemetryData.map(p => new Date(p.ts).toLocaleTimeString()),
                datasets: [{
                    label: 'Velocidade (km/h)',
                    data: telemetryData.map(p => p.speed),
                    borderColor: 'rgba(75, 192, 192, 1)',
                    tension: 0.1,
                    fill: false,
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    title: { display: true, text: 'Velocidade vs. Tempo' }
                }
            }
        });
    };

    /**
     * Manipulador de eventos para cliques no gráfico de GPS.
     * Envia a detecção para a API e atualiza a UI com o resultado.
     */
    const handleChartClick = async (event) => {
        // Identifica o ponto mais próximo do clique do mouse.
        const points = gpsChart.getElementsAtEventForMode(event, 'nearest', { intersect: true }, true);
        if (points.length === 0) return; // Se não clicou em um ponto, não faz nada.

        const pointIndex = points[0].index;
        const telemetryPoint = telemetryData[pointIndex];

        // Atualiza a UI para mostrar que estamos processando a ação.
        actionStatusElement.textContent = `Analisando ponto ID ${telemetryPoint.id}...`;
        mitigationPanel.style.display = 'none'; // Esconde o painel de mitigação anterior.

        try {
            // Envia o ID do ponto clicado para a API de detecção.
            const response = await fetch('/api/detect', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ telemetry_id: telemetryPoint.id }),
            });
            const result = await response.json();

            // Com base na resposta da API, atualiza a UI.
            if (result.result === true) {
                actionStatusElement.textContent = `Correto! ID ${telemetryPoint.id} é uma anomalia.`;
                updatePointColor(pointIndex, colors.correctGuess);
                await showMitigations(telemetryPoint.id, telemetryPoint.anomaly_type);
            } else {
                actionStatusElement.textContent = `Incorreto. ID ${telemetryPoint.id} é normal.`;
                updatePointColor(pointIndex, colors.incorrectGuess);
            }

            // Após a ação, busca o novo placar e badges do backend.
            await updateGameState();

        } catch (error) {
            console.error('Erro durante a detecção:', error);
            actionStatusElement.textContent = 'Falha de comunicação com o servidor.';
        }
    };

    /**
     * Atualiza a cor de um ponto específico no gráfico de GPS após ser clicado.
     */
    const updatePointColor = (index, color) => {
        if (gpsChart && gpsChart.data.datasets[0].backgroundColor[index]) {
            gpsChart.data.datasets[0].backgroundColor[index] = color;
            gpsChart.update('none'); // 'none' para uma atualização sem animação, mais rápida.
        }
    };

    /**
     * Se uma detecção for correta, busca e exibe as ações de mitigação sugeridas.
     */
    const showMitigations = async (telemetryId, threatType) => {
        try {
            const response = await fetch(`/api/mitigate/${telemetryId}`, { method: 'POST' });
            const mitigationData = await response.json();

            if (mitigationData.suggested_actions) {
                threatTypeElement.textContent = threatType || 'N/A';
                mitigationListElement.innerHTML = ''; // Limpa a lista antiga.
                mitigationData.suggested_actions.forEach(action => {
                    const li = document.createElement('li');
                    li.textContent = action;
                    mitigationListElement.appendChild(li);
                });
                mitigationPanel.style.display = 'block'; // Mostra o painel.
            }
        } catch (error) {
            console.error('Falha ao buscar mitigações:', error);
        }
    };

    // --- Ponto de Entrada da Aplicação ---
    // Inicia todo o processo assim que a página termina de carregar.
    initializeDashboard();
});