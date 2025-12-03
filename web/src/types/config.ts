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

export type UIConfig = {
  background?: BackgroundConfig;
  font?: string;
};

export type HerbstConfig = {
  title: string;
  ui: UIConfig;
  weather: WeatherConfig;
  services: Service[];
  theme: string;
  themeVars: Record<string, string>;
};
