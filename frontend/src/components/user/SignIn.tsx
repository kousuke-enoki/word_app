import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Link } from 'react-router-dom'
import axiosInstance from '@/axiosConfig'
import { useTheme } from '@/contexts/themeContext'


type SettingResponse = {
  isLineAuth: boolean;
};

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')
  const [lineAuthEnabled, setLineAuthEnabled] = useState<boolean>(false);
  const [loadingSetting,  setLoadingSetting]  = useState<boolean>(true);
  const navigate = useNavigate()
  const { setTheme } = useTheme()

  useEffect(() => {
    let isMounted = true;                         // アンマウント対策

    (async () => {
      try {
        const { data } = await axiosInstance.get<SettingResponse>(
          '/setting/auth',
        );
        if (!isMounted) return;
        setLineAuthEnabled(data.isLineAuth);
      } finally {
        // eslint-disable-next-line @typescript-eslint/no-unused-expressions
        isMounted && setLoadingSetting(false);
      }
    })();

    return () => {
      isMounted = false;
    };
  }, [setTheme]);

  const handleSignIn = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    try {
      const response = await axiosInstance.post('/users/sign_in', {
        email,
        password,
      })
      const token = response.data.token
      localStorage.setItem('token', token)
      localStorage.setItem('logoutMessage', 'サインイン成功！')
      const res = await axiosInstance.get('/setting/user_config')
      setTheme(res.data.is_dark_mode ? 'dark' : 'light')

      setTimeout(() => {
        navigate('/mypage')
      })
    } catch {
      setMessage('Sign in failed. Please try again.')
    }
  }

  const handleLineLogin = () => {
    window.location.href = `${import.meta.env.VITE_API_URL}/users/auth/line/login`;
  };

  if (loadingSetting) return <p>Loading…</p>;

  return (
    <div>
      <h1>サインイン</h1>
      <form onSubmit={handleSignIn}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit">サインイン</button>
      </form>
      {message && <p>{message}</p>}
      <div>
        {lineAuthEnabled && (
          <button type="button" onClick={handleLineLogin}>
            LINEでログイン
          </button>
        )}
      </div>
      <div>
        <p>
          <Link to="/sign_up">サインアップはここから！</Link>
        </p>
      </div>
    </div>
  )
}

export default SignIn
