<template>
  <a-config-provider :theme="themeConfig">
    <a-layout style="height: 100vh">
      <a-layout-sider :style="{
        height: '100vh',
        boxShadow: '2px 0 8px 0 rgba(29, 35, 41, 0.05)',
        position: 'fixed',
        left: 0,
        top: 0,
        bottom: 0,
        zIndex: 100,
      }" :collapsed="collapsed" :width="190" breakpoint="lg" @collapse="onCollapse">
        <div class="logo-container">
          <img src="@/assets/images/maxx_logo.svg" alt="logo" class="logo" />
        </div>
        <a-menu v-model:selected-keys="selectedKeys" :collapsed="collapsed" :style="{ borderRight: 'none' }"
          @menu-item-click="onMenuClick">
          <a-menu-item key="crack">
            <template #icon>
              <icon-lock />
            </template>
            Service Crack
          </a-menu-item>
          <a-menu-item key="port">
            <template #icon>
              <icon-scan />
            </template>
            Port Scan
          </a-menu-item>
          <a-menu-item key="vuln">
            <template #icon>
              <icon-bug />
            </template>
            Vulnerability Scan
          </a-menu-item>
        </a-menu>

        <div class="sidebar-footer">
          <a-tooltip :content="collapsed ? 'Expand' : 'Collapse'" position="right">
            <a-button shape="circle" class="collapse-btn" @click="toggleCollapse">
              <template #icon>
                <icon-left v-if="!collapsed" />
                <icon-right v-else />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip :content="darkTheme ? 'Light Mode' : 'Dark Mode'" position="right">
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
              <a-breadcrumb-item>Security Tools</a-breadcrumb-item>
              <a-breadcrumb-item>{{ currentTitle }}</a-breadcrumb-item>
            </a-breadcrumb>
          </div>
          <div class="header-right">
            <a-avatar :size="32" :style="{ backgroundColor: '#3370ff' }">
              <icon-user />
            </a-avatar>
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
import { ref, computed, reactive } from 'vue';
import type { Component } from 'vue';
import PortScanner from './components/PortScanner.vue';
import VulnScanner from './components/VulnScanner.vue';
import ServiceCrack from './components/ServiceCrack.vue';
import {
  IconScan,
  IconBug,
  IconLock,
  IconLeft,
  IconRight,
  IconMoonFill,
  IconSunFill,
  IconUser
} from '@arco-design/web-vue/es/icon';

const selectedKeys = ref(['crack']);
const collapsed = ref(false);
const darkTheme = ref(false);

const componentMap: Record<string, Component> = {
  crack: ServiceCrack,
  port: PortScanner,
  vuln: VulnScanner,
};

const titleMap: Record<string, string> = {
  crack: 'Service Crack',
  port: 'Port Scan',
  vuln: 'Vulnerability Scan',
};

const currentComponent = computed(() => componentMap[selectedKeys.value[0]]);
const currentTitle = computed(() => titleMap[selectedKeys.value[0]]);

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
  collapsed.value = !collapsed.value;
}

function onCollapse(val: boolean) {
  collapsed.value = val;
}

function toggleTheme() {
  darkTheme.value = !darkTheme.value;

  if (darkTheme.value) {
    document.body.setAttribute('arco-theme', 'dark');
    themeConfig.components.Menu.itemSelectedBg = 'var(--color-primary-light-1)';
    themeConfig.components.Menu.itemHoverBg = 'var(--color-primary-light-2)';
  } else {
    document.body.removeAttribute('arco-theme');
    themeConfig.components.Menu.itemSelectedBg = 'var(--color-primary-light-1)';
    themeConfig.components.Menu.itemHoverBg = 'var(--color-primary-light-2)';
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
}

.logo {
  height: 32px;
  transition: all 0.2s;
}

.logo-collapsed {
  height: 28px;
  transition: all 0.2s;
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

.content {
  padding: 24px;
  background-color: var(--color-bg-1);
  min-height: calc(100vh - 64px);
  overflow: auto;
}
</style>