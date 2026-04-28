import { Users } from 'lucide-react';

export const NAV_ITEMS = [
  { id: 'users', label: 'USUÁRIOS', icon: Users },
];

export const DEFAULT_COMMANDS = [
  {
    id: 'std_start',
    trigger: '/start',
    text: '👋 <b>Olá, $user_name!</b>\n\nBem-vindo ao seu assistente digital inteligente. Estou aqui para ajudar a gerenciar suas atividades, economia e inventário.\n\n⚡ <b>Status do Sistema:</b> 🟢 Ativo\n🤖 <b>ID do Bot:</b> <code>$bot_id</code>\n\n📱 Use os botões abaixo para acessar sua central de controle ou ver os comandos disponíveis.',
    buttons: [
      [{id: 'dash', text: '🖥️ Abrir Dashboard', type: 'webapp', url: '{dashboard}'}],
      [{id: 'help', text: '📖 Comandos', type: 'callback', data: 'help:main'}, {id: 'profile', text: '👤 Meu Perfil', type: 'callback', data: 'profile:back'}]
    ],
    action: 'send'
  },
  {
    id: 'std_saldo',
    trigger: '/saldo',
    text: '💰 Olá <b>$user_name</b>, seu saldo atual é: <b>$user_saldo</b>',
    buttons: [[]],
    action: 'send'
  },
  {
    id: 'std_perfil',
    trigger: '/perfil',
    text: '👤 <b>PERFIL DO USUÁRIO</b>\n\n🏷 <b>Nome:</b> $user_name\n🆔 <b>ID:</b> <code>$user_id</code>\n💰 <b>Saldo:</b> $user_saldo\n\n📊 <b>RANKINGS</b>\n💬 <b>Mensagens:</b> #?\n🏆 <b>Riqueza:</b> #?',
    buttons: [
      [{ id: 'profile:inv', text: '🎒 Inventário', type: 'callback', data: 'profile:inv' }],
      [{ id: 'profile:ext', text: '📜 Extrato', type: 'callback', data: 'profile:ext' }]
    ],
    action: 'send'
  },
  {
    id: 'std_loja',
    trigger: '/loja',
    text: '🛒 <b>LOJA INTEGRADA</b>\n\nConfira nossos itens disponíveis e turbine sua experiência!',
    buttons: [[{ id: 'btn_shop', text: '🛍️ Abrir Loja', type: 'webapp', url: '{dashboard}/shop' }]],
    action: 'send'
  },
  {
    id: 'std_rank',
    trigger: '/rank',
    text: '📊 <b>Ranking de Mensagens</b>\n\nVeja quem são os usuários mais ativos do grupo!',
    buttons: [[{ id: 'btn_rank', text: '📊 Ver Ranking', type: 'webapp', url: '{dashboard}/ranking' }]],
    action: 'send'
  },
  {
    id: 'std_inventario',
    trigger: '/inventario',
    text: '🎒 <b>SEU INVENTÁRIO</b>\n\nOlá <b>$user_name</b>, aqui estão seus itens:',
    buttons: [[{ id: 'btn_inv', text: '📦 Ver Itens', type: 'webapp', url: '{dashboard}' }]],
    action: 'send'
  },
  {
    id: 'std_ajuda',
    trigger: '/ajuda',
    text: '📖 <b>COMANDOS DISPONÍVEIS</b>\n\n💰 <code>/saldo</code> - Veja seu saldo atual\n🎒 <code>/inventario</code> - Veja seus itens e use-os\n🛒 <code>/loja</code> - Abre a loja de itens\n🏆 <code>/rank</code> - Ranking de mensagens e riqueza\n👤 <code>/perfil</code> - Seu perfil completo e dashboard\n🎁 <code>/resgatar</code> - Lista e resgata recompensas\n\n💡 <i>Use o dashboard web para uma experiência completa!</i>',
    buttons: [[]],
    action: 'send'
  },
  {
    id: 'std_admin',
    trigger: '/admin',
    text: '🛠 <b>GESTOR DO BOT</b>\n\nOlá <b>$user_name</b>, aqui estão os comandos de administração:\n\n• <code>/addadmin [ID]</code> - Promover Admin\n• <code>/removeadmin [ID]</code> - Remover Admin\n• <code>/addmoney [ID] [VALOR]</code> - Adicionar Saldo\n• <code>/additem [ID] [ITEM_ID]</code> - Adicionar Item\n\nUse o painel para gestão completa:',
    buttons: [[{ id: 'btn_admin_webapp', text: '⚙️ Abrir Painel Admin', type: 'webapp', url: '{dashboard}/admin' }]],
    action: 'send'
  },
  {
    id: 'std_ajuda_admin',
    trigger: '/ajudaadmin',
    text: '📖 <b>COMANDOS DE ADMINISTRAÇÃO</b>\n\n• <code>/addadmin [ID]</code> - Promove usuário\n• <code>/removeadmin [ID]</code> - Remove privilégios\n• <code>/addmoney [ID] [QUANTIA]</code> - Adiciona saldo\n• <code>/additem [ID] [ITEM_ID]</code> - Adiciona item',
    buttons: [[{ id: 'btn_admin_go', text: '⚙️ Gerenciar no WebApp', type: 'webapp', url: '{dashboard}/admin' }]],
    action: 'send'
  },
  {
    id: 'cb_profile_inv',
    trigger: 'profile:inv',
    text: '🎒 <b>SEU INVENTÁRIO</b>\n\nAqui estão os itens que você possui atualmente:',
    buttons: [[{ id: 'profile:back', text: '🔙 Voltar ao Perfil', type: 'callback', data: 'profile:back' }]],
    action: 'edit'
  },
  {
    id: 'cb_profile_ext',
    trigger: 'profile:ext',
    text: '📜 <b>EXTRATO RECENTE</b>\n\nVeja suas últimas movimentações financeiras:',
    buttons: [[{ id: 'profile:back', text: '🔙 Voltar ao Perfil', type: 'callback', data: 'profile:back' }]],
    action: 'edit'
  },
  {
    id: 'cb_profile_back',
    trigger: 'profile:back',
    text: '👤 <b>PERFIL DO USUÁRIO</b>\n\n🏷 <b>Nome:</b> $user_name\n🆔 <b>ID:</b> <code>$user_id</code>\n💰 <b>Saldo:</b> $user_saldo',
    buttons: [
      [{ id: 'profile:inv', text: '🎒 Inventário', type: 'callback', data: 'profile:inv' }],
      [{ id: 'profile:ext', text: '📜 Extrato', type: 'callback', data: 'profile:ext' }]
    ],
    action: 'edit'
  },
  {
    id: 'cb_help_main',
    trigger: 'help:main',
    text: '📖 <b>COMO FUNCIONA?</b>\n\nNeste bot você pode ganhar moedas participando do grupo e comprar itens exclusivos na loja!',
    buttons: [[{ id: 'profile:back', text: '🔙 Voltar', type: 'callback', data: 'profile:back' }]],
    action: 'edit'
  }
];
