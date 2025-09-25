import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import api from '../services/api';

function Dashboard() {
  const [monthlyBilling, setMonthlyBilling] = useState([]);
  const [billingByProduct, setBillingByProduct] = useState([]);
  const [billingByPartner, setBillingByPartner] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      setError('');
      
      const [monthlyRes, productRes, partnerRes] = await Promise.all([
        api.get('/api/reports/billing/monthly').catch(err => {
          console.warn('Erro ao carregar dados mensais:', err);
          return { data: [] };
        }),
        api.get('/api/reports/billing/by-product').catch(err => {
          console.warn('Erro ao carregar dados por produto:', err);
          return { data: [] };
        }),
        api.get('/api/reports/billing/by-partner').catch(err => {
          console.warn('Erro ao carregar dados por parceiro:', err);
          return { data: [] };
        })
      ]);

      setMonthlyBilling(monthlyRes.data || []);
      setBillingByProduct(productRes.data || []);
      setBillingByPartner(partnerRes.data || []);
      
      // Se não há dados, mostrar mensagem informativa
      if (!monthlyRes.data?.length && !productRes.data?.length && !partnerRes.data?.length) {
        setError('Nenhum dado encontrado. Faça upload de um arquivo para visualizar os relatórios.');
      }
    } catch (err) {
      setError('Erro ao carregar dados do dashboard: ' + (err.response?.data?.message || err.message));
      console.error('Erro:', err);
    } finally {
      setLoading(false);
    }
  };

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

  if (loading) {
    return <div className="loading">Carregando dashboard...</div>;
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error">
          <h2>⚠️ {error}</h2>
          <p>Para começar a usar o dashboard:</p>
          <ol>
            <li>Vá para a página de <strong>Upload</strong></li>
            <li>Faça upload de um arquivo Excel ou CSV</li>
            <li>Volte ao dashboard para visualizar os dados</li>
          </ol>
          <button onClick={() => window.location.href = '/upload'} className="btn-primary">
            Ir para Upload
          </button>
        </div>
      </div>
    );
  }

  const totalMonthly = monthlyBilling.reduce((sum, item) => sum + item.total, 0);
  const totalProducts = billingByProduct.reduce((sum, item) => sum + item.total, 0);
  const totalPartners = billingByPartner.reduce((sum, item) => sum + item.total, 0);

  return (
    <div>
      <h1>📊 Dashboard</h1>
      
      {/* Cards de estatísticas */}
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Faturamento Total Mensal</h3>
          <div className="value">${totalMonthly.toFixed(2)}</div>
        </div>
        <div className="stat-card">
          <h3>Total por Produtos</h3>
          <div className="value">${totalProducts.toFixed(2)}</div>
        </div>
        <div className="stat-card">
          <h3>Total por Parceiros</h3>
          <div className="value">${totalPartners.toFixed(2)}</div>
        </div>
        <div className="stat-card">
          <h3>Meses com Dados</h3>
          <div className="value">{monthlyBilling.length}</div>
        </div>
      </div>

      {/* Gráfico de faturamento mensal */}
      <div className="chart-container">
        <h2>📈 Faturamento Mensal</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={monthlyBilling}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip formatter={(value) => [`$${value.toFixed(2)}`, 'Total']} />
            <Line 
              type="monotone" 
              dataKey="total" 
              stroke="#8884d8" 
              strokeWidth={2}
              dot={{ fill: '#8884d8', strokeWidth: 2, r: 4 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Gráfico de faturamento por produto */}
      <div className="chart-container">
        <h2>🛍️ Faturamento por Produto</h2>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={billingByProduct.slice(0, 5)}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={({ name, percent }) => `${name} (${(percent * 100).toFixed(0)}%)`}
              outerRadius={80}
              fill="#8884d8"
              dataKey="total"
            >
              {billingByProduct.slice(0, 5).map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip formatter={(value) => [`$${value.toFixed(2)}`, 'Total']} />
          </PieChart>
        </ResponsiveContainer>
      </div>

      {/* Tabela de top produtos */}
      <div className="card">
        <h2>🏆 Top Produtos por Faturamento</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Produto</th>
              <th>Categoria</th>
              <th>Total</th>
              <th>Registros</th>
            </tr>
          </thead>
          <tbody>
            {billingByProduct.slice(0, 10).map((item, index) => (
              <tr key={index}>
                <td>{item.product_name}</td>
                <td>{item.category}</td>
                <td>${item.total.toFixed(2)}</td>
                <td>{item.count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Tabela de top parceiros */}
      <div className="card">
        <h2>🤝 Top Parceiros por Faturamento</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Parceiro</th>
              <th>Total</th>
              <th>Registros</th>
            </tr>
          </thead>
          <tbody>
            {billingByPartner.map((item, index) => (
              <tr key={index}>
                <td>{item.partner_name}</td>
                <td>${item.total.toFixed(2)}</td>
                <td>{item.count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default Dashboard;
