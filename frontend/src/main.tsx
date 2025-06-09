import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
// import ReactDOM from 'react-dom/client'
import './index.css'
import App from './App'
// import reportWebVitals from './reportWebVitals'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
// ダークモード等の切り替え用
import { ThemeProvider } from '@/context/ThemeProvider'

// const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement)
const queryClient = new QueryClient()
// root.render(

createRoot(document.getElementById('root')!).render(

  <StrictMode>
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
      </ThemeProvider>
    </StrictMode>,
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
// reportWebVitals()
