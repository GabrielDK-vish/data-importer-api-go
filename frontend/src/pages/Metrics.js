import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Metrics() {
  const [metrics, setMetrics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadMetricsData();
  }, []);

  const loadMetricsData = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response = await api.get('/api/metrics/processing');
      setMetrics(response.data || []);
    } catch (err) {
      console.error('Erro ao carregar métricas de processamento:', err);
      setError('Falha ao carregar métricas de processamento. Tente novamente mais tarde.');
    } finally {
      setLoading(false);
    }
  };

  // Formatar data e hora
  const formatDateTime = (timestamp) => {
    if (!timestamp) return 'N/A';
    const date = new Date(timestamp);
    return date.toLocaleString('pt-BR');
  };

  // Formatar duração em milissegundos para formato legível
  const formatDuration = (ms) => {
    if (!ms) return 'N/A';
    
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    
    if (minutes > 0) {
      return `${minutes}m ${remainingSeconds}s`;
    }
    return `${seconds}.${ms % 1000}s`;
  };

  return (
    <div className="metrics-container">
      <h2>Métricas de Processamento</h2>
      
      {loading ? (
        <div className="loading">Carregando métricas...</div>
      ) : error ? (
        <div className="error-message">{error}</div>
      ) : metrics.length === 0 ? (
        <div className="no-data">Nenhuma métrica de processamento disponível.</div>
      ) : (
        <div className="table-container">
          <table className="data-table">
            <thead>
              <tr>
                <th>Arquivo</th>
                <th>Data de Processamento</th>
                <th>Tempo de Processamento</th>
                <th>Registros Processados</th>
              </tr>
            </thead>
            <tbody>
              {metrics.map((metric) => (
                <tr key={metric.id}>
                  <td>{metric.file_name || 'Desconhecido'}</td>
                  <td>{formatDateTime(metric.created_at)}</td>
                  <td>{formatDuration(metric.duration_ms)}</td>
                  <td>{metric.records_count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default Metrics;