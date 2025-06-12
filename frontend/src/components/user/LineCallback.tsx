import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from '@/axiosConfig';

const LineCallback: React.FC = () => {
  const nav = useNavigate();

  useEffect(() => {
    (async () => {
      const qs = new URLSearchParams(location.search);
      try {
        const { data } = await axios.get('/users/auth/line/callback', {
          params: { code: qs.get('code'), state: qs.get('state') },
        });
        localStorage.setItem('token', data.token);
        nav('/mypage');
      } catch {
        nav('/signin?err=line');
      }
    })();
  }, [nav]);

  return <p>LINE 認証処理中...</p>;
};

export default LineCallback;
