import React from "react";
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Customers from './pages/Customers';
import Reports from './pages/Reports';
import Upload from './pages/Upload';
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
              to="/upload" 
              className={location.pathname === '/upload' ? 'active' : ''}
            >
              Upload
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
            path="/upload" 
            element={isAuthenticated ? <Upload /> : <Login />} 
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
