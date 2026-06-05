export type AppView = 'create-frame' | 'copy-image' | 'frame-library' | 'settings';

export type NavItem = {
  id: AppView;
  label: string;
  icon: string;
  disabled?: boolean;
  description?: string;
};

export const NAV_ITEMS: NavItem[] = [
  {
    id: 'create-frame',
    label: 'Create Frame',
    icon: '🖼️',
  },
  {
    id: 'copy-image',
    label: 'Copy Image',
    icon: '📋',
  },
  {
    id: 'frame-library',
    label: 'Frame Library',
    icon: '📚',
    disabled: true,
    description: 'Browse and manage your saved frame templates in one place.',
  },
  {
    id: 'settings',
    label: 'Settings',
    icon: '⚙️',
    disabled: true,
    description: 'Configure app preferences, defaults, and output options.',
  },
];
