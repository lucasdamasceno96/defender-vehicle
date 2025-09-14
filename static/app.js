// static/app.js

document.addEventListener('DOMContentLoaded', () => {
    // --- Chart instances and data store ---
    let gpsChart, speedChart;
    let telemetryData = [];
    const scoreElement = document.getElementById('playerScore');
    const actionStatusElement = document.getElementById('actionStatus');
    const mitigationPanel = document.getElementById('mitigationPanel');
    const threatTypeElement = document.getElementById('threatType');
    const mitigationListElement = document.getElementById('mitigationList');

    // --- Chart Colors ---
    const colors = {
        normal: 'rgba(54, 162, 235, 0.6)',
        anomaly: 'rgba(255, 206, 86, 0.8)',
        correctGuess: 'rgba(75, 192, 192, 1)',
        incorrectGuess: 'rgba(255, 99, 132, 1)',
    };

    // --- Fetch initial data and render charts ---
    const initializeDashboard = async () => {
        try {
            // Fetch all 100 points for the client-side simulation
            const response = await fetch('/api/telemetry?limit=100');
            telemetryData = await response.json();
            renderGpsChart();
            renderSpeedChart();
        } catch (error) {
            console.error('Failed to load telemetry data:', error);
            actionStatusElement.textContent = 'Error loading data.';
        }
    };

    // --- Chart Rendering Functions ---
    const renderGpsChart = () => {
        const ctx = document.getElementById('gpsChart').getContext('2d');
        const data = {
            datasets: [{
                label: 'GPS Path',
                data: telemetryData.map(p => ({ x: p.lon, y: p.lat })),
                backgroundColor: telemetryData.map(p => p.is_anomaly ? colors.anomaly : colors.normal),
                pointRadius: 5,
                pointHoverRadius: 8,
            }]
        };
        gpsChart = new Chart(ctx, {
            type: 'scatter',
            data: data,
            options: {
                responsive: true,
                plugins: {
                    legend: { display: false },
                    title: { display: true, text: 'GPS Path (Longitude vs. Latitude)' }
                },
                onClick: handleChartClick,
            }
        });
    };

    const renderSpeedChart = () => {
        const ctx = document.getElementById('speedChart').getContext('2d');
        const labels = telemetryData.map(p => new Date(p.ts).toLocaleTimeString());
        const data = {
            labels: labels,
            datasets: [{
                label: 'Speed (km/h)',
                data: telemetryData.map(p => p.speed),
                borderColor: 'rgba(75, 192, 192, 1)',
                tension: 0.1,
                fill: false,
            }]
        };
        speedChart = new Chart(ctx, {
            type: 'line',
            data: data,
            options: {
                responsive: true,
                plugins: {
                    title: { display: true, text: 'Speed Over Time' }
                }
            }
        });
    };

    // --- Interactivity Handlers ---
    const handleChartClick = async (event) => {
        const points = gpsChart.getElementsAtEventForMode(event, 'nearest', { intersect: true }, true);
        if (points.length === 0) return;

        const pointIndex = points[0].index;
        const telemetryPoint = telemetryData[pointIndex];

        actionStatusElement.textContent = `Analyzing point ID ${telemetryPoint.id}...`;
        mitigationPanel.style.display = 'none'; // Hide previous mitigation

        try {
            const response = await fetch('/api/detect', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ telemetry_id: telemetryPoint.id }),
            });
            const result = await response.json();

            // Update UI based on detection result
            if (result.result === true) {
                actionStatusElement.textContent = `Correct! ID ${telemetryPoint.id} is an anomaly. +10 points.`;
                updatePointColor(pointIndex, colors.correctGuess);
                await showMitigations(telemetryPoint.id, telemetryPoint.anomaly_type);
            } else {
                actionStatusElement.textContent = `Incorrect. ID ${telemetryPoint.id} is normal. -5 points.`;
                updatePointColor(pointIndex, colors.incorrectGuess);
            }
            // A real app would get the score from the backend, but we'll simulate for now.
            let currentScore = parseInt(scoreElement.textContent);
            scoreElement.textContent = result.result ? currentScore + 10 : currentScore - 5;
        } catch (error) {
            console.error('Error during detection:', error);
            actionStatusElement.textContent = 'Failed to communicate with the server.';
        }
    };

    const updatePointColor = (index, color) => {
        gpsChart.data.datasets[0].backgroundColor[index] = color;
        gpsChart.update();
    };

    const showMitigations = async (telemetryId, threatType) => {
        try {
            const response = await fetch(`/api/mitigate/${telemetryId}`, { method: 'POST' });
            const mitigationData = await response.json();

            if (mitigationData.suggested_actions) {
                threatTypeElement.textContent = threatType || 'N/A';
                mitigationListElement.innerHTML = ''; // Clear previous list
                mitigationData.suggested_actions.forEach(action => {
                    const li = document.createElement('li');
                    li.textContent = action;
                    mitigationListElement.appendChild(li);
                });
                mitigationPanel.style.display = 'block';
            }
        } catch (error) {
            console.error('Failed to fetch mitigations:', error);
        }
    };


    // --- Start the application ---
    initializeDashboard();
});