<template>
  <a-config-provider :theme="themeConfig" :locale="arcoLocale">
    <a-layout style="height: 100vh">
      <a-layout-sider :style="{
        height: '100vh',
        boxShadow: '2px 0 8px 0 rgba(29, 35, 41, 0.05)',
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0,
        zIndex: 100,
      }" :collapsed="collapsed" :collapsible="!darkTheme" :width="190" breakpoint="lg" @collapse="onCollapse">
        <div class="logo-container" :class="{ 'collapsed-logo': collapsed }">
          <img src="@/assets/images/maxx_logo.svg" alt="logo" class="logo" />
        </div>
        <a-menu v-model:selected-keys="selectedKeys" :collapsed="collapsed" :style="{ borderRight: 'none' }"
          @menu-item-click="onMenuClick">
          <a-menu-item key="crack">
            <template #icon>
              <icon-lock />
            </template>
            {{ $t('common.serviceCrack') }}
          </a-menu-item>
          <a-menu-item key="port">
            <template #icon>
              <icon-scan />
            </template>
            {{ $t('common.portScan') }}
          </a-menu-item>
          <a-menu-item key="vuln">
            <template #icon>
              <icon-bug />
            </template>
            {{ $t('common.vulnerabilityScan') }}
          </a-menu-item>
          <a-menu-item key="config">
            <template #icon>
              <icon-edit />
            </template>
            {{ $t('common.config') }}
          </a-menu-item>
        </a-menu>

        <div class="sidebar-footer">
          <a-tooltip :content="collapsed ? $t('common.expand') : $t('common.collapse')" position="right">
            <a-button shape="circle" class="collapse-btn" @click="toggleCollapse" :disabled="darkTheme">
              <template #icon>
                <icon-left v-if="!collapsed" />
                <icon-right v-else />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip :content="darkTheme ? $t('common.lightMode') : $t('common.darkMode')" position="right">
            <a-button shape="circle" class="theme-toggle-btn" @click="toggleTheme">
              <template #icon>
                <icon-moon-fill v-if="darkTheme" />
                <icon-sun-fill v-else />
              </template>
            </a-button>
          </a-tooltip>
        </div>
      </a-layout-sider>

      <a-layout :style="{
        marginLeft: collapsed ? '48px' : '220px',
        transition: 'margin-left 0.2s cubic-bezier(0.34, 0.69, 0.1, 1)'
      }">
        <a-layout-header class="header">
          <div class="header-left">
            <a-breadcrumb>
              <a-breadcrumb-item>{{ $t('common.securityTools') }}</a-breadcrumb-item>
              <a-breadcrumb-item>{{ currentTitle }}</a-breadcrumb-item>
            </a-breadcrumb>
          </div>
          <div class="header-right">
            <a-button class="header-btn" type="text" @click="toggleLocale">
              {{ currentLocale === 'zh' ? '中文' : 'EN' }}
            </a-button>
            <a-tooltip :content="$t('common.refresh')">
              <a-button class="header-btn" type="text" @click="reloadApp" :style="{ marginRight: '10px' }">
                <template #icon>
                  <icon-refresh />
                </template>
              </a-button>
            </a-tooltip>
          </div>
        </a-layout-header>

        <a-layout-content class="content">
          <component :is="currentComponent" />
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-config-provider>
</template>

<script lang="ts" setup>
import { ref, computed, reactive, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import type { Component } from 'vue';
import PortScanner from './components/PortScanner.vue';
import VulnScanner from './components/VulnScanner.vue';
import ServiceCrack from './components/ServiceCrack.vue';
import {
  IconScan,
  IconBug,
  IconLock,
  IconEdit,
  IconLeft,
  IconRight,
  IconMoonFill,
  IconSunFill,
  IconUser,
  IconRefresh,
} from '@arco-design/web-vue/es/icon';
import enUS from '@arco-design/web-vue/es/locale/lang/en-us';
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import Config from './components/Config.vue';
import { WindowReload } from '../wailsjs/runtime/runtime';

const { t, locale } = useI18n();

const currentLocale = ref(locale.value);

const reloadApp = () => {
  WindowReload()
};

const arcoLocale = computed(() => {
  return currentLocale.value === 'zh' ? zhCN : enUS;
});

const toggleLocale = () => {
  currentLocale.value = currentLocale.value === 'zh' ? 'en' : 'zh';
  locale.value = currentLocale.value;
  localStorage.setItem('locale', currentLocale.value);
};

onMounted(() => {
  const savedLocale = localStorage.getItem('locale');
  if (savedLocale) {
    currentLocale.value = savedLocale;
    locale.value = savedLocale;
  }
});

const selectedKeys = ref(['crack']);
const collapsed = ref(false);
const darkTheme = ref(false);

const componentMap: Record<string, Component> = {
  crack: ServiceCrack,
  port: PortScanner,
  vuln: VulnScanner,
  config: Config,
};

const titleMap = computed<Record<string, string>>(() => ({
  crack: t('common.serviceCrack'),
  port: t('common.portScan'),
  vuln: t('common.vulnerabilityScan'),
  config: t('common.config')
}));

const currentComponent = computed(() => componentMap[selectedKeys.value[0]])
const currentTitle = computed(() => computed(() => titleMap.value[selectedKeys.value[0]]))

const themeConfig = reactive({
  components: {
    Layout: {
      siderBorder: 'none'
    },
    Menu: {
      itemSelectedBg: 'var(--color-primary-light-1)',
      itemHoverBg: 'var(--color-primary-light-2)',
    }
  },
  dark: darkTheme
});

function onMenuClick(key: string) {
  selectedKeys.value = [key];
}

function toggleCollapse() {
  if (!darkTheme.value) {
    collapsed.value = !collapsed.value;
  }
}

function onCollapse(val: boolean) {
  collapsed.value = val;
}

function toggleTheme() {
  darkTheme.value = !darkTheme.value;

  if (darkTheme.value) {
    collapsed.value = true;
  }

  if (darkTheme.value) {
    document.body.setAttribute('arco-theme', 'dark');
  } else {
    document.body.removeAttribute('arco-theme');
  }
}
</script>

<style scoped>
.logo-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  padding: 12px;
  box-sizing: border-box;
  transition: all 0.2s ease;
}

.logo {
  height: 32px;
  width: auto;
  transition: all 0.2s ease;
}

.logo-container.collapsed-logo {
  padding: 12px 6px;
}

.logo-container.collapsed-logo .logo {
  height: 28px;
  transform: scale(0.5);
}

.sidebar-footer {
  position: absolute;
  bottom: 20px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.collapse-btn,
.theme-toggle-btn {
  background-color: var(--color-bg-2);
  border: 1px solid var(--color-border);
  color: var(--color-text-1);
}

.collapse-btn:hover,
.theme-toggle-btn:hover {
  background-color: var(--color-fill-2);
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 24px;
  background-color: var(--color-bg-2);
  border-bottom: 1px solid var(--color-border);
}

.header-right {
  display: flex;
  align-items: center;
  height: 100%;
}

.header-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 32px;
  margin-right: 8px;
  padding: 0 8px;
}

.header-btn .arco-icon {
  font-size: 16px;
  vertical-align: middle;
}

.content {
  padding: 24px;
  background-color: var(--color-bg-1);
  min-height: calc(100vh - 64px);
  overflow: auto;
}
</style>