import React from "react";
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Customers from './pages/Customers';
import Reports from './pages/Reports';
import Metrics from './pages/Metrics';
// Upload removido: fluxo agora importa automaticamente no login
import { AuthProvider, useAuth } from './services/AuthContext';
import './index.css';

function Navigation() {
  const location = useLocation();
  const { isAuthenticated, logout } = useAuth();

  if (!isAuthenticated) {
    return null;
  }

  return (
    <nav className="nav">
      <div className="container">
        <ul>
          <li>
            <Link 
              to="/" 
              className={location.pathname === '/' ? 'active' : ''}
            >
              Dashboard
            </Link>
          </li>
          <li>
            <Link 
              to="/customers" 
              className={location.pathname === '/customers' ? 'active' : ''}
            >
              Clientes
            </Link>
          </li>
          <li>
            <Link 
              to="/reports" 
              className={location.pathname === '/reports' ? 'active' : ''}
            >
              Relatórios
            </Link>
          </li>
          <li>
            <Link 
              to="/metrics" 
              className={location.pathname === '/metrics' ? 'active' : ''}
            >
              Métricas
            </Link>
          </li>
          
          <li>
            <button onClick={logout} className="btn">
              Sair
            </button>
          </li>
        </ul>
      </div>
    </nav>
  );
}

function AppContent() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="App">
      <header className="header">
        <div className="container">
          <h1>Data Importer Dashboard</h1>
          <p>Visualização de dados de faturamento e uso</p>
        </div>
      </header>

      <Navigation />

      <main className="container">
        <Routes>
          <Route 
            path="/login" 
            element={isAuthenticated ? <Dashboard /> : <Login />} 
          />
          <Route 
            path="/" 
            element={isAuthenticated ? <Dashboard /> : <Login />} 
          />
          <Route 
            path="/customers" 
            element={isAuthenticated ? <Customers /> : <Login />} 
          />
          <Route 
            path="/reports" 
            element={isAuthenticated ? <Reports /> : <Login />} 
          />
          <Route 
            path="/metrics" 
            element={isAuthenticated ? <Metrics /> : <Login />} 
          />
          
        </Routes>
      </main>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <AppContent />
      </Router>
    </AuthProvider>
  );
}

export default App;
