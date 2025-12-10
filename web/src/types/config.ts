export type Service = {
  name: string;
  url: string;
  icon?: string;
  onlineBadge?: boolean;
};

export type ServiceSection = {
  title: string;
  services: Service[];
};

export type BackgroundConfig = {
  image?: string;
  blur?: number;
};

export type WeatherConfig = {
  enabled: boolean;
  apiKey: string;
  location?: string; // City name, "zip:CODE,COUNTRY", or empty for lat/lon
  lat: number;
  lon: number;
  units: "metric" | "imperial" | "standard";
};

export type WeatherData = {
  temp: number;
  feelsLike: number;
  humidity: number;
  description: string;
  icon: string;
  city: string;
};

export type DockerConfig = {
  enabled: boolean;
  socketPath: string;
  agentsConfigured: boolean;
};

export type SystemConfig = {
  enabled: boolean;
  diskPath: string;
};

export type DockerContainer = {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  created: number;
};

export type ClockConfig = {
  timeFormat: "24h" | "12h";
  dateFormat: "short" | "numeric";
};

export type UIConfig = {
  background?: BackgroundConfig;
  font?: string;
  clock?: ClockConfig;
};

export type HerbstConfig = {
  title: string;
  ui: UIConfig;
  weather: WeatherConfig;
  docker: DockerConfig;
  system: SystemConfig;
  services: Service[];
  sections: ServiceSection[];
  theme: string;
  themeVars: Record<string, string>;
};
