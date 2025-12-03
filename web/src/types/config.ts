export type Service = {
  name: string;
  url: string;
  icon?: string;
  onlineBadge?: boolean;
};

export type BackgroundConfig = {
  image?: string;
  blur?: number;
};

export type UIConfig = {
  background?: BackgroundConfig;
  font?: string;
};

export type HerbstConfig = {
  title: string;
  ui: UIConfig;
  services: Service[];
  theme: string;
  themeVars: Record<string, string>;
};
