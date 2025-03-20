import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axiosInstance from '../../axiosConfig'

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')
  const navigate = useNavigate()

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

      setTimeout(() => {
        navigate('/mypage')
      })
    } catch (error) {
      setMessage('Sign in failed. Please try again.')
    }
  }

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
    </div>
  )
}

export default SignIn
