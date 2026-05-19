import { useState, useEffect, useCallback } from 'react';

type Theme = 'light' | 'dark';

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    // 1. Check if manual preference exists
    const manual = localStorage.getItem('theme-manual');
    const saved = localStorage.getItem('theme') as Theme | null;
    if (manual && saved) return saved;

    // 2. Auto night check (18h - 06h)
    const hour = new Date().getHours();
    if (hour >= 18 || hour < 6) return 'dark';

    // 3. Telegram SDK
    const tgScheme = window.Telegram?.WebApp?.colorScheme;
    if (tgScheme === 'light' || tgScheme === 'dark') return tgScheme;

    // 4. Device preference
    if (window.matchMedia?.('(prefers-color-scheme: light)').matches) return 'light';
    return 'dark';
  });

  useEffect(() => {
    // Periodic check for night time (every minute)
    const interval = setInterval(() => {
      if (!localStorage.getItem('theme-manual')) {
        const hour = new Date().getHours();
        const isNight = hour >= 18 || hour < 6;
        setThemeState(isNight ? 'dark' : 'light');
      }
    }, 60000);

    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);

    // Update Telegram header/bg
    const tg = window.Telegram?.WebApp;
    if (tg) {
      try {
        tg.setHeaderColor(theme === 'light' ? '#faf6ee' : '#0f1115');
        tg.setBackgroundColor(theme === 'light' ? '#faf6ee' : '#0f1115');
      } catch {}
    }

    // Meta theme-color
    let meta = document.querySelector('meta[name="theme-color"]') as HTMLMetaElement;
    if (!meta) {
      meta = document.createElement('meta');
      meta.name = 'theme-color';
      document.head.appendChild(meta);
    }
    meta.content = theme === 'light' ? '#faf6ee' : '#0f1115';
  }, [theme]);

  // Listen to device changes
  useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)');
    const handler = (e: MediaQueryListEvent) => {
      if (!localStorage.getItem('theme-manual')) {
        setThemeState(e.matches ? 'dark' : 'light');
      }
    };
    mq.addEventListener('change', handler);
    return () => mq.removeEventListener('change', handler);
  }, []);

  const toggleTheme = useCallback(() => {
    setThemeState(prev => {
      const next = prev === 'dark' ? 'light' : 'dark';
      localStorage.setItem('theme-manual', '1');
      return next;
    });
  }, []);

  return { theme, toggleTheme };
}
