import React, { useState } from 'react';
import api from '../services/api';

function Upload() {
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile) {
      // Verificar tipo de arquivo
      const fileName = selectedFile.name.toLowerCase();
      if (fileName.endsWith('.csv') || fileName.endsWith('.xlsx')) {
        setFile(selectedFile);
        setError('');
      } else {
        setError('Por favor, selecione um arquivo CSV ou Excel (.xlsx)');
        setFile(null);
      }
    }
  };

  const handleUpload = async (e) => {
    e.preventDefault();
    
    if (!file) {
      setError('Por favor, selecione um arquivo');
      return;
    }

    setUploading(true);
    setError('');
    setResult(null);

    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await api.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      setResult(response.data);
      setFile(null);
      
      // Limpar input de arquivo
      const fileInput = document.getElementById('file-input');
      if (fileInput) {
        fileInput.value = '';
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Erro ao fazer upload do arquivo');
      console.error('Erro no upload:', err);
    } finally {
      setUploading(false);
    }
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div>
      <h1>📤 Upload de Arquivos</h1>
      
      <div className="card">
        <h2>Importar Dados</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          Faça upload de arquivos CSV ou Excel (.xlsx) para importar dados para o sistema.
        </p>

        <form onSubmit={handleUpload}>
          <div className="form-group">
            <label htmlFor="file-input">Selecionar Arquivo:</label>
            <input
              id="file-input"
              type="file"
              accept=".csv,.xlsx"
              onChange={handleFileChange}
              style={{
                width: '100%',
                padding: '8px 12px',
                border: '1px solid #ddd',
                borderRadius: '4px',
                fontSize: '14px'
              }}
            />
            {file && (
              <div style={{ marginTop: '10px', padding: '10px', backgroundColor: '#f8f9fa', borderRadius: '4px' }}>
                <strong>Arquivo selecionado:</strong> {file.name}<br />
                <strong>Tamanho:</strong> {formatFileSize(file.size)}<br />
                <strong>Tipo:</strong> {file.type || 'Arquivo'}
              </div>
            )}
          </div>

          {error && <div className="error">{error}</div>}

          <button 
            type="submit" 
            className="btn" 
            disabled={!file || uploading}
            style={{ width: '100%', marginTop: '15px' }}
          >
            {uploading ? '⏳ Processando...' : '📤 Fazer Upload'}
          </button>
        </form>

        {result && (
          <div style={{ 
            marginTop: '20px', 
            padding: '15px', 
            backgroundColor: '#d4edda', 
            border: '1px solid #c3e6cb', 
            borderRadius: '4px',
            color: '#155724'
          }}>
            <h3>✅ Upload Concluído com Sucesso!</h3>
            <p><strong>Mensagem:</strong> {result.message}</p>
            <div style={{ marginTop: '10px' }}>
              <strong>Dados processados:</strong>
              <ul style={{ margin: '5px 0', paddingLeft: '20px' }}>
                <li>Parceiros: {result.data.partners}</li>
                <li>Clientes: {result.data.customers}</li>
                <li>Produtos: {result.data.products}</li>
                <li>Registros de uso: {result.data.usages}</li>
              </ul>
            </div>
          </div>
        )}
      </div>

      <div className="card">
        <h2>📋 Instruções</h2>
        <div style={{ marginBottom: '15px' }}>
          <h3>Formatos Suportados:</h3>
          <ul>
            <li><strong>CSV (.csv)</strong> - Arquivo separado por vírgulas</li>
            <li><strong>Excel (.xlsx)</strong> - Planilha do Microsoft Excel</li>
          </ul>
        </div>

        <div style={{ marginBottom: '15px' }}>
          <h3>Colunas Obrigatórias:</h3>
          <ul>
            <li><code>partner_id</code> - ID do parceiro</li>
            <li><code>customer_id</code> - ID do cliente</li>
            <li><code>product_id</code> - ID do produto</li>
            <li><code>usage_date</code> - Data do uso</li>
            <li><code>quantity</code> - Quantidade</li>
            <li><code>unit_price</code> - Preço unitário</li>
          </ul>
        </div>

        <div>
          <h3>Colunas Opcionais:</h3>
          <ul>
            <li><code>partner_name</code>, <code>mpn_id</code>, <code>tier2_mpn_id</code></li>
            <li><code>customer_name</code>, <code>customer_domain_name</code>, <code>country</code></li>
            <li><code>sku_id</code>, <code>sku_name</code>, <code>product_name</code>, <code>meter_type</code>, <code>category</code>, <code>sub_category</code>, <code>unit_type</code></li>
            <li><code>invoice_number</code>, <code>charge_start_date</code>, <code>billing_pre_tax_total</code>, <code>resource_location</code>, <code>tags</code>, <code>benefit_type</code></li>
          </ul>
        </div>
      </div>

      <div className="card">
        <h2>💡 Dicas</h2>
        <ul>
          <li>Certifique-se de que a primeira linha contém os cabeçalhos das colunas</li>
          <li>Use vírgulas como separador decimal (ex: 1,50)</li>
          <li>Datas podem estar nos formatos: YYYY-MM-DD, DD/MM/YYYY, DD-MM-YYYY</li>
          <li>Para Excel, datas em formato serial são convertidas automaticamente</li>
          <li>Arquivos grandes são processados em lotes para melhor performance</li>
        </ul>
      </div>
    </div>
  );
}

export default Upload;
