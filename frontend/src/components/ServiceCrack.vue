<template>
  <div class="service-crack-container">
    <a-card :title="$t('common.scanConfiguration')" class="config-card">
      <a-form :model="formState" layout="vertical">
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-item :label="$t('common.targetAddress')">
              <a-textarea 
                v-model="formState.target" 
                :placeholder="$t('placeholder.targetAddress')"
                :auto-size="{minRows:5,maxRows:5}" 
              />
            </a-form-item>
          </a-col>
          <a-col :span="7">
            <a-form-item :label="$t('common.username')">
              <a-textarea 
                v-model="formState.username" 
                :placeholder="$t('placeholder.username')"
                :auto-size="{minRows:5,maxRows:5}"
              />
            </a-form-item>
          </a-col>
          <a-col :span="7">
            <a-form-item :label="$t('common.password')">
              <a-textarea 
                v-model="formState.password" 
                :placeholder="$t('placeholder.password')"
                :auto-size="{minRows:5,maxRows:5}"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="4">
            <a-form-item :label="$t('common.threads')">
              <a-input-number v-model="formState.threads" :min="1" :max="100" />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item :label="$t('common.interval')">
              <a-input-number v-model="formState.interval" :min="0" :max="10000" />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item :label="$t('common.maxRuntime')">
              <a-input-number v-model="formState.maxRuntime" :min="1" :max="120" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <div class="action-buttons">
      <a-button 
        type="primary" 
        @click="handleScan"
        :loading="scanLoading"
      >
        {{ isScanning ? $t('common.cancel') : $t('common.scan') }}
      </a-button>
      
      <a-button 
        type="secondary" 
        @click="handlePause"
        :disabled="!isScanning || isPaused"
      >
        {{ isPaused ? $t('common.continue') : $t('common.pause') }}
      </a-button>
      
      <a-button @click="showLogDrawer = true">
        <template #icon><icon-file /></template>
        {{ $t('common.log') }}
      </a-button>
    </div>

    <a-progress 
      :percent="progress" 
      :status="progressStatus" 
      class="progress-bar"
      :show-text="false"
    />

    <a-card :title="$t('common.results')" class="results-card">
      <a-table 
        :columns="columns" 
        :data="results" 
        :pagination="false"
        :loading="tableLoading"
      >
        <template #password="{ record }">
          <a-tooltip :content="record.password">
            <span class="password-cell">{{ record.password }}</span>
          </a-tooltip>
        </template>
      </a-table>
    </a-card>

    <a-drawer 
      :title="$t('common.executionLog')" 
      :visible="showLogDrawer" 
      @cancel="showLogDrawer = false"
      :width="600"
    >
      <div class="log-content">
        <a-timeline>
          <a-timeline-item v-for="(log, index) in logs" :key="index">
            {{ log.time }} - {{ log.message }}
          </a-timeline-item>
        </a-timeline>
      </div>
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive,computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface FormState {
  target: string;
  serviceType: string;
  username: string;
  password: string;
  threads: number;
  interval: number;
  maxRuntime: number;
  dictionaryFile: string;
}

interface ResultItem {
  target: string;
  username: string;
  password: string;
  status: string;
}

interface LogItem {
  time: string;
  message: string;
}

const formState = reactive<FormState>({
  target: '',
  serviceType: 'ssh',
  username: '',
  password: '',
  threads: 10,
  interval: 100,
  maxRuntime: 30,
  dictionaryFile: ''
});

const isScanning = ref(false);
const isPaused = ref(false);
const scanLoading = ref(false);
const tableLoading = ref(false);

const progress = ref(0);
const progressStatus = ref<'normal' | 'success' | 'warning' | 'danger'>('normal');

const results = ref<ResultItem[]>([]);
const columns = computed(() => [
  { title: t('common.target'), dataIndex: 'target' },
  { title: t('common.username'), dataIndex: 'username' },
  { 
    title: t('common.password'), 
    dataIndex: 'password',
    slotName: 'password'
  },
  { title: t('common.service'), dataIndex: 'service' }
]);

const logs = ref<LogItem[]>([]);
const showLogDrawer = ref(false);

const addLog = (message: string) => {
  logs.value.push({
    time: new Date().toLocaleTimeString(),
    message
  });
};

let scanInterval: number;
let scanTimer: number;

const handleScan = () => {
  if (isScanning.value) {
    clearInterval(scanInterval);
    clearTimeout(scanTimer);
    isScanning.value = false;
    progressStatus.value = 'danger';
    addLog(t('log.scanCancelled'));
    Message.warning(t('message.scanCancelled'));
  } else {
    if (!formState.target) {
      Message.error(t('message.enterTargetAddress'));
      return;
    }
    
    scanLoading.value = true;
    tableLoading.value = true;
    addLog(t('log.startingScan'));
    
    setTimeout(() => {
      isScanning.value = true;
      scanLoading.value = false;
      progressStatus.value = 'normal';
      addLog(t('log.scanStarted', { target: formState.target }));
      
      let currentProgress = 0;
      scanInterval = setInterval(() => {
        if (!isPaused.value) {
          currentProgress += Math.random() * 5;
          progress.value = Math.min(currentProgress, 100);
          
          if (Math.random() > 0.9) {
            const newResult: ResultItem = {
              target: formState.target,
              username: formState.username || `user${Math.floor(Math.random() * 100)}`,
              password: formState.password || `pass${Math.floor(Math.random() * 1000)}`,
              status: t('common.cracked')
            };
            results.value.push(newResult);
            addLog(t('log.foundCredentials', { username: newResult.username }));
          }
          
          if (currentProgress >= 100) {
            clearInterval(scanInterval);
            isScanning.value = false;
            progressStatus.value = 'success';
            addLog(t('log.scanCompleted'));
            Message.success(t('message.scanCompleted'));
          }
        }
      }, 500);
      
      scanTimer = setTimeout(() => {
        if (isScanning.value) {
          clearInterval(scanInterval);
          isScanning.value = false;
          progress.value = 100;
          progressStatus.value = 'warning';
          addLog(t('log.scanStopped'));
          Message.warning(t('message.scanStopped'));
        }
      }, formState.maxRuntime * 60 * 1000);
      
      tableLoading.value = false;
    }, 1000);
  }
};

const handlePause = () => {
  isPaused.value = !isPaused.value;
  if (isPaused.value) {
    addLog(t('log.scanPaused'));
    Message.info(t('message.scanPaused'));
  } else {
    addLog(t('log.scanResumed'));
    Message.info(t('message.scanResumed'));
  }
};
</script>

<style scoped>
/* 原有样式保持不变 */
.service-crack-container {
  padding: 20px;
}

.config-card {
  margin-bottom: 20px;
}

.action-buttons {
  margin: 20px 0;
  display: flex;
  gap: 12px;
}

.progress-bar {
  margin-bottom: 20px;
}

.results-card {
  margin-top: 20px;
}

.password-cell {
  display: inline-block;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.log-content {
  height: calc(100vh - 100px);
  overflow-y: auto;
  padding-right: 10px;
}
</style>