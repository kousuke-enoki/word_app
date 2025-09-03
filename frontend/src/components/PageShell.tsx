import { Link } from 'react-router-dom'

import { useTheme } from '@/contexts/themeContext'

export const PageShell: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="min-h-screen bg-[var(--bg)] text-[var(--fg)]">
      <Header />
      <main>{children}</main>
      <Footer />
    </div>
  )
}

const Header: React.FC = () => {
  return (
    <header className="sticky top-0 z-40 border-b border-[var(--border)] bg-[var(--bg)]/80 backdrop-blur">
      <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-3">
        <Link to="/" className="text-lg font-bold tracking-wide">
          <span className="mr-1">📘</span> word app
        </Link>
        <nav className="flex items-center gap-2">
          <Link className="rounded-lg px-3 py-1.5 text-sm hover:bg-[var(--container_bg)]" to="/sign_in">
            サインイン
          </Link>
          <ThemeToggle />
        </nav>
      </div>
    </header>
  )
}

const ThemeToggle: React.FC = () => {
const { theme, setTheme } = useTheme()
const next = theme === 'dark' ? 'light' : 'dark'
  return (
    <button
      onClick={() => setTheme(next)}
      className="rounded-lg border border-[var(--border)] px-3 py-1.5 text-sm hover:bg-[var(--container_bg)]"
      aria-label="テーマ切替"
      title="テーマ切替"
    >
      {theme === 'dark' ? '🌞 Light' : '🌙 Dark'}
    </button>
  )
}



const Footer: React.FC = () => (
  <footer className="mt-16 border-t border-[var(--border)]">
    <div className="mx-auto max-w-5xl px-4 py-8 text-xs opacity-70">
      © {new Date().getFullYear()} DictQuiz / word app
    </div>
  </footer>
)
