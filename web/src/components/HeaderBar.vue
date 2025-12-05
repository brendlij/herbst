<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from "vue";
import Logo from "./Logo.vue";
import type { WeatherConfig } from "../types/config";

const props = defineProps<{
  title: string;
  weather: WeatherConfig;
  activeTab: string;
  dockerEnabled: boolean;
  dockerAgentsConfigured: boolean;
}>();

const emit = defineEmits<{
  (e: "tabChange", tab: string): void;
}>();

const time = ref("");
const date = ref("");
const use24h = ref(true);
const useGermanDate = ref(true); // true = "3. Dez 2025", false = "03.12.2025"

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

  time.value = now.toLocaleTimeString("de-DE", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: !use24h.value,
  });

  if (useGermanDate.value) {
    // German format: "3. Dez 2025"
    date.value = now.toLocaleDateString("de-DE", {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
  } else {
    // Numeric format: "03.12.2025"
    date.value = now.toLocaleDateString("de-DE", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  }
}

function toggleHourFormat() {
  use24h.value = !use24h.value;
  updateClock();
}

function toggleDateFormat() {
  useGermanDate.value = !useGermanDate.value;
  updateClock();
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
  setInterval(updateClock, 1000);

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
    <!-- Left section - logo and title -->
    <div class="header-left">
      <div class="logo-icon">
        <Logo color="var(--color-accent)" />
      </div>
      <span class="title">{{ title }}</span>
    </div>

    <!-- Center section - tabs -->
    <nav class="header-tabs">
      <div
        class="tab"
        :class="{ active: activeTab === 'services' }"
        @click="emit('tabChange', 'services')"
      >
        Services
      </div>
      <div
        v-if="dockerEnabled"
        class="tab"
        :class="{ active: activeTab === 'docker' }"
        @click="emit('tabChange', 'docker')"
      >
        <div>Docker</div>
        <div class="tab-subtitle">Local</div>
      </div>
      <div
        v-if="dockerAgentsConfigured"
        class="tab"
        :class="{ active: activeTab === 'docker-nodes' }"
        @click="emit('tabChange', 'docker-nodes')"
      >
        <div>Docker</div>
        <div class="tab-subtitle">Agents</div>
      </div>
      <div
        class="tab"
        :class="{ active: activeTab === 'config' }"
        @click="emit('tabChange', 'config')"
      >
        <div>Configuration</div>
      </div>
    </nav>

    <!-- Right section - weather and datetime -->
    <div class="header-right">
      <!-- Weather -->
      <div v-if="weatherData" class="weather-info">
        <img
          :src="getWeatherIconUrl(weatherData.icon)"
          :alt="weatherData.description"
          class="weather-icon"
        />
        <div class="weather-text">
          <span class="weather-temp">{{
            formatTemp(weatherData.temp, weatherData.units)
          }}</span>
          <span class="weather-city">{{ weatherData.city }}</span>
        </div>
      </div>

      <!-- DateTime stacked -->
      <div class="datetime">
        <span class="time" @click="toggleHourFormat">{{ time }}</span>
        <span class="date" @click="toggleDateFormat">{{ date }}</span>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  padding-left: 12px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  border-radius: 2em;

  /* Backdrop blur for transparent themes */
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.header-tabs {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.tab {
  padding: 10px 18px;
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--color-text-muted);
  border-radius: 10px;
  cursor: pointer;
  position: relative;
  transition: 0.2s ease;
  user-select: none;
}

.tab:hover {
  color: var(--color-text);
  background: var(--color-bg);
}

.tab.active {
  color: var(--color-text);
}

.tab.active::after {
  content: "";
  position: absolute;
  bottom: -2px;
  left: 20%;
  right: 20%;
  height: 2px;
  background: var(--color-accent);
  border-radius: 4px;
}

.tab-subtitle {
  position: absolute;
  top: 27px;
  font-size: 0.6rem;
  color: var(--color-text-muted);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  justify-content: flex-end;
  min-width: 0;
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
  margin-left: -2px;
}

.title {
  font-size: 1.3rem;
  font-weight: 600;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.datetime {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  font-size: 0.85rem;
}

.time {
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s;
  color: var(--color-text);
}

.time:hover {
  color: var(--color-accent);
}

.date {
  color: var(--color-text);
  opacity: 0.6;
  font-size: 0.8rem;
  cursor: pointer;
  user-select: none;
  transition: color 0.2s, opacity 0.2s;
}

.date:hover {
  color: var(--color-accent);
  opacity: 1;
}

.weather-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-right: 16px;
  border-right: 1px solid var(--color-border);
}

.weather-icon {
  width: 36px;
  height: 36px;
  object-fit: contain;
}

.weather-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}

.weather-temp {
  font-weight: 600;
  font-size: 1rem;
  color: var(--color-text);
}

.weather-city {
  color: var(--color-text);
  opacity: 0.6;
  font-size: 0.8rem;
}
</style>
