export type Service = {
  name: string;
  url: string;
  icon?: string;
};

export type UIConfig = {
  background?: string;
};

export type HerbstConfig = {
  title: string;
  ui: UIConfig;
  services: Service[];
  theme: string;
  themeVars: Record<string, string>;
};
