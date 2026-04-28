/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideInRight: {
          '0%': { transform: 'translateX(100%)', opacity: '0' },
          '100%': { transform: 'translateX(0)', opacity: '1' },
        }
      },
      animation: {
        'fade-in': 'fadeIn 0.2s ease-out forwards',
        'slide-up': 'slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards',
        'slide-in': 'slideInRight 0.3s ease-out forwards',
      },
      fontFamily: {
        sans: ['"Inter"', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', '"Fira Code"', 'monospace'],
      },
      colors: {
        surface: '#050505', // Fundo principal
        panel: '#0f0f0f',   // Painéis e Cards
        primary: '#D9FF00', // Verde Neon (Principal)
        secondary: '#00F0FF', // Ciano Neon
        accent: '#FF003C',  // Vermelho/Rosa Neon
        muted: '#444444',
        text: '#E0E0E0',
        border: '#1a1a1a',
        success: '#22C55E',
        warning: '#FACC15',
        error: '#EF4444',
      },
      backgroundImage: {
        'neon-gradient': 'linear-gradient(180deg, rgba(217, 255, 0, 0.1) 0%, rgba(0, 0, 0, 0) 100%)',
        'cyber-grid': 'radial-gradient(circle, #1a1a1a 1px, transparent 1px)',
      },
      boxShadow: {
        'neon': '0 0 10px rgba(217, 255, 0, 0.5), 0 0 20px rgba(217, 255, 0, 0.2)',
        'neon-strong': '0 0 15px rgba(217, 255, 0, 0.8), 0 0 30px rgba(217, 255, 0, 0.4)',
        'accent-neon': '0 0 10px rgba(255, 0, 60, 0.5)',
      },
      spacing: {
        xs: '4px',
        sm: '8px',
        md: '16px',
        lg: '24px',
        xl: '48px',
      },
    },
  },
  plugins: [],
}
