<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import Logo from "./Logo.vue";
import type { WeatherConfig } from "../types/config";

const props = defineProps<{
  title: string;
  weather: WeatherConfig;
}>();

const time = ref("");
const date = ref("");
const cpuLoad = ref("0%");
const uptime = ref("0h");
const isOnline = ref(true);
const use24h = ref(true);

// Weather state
const weatherData = ref<{
  temp: number;
  feelsLike: number;
  humidity: number;
  description: string;
  icon: string;
  city: string;
  units: string;
} | null>(null);
const weatherError = ref<string | null>(null);
let weatherInterval: ReturnType<typeof setInterval> | null = null;

function updateClock() {
  const now = new Date();

  time.value = now.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: !use24h.value,
  });

  date.value = now.toLocaleDateString([], {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
}

function toggleHourFormat() {
  use24h.value = !use24h.value;
  updateClock();
}

// Fake CPU/Uptime (du kannst später echte Endpoints reinballern)
function updateSystemInfo() {
  cpuLoad.value = Math.round(Math.random() * 20) + "%";
  uptime.value = `${Math.round(Math.random() * 48)}h`;
}

async function fetchWeather() {
  try {
    const response = await fetch("/api/weather");
    const data = await response.json();

    if (data.enabled && !data.error) {
      weatherData.value = {
        temp: data.temp,
        feelsLike: data.feelsLike,
        humidity: data.humidity,
        description: data.description,
        icon: data.icon,
        city: data.city,
        units: data.units,
      };
      weatherError.value = null;
    } else if (data.error) {
      weatherError.value = data.error;
      weatherData.value = null;
    }
  } catch (err) {
    weatherError.value = "Failed to fetch weather";
    weatherData.value = null;
  }
}

function getWeatherIconUrl(iconCode: string): string {
  return `https://openweathermap.org/img/wn/${iconCode}@2x.png`;
}

function formatTemp(temp: number, units: string): string {
  const rounded = Math.round(temp);
  switch (units) {
    case "metric":
      return `${rounded}°C`;
    case "imperial":
      return `${rounded}°F`;
    default:
      return `${rounded}K`;
  }
}

onMounted(() => {
  updateClock();
  updateSystemInfo();
  setInterval(updateClock, 1000);
  setInterval(updateSystemInfo, 5000);

  // Fetch weather immediately and then every hour
  fetchWeather();
  weatherInterval = setInterval(fetchWeather, 60 * 60 * 1000); // 1 hour
});

// Refetch weather when config changes (location, units, etc.)
watch(
  () => [props.weather.location, props.weather.units, props.weather.enabled],
  () => {
    if (props.weather.enabled) {
      fetchWeather();
    } else {
      weatherData.value = null;
    }
  }
);

onUnmounted(() => {
  if (weatherInterval) {
    clearInterval(weatherInterval);
  }
});
</script>

<template>
  <header class="header-bar">
    <div class="header-left">
      <div class="logo-icon">
        <Logo />
      </div>
      <span class="title">{{ title }}</span>
    </div>

    <div class="status-bar">
      <!-- Weather -->
      <div v-if="weatherData" class="weather-info">
        <img
          :src="getWeatherIconUrl(weatherData.icon)"
          :alt="weatherData.description"
          class="weather-icon"
        />
        <span class="weather-temp">{{
          formatTemp(weatherData.temp, weatherData.units)
        }}</span>
        <span class="weather-divider">·</span>
        <span class="weather-city">{{ weatherData.city }}</span>
      </div>

      <!-- Online Badge -->
      <span class="dot" :class="{ online: isOnline }"></span>

      <!-- Time -->
      <span class="time" @click="toggleHourFormat">{{ time }}</span>

      <!-- Date -->
      <span class="date">{{ date }}</span>

      <!-- CPU + Uptime -->
      <span class="meta">CPU {{ cpuLoad }} · UP {{ uptime }}</span>
    </div>
  </header>
</template>

<style scoped>
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  padding-left: 12px; /* Less padding on left for logo integration */
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  border-radius: 3em;

  /* Backdrop blur for transparent themes */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  padding: 8px;
  box-sizing: border-box;
  margin-left: -2px; /* Slight overlap with header edge */
}

.title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--color-text);
}

.status-bar {
  display: flex;
  align-items: center;
  gap: 12px; /* perfekt */
  font-size: 0.85rem;
  color: var(--color-text-muted);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-text-muted);
  transition: background 0.2s;
}

.dot.online {
  background: #4ade80; /* green-400 */
}

.time {
  font-weight: 600;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
}
.time:hover {
  color: var(--color-accent);
}

.date {
  color: var(--color-text-muted);
}

.meta {
  color: var(--color-text-muted);
  font-size: 0.78rem;
  opacity: 0.8;
}

.weather-info {
  display: flex;
  align-items: center;
  gap: 4px;
  padding-right: 8px;
  border-right: 1px solid var(--color-border);
  margin-right: 4px;
}

.weather-icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
  margin: -4px;
}

.weather-temp {
  font-weight: 600;
  color: var(--color-text);
}

.weather-divider {
  color: var(--color-text-muted);
  opacity: 0.6;
}

.weather-city {
  color: var(--color-text-muted);
  font-size: 0.8rem;
}
</style>
