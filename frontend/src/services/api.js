import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'https://data-importer-api-go.onrender.com',
  timeout: 10000,
});

// Função para obter produtos por parceiro
export const getPartnerProducts = (partnerId) => {
  return api.get(`/api/partners/${partnerId}/products`);
};

// Interceptor para adicionar token automaticamente
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Interceptor para tratar erros de autenticação
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.log('API Error:', error.response?.status, error.response?.data);
    
    // Só redirecionar para login se for erro 401 e não estiver já na página de login
    if (error.response?.status === 401 && !window.location.pathname.includes('/login')) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    
    // Para outros erros, não redirecionar automaticamente
    return Promise.reject(error);
  }
);

export default api;
