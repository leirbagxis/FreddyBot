import { useState, useEffect } from 'react';
import { fetchBotConfig, getBotId } from '../api';

export interface BotConfig {
  bot_username: string;
  bot_name: string;
  currency_name: string;
  currency_symbol: string;
}

const DEFAULT_CONFIG: BotConfig = {
  bot_username: 'Bot',
  bot_name: 'Meu Bot',
  currency_name: 'Reais',
  currency_symbol: 'R$'
};

export function useBotConfig() {
  const botId = getBotId();
  const [config, setConfig] = useState<BotConfig>(() => {
    const cached = localStorage.getItem(`bot_${botId}_config`);
    return cached ? JSON.parse(cached) : DEFAULT_CONFIG;
  });

  useEffect(() => {
    if (!botId) return;

    const loadConfig = async () => {
      try {
        const data = await fetchBotConfig(botId);
        if (data) {
          setConfig(data);
          localStorage.setItem(`bot_${botId}_config`, JSON.stringify(data));
        }
      } catch (err) {
        console.error("Failed to load bot config", err);
      }
    };

    loadConfig();
  }, [botId]);

  return config;
}
