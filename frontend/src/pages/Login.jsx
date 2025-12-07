import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/api';

export default function Login() {
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('password123');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const res = await login(username, password);
      // Token'ı tarayıcı hafızasına kaydet
      localStorage.setItem('token', res.data.token); 
      navigate('/map'); // Haritaya yönlendir
    } catch (err) {
      setError('Giriş başarısız! Kullanıcı adı veya şifre yanlış olabilir.');
    }
  };

  return (
    <div className="flex h-screen items-center justify-center bg-gray-900">
      <form onSubmit={handleLogin} className="w-96 bg-white p-8 rounded-lg shadow-xl">
        <h1 className="text-3xl font-bold text-center mb-6 text-gray-800">TaxiHub 🚕</h1>
        
        {error && <div className="bg-red-100 text-red-700 p-2 mb-4 rounded text-sm">{error}</div>}

        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-bold mb-2">Kullanıcı Adı</label>
          <input
            className="w-full p-3 border rounded focus:outline-none focus:ring-2 focus:ring-taxi"
            placeholder="admin"
            value={username}
            onChange={e => setUsername(e.target.value)}
          />
        </div>
        
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-bold mb-2">Şifre</label>
          <input
            className="w-full p-3 border rounded focus:outline-none focus:ring-2 focus:ring-taxi"
            type="password"
            placeholder="password123"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        </div>

        <button className="w-full bg-taxi text-gray-900 font-bold p-3 rounded hover:bg-yellow-400 transition duration-200">
          GİRİŞ YAP
        </button>
      </form>
    </div>
  );
}