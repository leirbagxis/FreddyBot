import { useState, useEffect, useCallback } from 'react';

type Theme = 'light' | 'dark' | 'telegram';
type BaseTheme = 'light' | 'dark';

const TELEGRAM_CSS_VARS = [
  '--bg', '--card', '--card-elevated', '--text', '--text-secondary', '--hint',
  '--accent', '--accent-rgb', '--accent-soft', '--accent-hover', '--accent-text',
  '--link', '--border', '--border-active', '--surface-rgb', '--surface',
  '--surface-hover', '--nav-bg', '--input-bg-rgb', '--input-bg', '--overlay',
  '--toggle-bg', '--toggle-knob',
];

const isTheme = (value: string | null): value is Theme => {
  return value === 'light' || value === 'dark' || value === 'telegram';
};

const getAutoTheme = (): BaseTheme => {
  const hour = new Date().getHours();
  if (hour >= 18 || hour < 6) return 'dark';

  const tgScheme = window.Telegram?.WebApp?.colorScheme;
  if (tgScheme === 'light' || tgScheme === 'dark') return tgScheme;

  if (window.matchMedia?.('(prefers-color-scheme: light)').matches) return 'light';
  return 'dark';
};

const hexToRgb = (hex: string, fallback: string) => {
  const normalized = hex.trim().replace('#', '');
  if (!/^[0-9a-fA-F]{6}$/.test(normalized)) return fallback;
  const value = parseInt(normalized, 16);
  return `${(value >> 16) & 255}, ${(value >> 8) & 255}, ${value & 255}`;
};

const applyTelegramVars = () => {
  const root = document.documentElement;
  const tg = window.Telegram?.WebApp;
  const params = tg?.themeParams || {};
  const scheme: BaseTheme = tg?.colorScheme === 'light' || tg?.colorScheme === 'dark'
    ? tg.colorScheme
    : getAutoTheme();

  const fallback = scheme === 'light'
    ? {
        bg: '#faf6ee', card: 'rgba(255, 255, 255, 0.7)', text: '#1e1b4b', hint: '#94a3b8',
        secondary: '#475569', accent: '#6366f1', border: 'rgba(99, 102, 241, 0.1)', surfaceRgb: '255, 255, 255',
        nav: 'rgba(248, 249, 255, 0.85)', overlay: 'rgba(0, 0, 0, 0.3)', inputRgb: '255, 255, 255',
      }
    : {
        bg: '#0a0c10', card: 'rgba(30, 41, 59, 0.4)', text: '#f1f5f9', hint: '#475569',
        secondary: '#94a3b8', accent: '#818cf8', border: 'rgba(255, 255, 255, 0.06)', surfaceRgb: '255, 255, 255',
        nav: 'rgba(10, 12, 16, 0.92)', overlay: 'rgba(0, 0, 0, 0.7)', inputRgb: '15, 23, 42',
      };

  const bg = params.bg_color || fallback.bg;
  const secondaryBg = params.secondary_bg_color || bg;
  const text = params.text_color || fallback.text;
  const hint = params.hint_color || fallback.hint;
  const link = params.link_color || fallback.accent;
  const button = params.button_color || link;
  const buttonText = params.button_text_color || '#ffffff';
  const surfaceRgb = hexToRgb(secondaryBg, fallback.surfaceRgb);
  const inputRgb = hexToRgb(secondaryBg, fallback.inputRgb);
  const accentRgb = hexToRgb(button, scheme === 'light' ? '99, 102, 241' : '129, 140, 248');

  root.style.setProperty('--bg', bg);
  root.style.setProperty('--card', `rgba(${surfaceRgb}, ${scheme === 'light' ? '0.78' : '0.42'})`);
  root.style.setProperty('--card-elevated', `rgba(${surfaceRgb}, ${scheme === 'light' ? '0.92' : '0.62'})`);
  root.style.setProperty('--text', text);
  root.style.setProperty('--text-secondary', params.subtitle_text_color || fallback.secondary);
  root.style.setProperty('--hint', hint);
  root.style.setProperty('--accent', button);
  root.style.setProperty('--accent-rgb', accentRgb);
  root.style.setProperty('--accent-soft', `rgba(${accentRgb}, 0.12)`);
  root.style.setProperty('--accent-hover', `rgba(${accentRgb}, 0.18)`);
  root.style.setProperty('--accent-text', buttonText);
  root.style.setProperty('--link', link);
  root.style.setProperty('--border', params.section_separator_color || fallback.border);
  root.style.setProperty('--border-active', button);
  root.style.setProperty('--surface-rgb', surfaceRgb);
  root.style.setProperty('--surface', `rgba(${surfaceRgb}, ${scheme === 'light' ? '0.5' : '0.08'})`);
  root.style.setProperty('--surface-hover', `rgba(${surfaceRgb}, ${scheme === 'light' ? '0.66' : '0.14'})`);
  root.style.setProperty('--nav-bg', params.header_bg_color || fallback.nav);
  root.style.setProperty('--input-bg-rgb', inputRgb);
  root.style.setProperty('--input-bg', `rgba(${inputRgb}, ${scheme === 'light' ? '0.55' : '0.45'})`);
  root.style.setProperty('--overlay', fallback.overlay);
  root.style.setProperty('--toggle-bg', params.secondary_bg_color || (scheme === 'light' ? '#e2e8f0' : '#1e293b'));
  root.style.setProperty('--toggle-knob', params.button_text_color || '#ffffff');

  return { bg, scheme };
};

const clearTelegramVars = () => {
  TELEGRAM_CSS_VARS.forEach(name => document.documentElement.style.removeProperty(name));
};

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    const manual = localStorage.getItem('theme-manual');
    const saved = localStorage.getItem('theme');
    if (manual && isTheme(saved)) return saved;
    return getAutoTheme();
  });

  useEffect(() => {
    const interval = setInterval(() => {
      if (!localStorage.getItem('theme-manual')) {
        setThemeState(getAutoTheme());
      }
    }, 60000);

    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);

    let bg = theme === 'light' ? '#faf6ee' : '#0f1115';
    if (theme === 'telegram') {
      bg = applyTelegramVars().bg;
    } else {
      clearTelegramVars();
    }

    const tg = window.Telegram?.WebApp;
    if (tg) {
      try {
        tg.setHeaderColor(bg);
        tg.setBackgroundColor(bg);
      } catch {}
    }

    let meta = document.querySelector('meta[name="theme-color"]') as HTMLMetaElement;
    if (!meta) {
      meta = document.createElement('meta');
      meta.name = 'theme-color';
      document.head.appendChild(meta);
    }
    meta.content = bg;
  }, [theme]);

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
      const next: Theme = prev === 'light' ? 'dark' : prev === 'dark' ? 'telegram' : 'light';
      localStorage.setItem('theme-manual', '1');
      return next;
    });
  }, []);

  return { theme, toggleTheme };
}
