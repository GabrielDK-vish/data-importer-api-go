import React, { useState } from 'react';
import { useAuth } from '../services/AuthContext';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    const result = await login(username, password);
    
    if (!result.success) {
      setError(result.error);
    }
    
    setLoading(false);
  };

  return (
    <div className="card" style={{ maxWidth: '400px', margin: '50px auto' }}>
      <h2>🔐 Login</h2>
      <p style={{ marginBottom: '20px', color: '#666' }}>
        Faça login para acessar o dashboard
      </p>
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Usuário:</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="admin, user ou demo"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="password">Senha:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="admin123, user123 ou demo123"
            required
          />
        </div>

        {error && <div className="error">{error}</div>}

        <button 
          type="submit" 
          className="btn" 
          disabled={loading}
          style={{ width: '100%' }}
        >
          {loading ? 'Entrando...' : 'Entrar'}
        </button>
      </form>

      <div style={{ marginTop: '20px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '4px' }}>
        <h4>Credenciais de teste:</h4>
        <ul style={{ margin: '10px 0', paddingLeft: '20px' }}>
          <li><strong>admin</strong> / admin123</li>
          <li><strong>user</strong> / user123</li>
          <li><strong>demo</strong> / demo123</li>
        </ul>
      </div>
    </div>
  );
}

export default Login;
