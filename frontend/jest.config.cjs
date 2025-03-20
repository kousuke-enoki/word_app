module.exports = {
  roots: ['<rootDir>/src'],
  testMatch: [
    '**/__tests__/**/*.+(ts|tsx|js)',
    '**/?(*.)+(spec|test).+(ts|tsx|js)',
  ],
  transform: {
    '^.+\\.(ts|tsx)$': [
      '@swc/jest',
      {
        jsc: {
          parser: {
            syntax: 'typescript',
            tsx: true
          },
          transform: {
            react: {
              runtime: 'automatic'
              // or "classic"
            }
          }
        }
      }
    ]
  },
  transformIgnorePatterns: ['<rootDir>/node_modules/'],
  testEnvironment: 'jsdom',
  globals: {
    'ts-jest': {
      // Vite向けに isolatedModules: true のままでも動くケースあり
      // もし問題があれば "isolatedModules": false を検討
    },
  },
}
