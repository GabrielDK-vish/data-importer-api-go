import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Customers() {
  const [customers, setCustomers] = useState([]);
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [customerUsage, setCustomerUsage] = useState([]);
  const [loading, setLoading] = useState(true);
  const [usageLoading, setUsageLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadCustomers();
  }, []);

  const loadCustomers = async () => {
    try {
      setLoading(true);
      setError('');
      const response = await api.get('/api/customers');
      setCustomers(response.data || []);
      
      if (!response.data?.length) {
        setError('Nenhum cliente encontrado. Faça upload de um arquivo para importar dados.');
      }
    } catch (err) {
      setError('Erro ao carregar clientes: ' + (err.response?.data?.message || err.message));
      console.error('Erro:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadCustomerUsage = async (customerId) => {
    try {
      setUsageLoading(true);
      setError(''); // Limpar erros anteriores
      const response = await api.get(`/api/customers/${customerId}/usage`);
      setCustomerUsage(response.data || []);
    } catch (err) {
      console.error('Erro ao carregar uso do cliente:', err);
      // Não definir erro global, apenas mostrar mensagem local
      setCustomerUsage([]);
    } finally {
      setUsageLoading(false);
    }
  };

  const handleCustomerClick = (customer) => {
    setSelectedCustomer(customer);
    loadCustomerUsage(customer.id);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  const formatCurrency = (value) => {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: 'USD'
    }).format(value);
  };

  if (loading) {
    return <div className="loading">Carregando clientes...</div>;
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error">
          <h2>{error}</h2>
          <p>Para começar a usar a lista de clientes:</p>
          <ol>
            <li>Vá para a página de <strong>Upload</strong></li>
            <li>Faça upload de um arquivo Excel ou CSV</li>
            <li>Volte aqui para visualizar os clientes</li>
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
      <h1>Clientes</h1>
      
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px' }}>
        {/* Lista de clientes */}
        <div className="card">
          <h2>Lista de Clientes ({customers.length})</h2>
          <div style={{ maxHeight: '500px', overflowY: 'auto' }}>
            {customers.map((customer) => (
              <div
                key={customer.id}
                onClick={() => handleCustomerClick(customer)}
                style={{
                  padding: '15px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  marginBottom: '10px',
                  cursor: 'pointer',
                  backgroundColor: selectedCustomer?.id === customer.id ? '#e3f2fd' : 'white',
                  transition: 'background-color 0.3s'
                }}
              >
                <h3 style={{ margin: '0 0 5px 0', color: '#333' }}>
                  {customer.customer_name}
                </h3>
                <p style={{ margin: '0 0 5px 0', color: '#666' }}>
                  <strong>ID:</strong> {customer.customer_id}
                </p>
                <p style={{ margin: '0 0 5px 0', color: '#666' }}>
                  <strong>Domínio:</strong> {customer.customer_domain_name}
                </p>
                <p style={{ margin: '0', color: '#666' }}>
                  <strong>País:</strong> {customer.country}
                </p>
              </div>
            ))}
          </div>
        </div>

        {/* Detalhes do cliente selecionado */}
        <div className="card">
          {selectedCustomer ? (
            <>
              <h2>Detalhes do Cliente</h2>
              <div style={{ marginBottom: '20px' }}>
                <h3>{selectedCustomer.customer_name}</h3>
                <p><strong>ID:</strong> {selectedCustomer.customer_id}</p>
                <p><strong>Domínio:</strong> {selectedCustomer.customer_domain_name}</p>
                <p><strong>País:</strong> {selectedCustomer.country}</p>
                <p><strong>Criado em:</strong> {formatDate(selectedCustomer.created_at)}</p>
              </div>

              <h3>Histórico de Uso</h3>
              {usageLoading ? (
                <div className="loading">Carregando uso...</div>
              ) : customerUsage.length > 0 ? (
                <div style={{ maxHeight: '400px', overflowY: 'auto' }}>
                  <table className="table">
                    <thead>
                      <tr>
                        <th>Data</th>
                        <th>Produto</th>
                        <th>Parceiro</th>
                        <th>Quantidade</th>
                        <th>Preço Unit.</th>
                        <th>Total</th>
                      </tr>
                    </thead>
                    <tbody>
                      {customerUsage.map((usage, index) => (
                        <tr key={index}>
                          <td>{formatDate(usage.usage_date)}</td>
                          <td>{usage.product?.product_name || 'N/A'}</td>
                          <td>{usage.partner?.partner_name || 'N/A'}</td>
                          <td>{usage.quantity}</td>
                          <td>{formatCurrency(usage.unit_price)}</td>
                          <td>{formatCurrency(usage.billing_pre_tax_total)}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div style={{ textAlign: 'center', padding: '20px', color: '#666' }}>
                  <p>Nenhum registro de uso encontrado para este cliente.</p>
                  <p style={{ fontSize: '14px', marginTop: '10px' }}>
                    Isso pode indicar que o cliente não possui dados de uso ou que ainda não foi feito upload de dados.
                  </p>
                </div>
              )}
            </>
          ) : (
            <div style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
              <h3>Selecione um cliente</h3>
              <p>Clique em um cliente da lista para ver os detalhes e histórico de uso.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Customers;
