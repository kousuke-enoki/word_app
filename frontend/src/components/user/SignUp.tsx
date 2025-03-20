import React, { useState } from 'react'
import axiosInstance from '../../axiosConfig'
import { useNavigate } from 'react-router-dom'

interface FieldError {
  field: string
  message: string
}

const SignUp: React.FC = () => {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')
  const [errors, setErrors] = useState<FieldError[]>([])
  const navigate = useNavigate()

  const handleSignUp = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    try {
      const response = await axiosInstance.post('/users/sign_up', {
        name,
        email,
        password,
      })
      const token = response.data.token
      localStorage.setItem('token', token)
      localStorage.setItem('logoutMessage', 'サインアップしました。')
      setMessage('Sign up successful!')
      setErrors([])
      setTimeout(() => {
        navigate('/mypage')
      })
    } catch (error: any) {
      const fieldErrors: FieldError[] = error.response?.data?.errors || []
      setErrors(fieldErrors)
      setMessage('')
    }
  }

  // フィールドごとにエラーメッセージを取得する関数
  const getErrorMessages = (field: string) => {
    return errors
      .filter((e) => e.field === field)
      .map((e, index) => (
        <p key={index} style={{ color: 'red' }}>
          {e.message}
        </p>
      ))
  }

  return (
    <div>
      <h1>サインアップ</h1>
      <form onSubmit={handleSignUp}>
        <div>
          <label htmlFor="name">Name:</label>
          <input
            type="name"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
          {getErrorMessages('name')}
        </div>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          {getErrorMessages('email')}
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
          {getErrorMessages('password')}
        </div>
        <button type="submit">サインアップ</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  )
}

export default SignUp
