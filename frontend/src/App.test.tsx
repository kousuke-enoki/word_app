import { render, screen, within } from '@testing-library/react'
import React from 'react'

import App from './App'

test('renders brand in header', () => {
  render(<App />)
  const header = screen.getByRole('banner') // <header> のランドマーク
  const brandLink = within(header).getByRole('link', { name: /DictQuiz/i })
  expect(brandLink).toBeInTheDocument()
})
