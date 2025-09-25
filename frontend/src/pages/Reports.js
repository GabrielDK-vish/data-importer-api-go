import React, { useState, useEffect } from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import api from '../services/api';

function Reports() {
  const [monthlyBilling, setMonthlyBilling] = useState([]);
  const [billingByProduct, setBillingByProduct] = useState([]);
  const [billingByPartner, setBillingByPartner] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadReports();
  }, []);

  const loadReports = async () => {
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
      setError('Erro ao carregar relatórios: ' + (err.response?.data?.message || err.message));
      console.error('Erro:', err);
    } finally {
      setLoading(false);
    }
  };

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D'];

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'USD'
    }).format(value);
  };

  if (loading) {
    return <div className="loading">Carregando relatórios...</div>;
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error">
          <h2>⚠️ {error}</h2>
          <p>Para começar a usar os relatórios:</p>
          <ol>
            <li>Vá para a página de <strong>Upload</strong></li>
            <li>Faça upload de um arquivo Excel ou CSV</li>
            <li>Volte aqui para visualizar os relatórios</li>
          </ol>
          <button onClick={() => window.location.href = '/upload'} className="btn-primary">
            Ir para Upload
          </button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1>📊 Relatórios</h1>
      
      {/* Faturamento mensal */}
      <div className="chart-container">
        <h2>📈 Faturamento Mensal</h2>
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={monthlyBilling}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="month" />
            <YAxis />
            <Tooltip 
              formatter={(value) => [formatCurrency(value), 'Total']}
              labelFormatter={(label) => `Mês: ${label}`}
            />
            <Bar dataKey="total" fill="#8884d8" />
          </BarChart>
        </ResponsiveContainer>
        
        <div style={{ marginTop: '20px' }}>
          <h3>Resumo Mensal</h3>
          <table className="table">
            <thead>
              <tr>
                <th>Mês</th>
                <th>Total</th>
                <th>Registros</th>
                <th>Média por Registro</th>
              </tr>
            </thead>
            <tbody>
              {monthlyBilling.map((item, index) => (
                <tr key={index}>
                  <td>{item.month}</td>
                  <td>{formatCurrency(item.total)}</td>
                  <td>{item.count}</td>
                  <td>{formatCurrency(item.total / item.count)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Faturamento por produto */}
      <div className="chart-container">
        <h2>🛍️ Faturamento por Produto</h2>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
          <div>
            <h3>Gráfico de Pizza</h3>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={billingByProduct.slice(0, 6)}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} (${(percent * 100).toFixed(0)}%)`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="total"
                >
                  {billingByProduct.slice(0, 6).map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [formatCurrency(value), 'Total']} />
              </PieChart>
            </ResponsiveContainer>
          </div>
          
          <div>
            <h3>Gráfico de Barras</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={billingByProduct.slice(0, 8)} layout="horizontal">
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis type="number" />
                <YAxis dataKey="product_name" type="category" width={100} />
                <Tooltip formatter={(value) => [formatCurrency(value), 'Total']} />
                <Bar dataKey="total" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
        
        <div style={{ marginTop: '20px' }}>
          <h3>Detalhes por Produto</h3>
          <table className="table">
            <thead>
              <tr>
                <th>Produto</th>
                <th>Categoria</th>
                <th>Total</th>
                <th>Registros</th>
                <th>Média por Registro</th>
              </tr>
            </thead>
            <tbody>
              {billingByProduct.map((item, index) => (
                <tr key={index}>
                  <td>{item.product_name}</td>
                  <td>{item.category}</td>
                  <td>{formatCurrency(item.total)}</td>
                  <td>{item.count}</td>
                  <td>{formatCurrency(item.total / item.count)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Faturamento por parceiro */}
      <div className="chart-container">
        <h2>🤝 Faturamento por Parceiro</h2>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
          <div>
            <h3>Gráfico de Pizza</h3>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={billingByPartner}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} (${(percent * 100).toFixed(0)}%)`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="total"
                >
                  {billingByPartner.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [formatCurrency(value), 'Total']} />
              </PieChart>
            </ResponsiveContainer>
          </div>
          
          <div>
            <h3>Gráfico de Barras</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={billingByPartner}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="partner_name" />
                <YAxis />
                <Tooltip formatter={(value) => [formatCurrency(value), 'Total']} />
                <Bar dataKey="total" fill="#8884d8" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
        
        <div style={{ marginTop: '20px' }}>
          <h3>Detalhes por Parceiro</h3>
          <table className="table">
            <thead>
              <tr>
                <th>Parceiro</th>
                <th>Total</th>
                <th>Registros</th>
                <th>Média por Registro</th>
              </tr>
            </thead>
            <tbody>
              {billingByPartner.map((item, index) => (
                <tr key={index}>
                  <td>{item.partner_name}</td>
                  <td>{formatCurrency(item.total)}</td>
                  <td>{item.count}</td>
                  <td>{formatCurrency(item.total / item.count)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default Reports;
