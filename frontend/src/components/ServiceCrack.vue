<template>
  <div class="service-crack-container">
    <a-card :title="$t('common.scanConfiguration')" class="config-card">
      <a-form :model="formState" layout="vertical">
        <a-row :gutter="24">
          <a-col :span="10">
            <a-form-item>
              <template #label>
                <label-component>
                  <span style="font-weight: bold">{{
                    $t("common.targetAddress")
                  }}</span>
                  <span
                    @click="openUploadTargetDialog"
                    style="cursor: pointer; color: blue"
                  >
                    点击上传
                  </span>
                </label-component>
              </template>
              <a-textarea
                v-model="formState.target"
                :placeholder="$t('placeholder.targetAddress')"
                :auto-size="{ minRows: 5, maxRows: 5 }"
              />
            </a-form-item>
          </a-col>
          <a-col :span="7">
            <a-form-item>
              <template #label>
                <label-component>
                  <span style="font-weight: bold">{{
                    $t("common.username")
                  }}</span>
                  <span
                    @click="openUploadUsernameDialog"
                    style="cursor: pointer; color: blue"
                  >
                    点击上传
                  </span>
                </label-component>
              </template>
              <a-textarea
                v-model="formState.username"
                :placeholder="$t('placeholder.username')"
                :auto-size="{ minRows: 5, maxRows: 5 }"
              />
            </a-form-item>
          </a-col>
          <a-col :span="7">
            <a-form-item>
              <template #label>
                <label-component>
                  <span style="font-weight: bold">{{
                    $t("common.password")
                  }}</span>
                  <span
                    @click="openUploadUsernameDialog"
                    style="cursor: pointer; color: blue"
                  >
                    点击上传
                  </span>
                </label-component>
              </template>
              <a-textarea
                v-model="formState.password"
                :placeholder="$t('placeholder.password')"
                :auto-size="{ minRows: 5, maxRows: 5 }"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="4">
            <a-form-item :label="$t('common.threads')">
              <a-input-number
                v-model="formState.threads"
                :min="1"
                :max="4096"
              />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item :label="$t('common.interval')">
              <a-input-number
                v-model="formState.interval"
                :min="0"
                :max="10000"
              />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item :label="$t('common.maxRuntime')">
              <a-input-number v-model="formState.maxRuntime" :min="0" />
            </a-form-item>
          </a-col>
          <a-col :span="4">
            <a-form-item :label="$t('common.proxies')">
              <a-select v-model="formState.proxies" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <div
      class="action-buttons"
      style="display: flex; justify-content: space-between; width: 100%"
    >
      <div class="left-buttons" style="display: flex; gap: 8px">
        <a-button
          type="outline"
          :status="isScanning ? 'danger' : 'normal'"
          :disabled="isPaused"
          @click="handleScan"
          :loading="scanLoading"
        >
          {{ isScanning ? $t("common.cancel") : $t("common.scan") }}
        </a-button>
        <a-button
          type="outline"
          @click="handlePause"
          :disabled="!(isScanning || isPaused)"
          :status="isScanning ? 'warning' : 'success'"
        >
          {{ isPaused ? $t("common.continue") : $t("common.pause") }}
        </a-button>

        <a-button type="outline" @click="showLogDrawer = true">
          <template #icon><icon-file /></template>
          {{ $t("common.log") }}
        </a-button>
      </div>

      <div class="right-buttons" style="display: flex; gap: 8px">
        <a-button type="outline" @click="handleDownloadCSV">
          <template #icon><icon-download /></template>
          {{ $t("common.downloadCSV") }}
        </a-button>
        <a-button type="outline" @click="handleGenerateReport">
          <template #icon><icon-file-pdf /></template>
          {{ $t("common.generateReport") }}
        </a-button>
      </div>
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

    <a-modal
      :visible="uploadTargetModalVisible"
      title="上传扫描目标"
      @cancel="uploadTargetModalVisible = false"
      @before-ok="handleTargetUploadOk"
    >
      <a-upload
        draggable
        :file-list="targetFileList"
        action="https://www.mocky.io/v2/5cc8019d300000980a055e76"
      >
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">选择文件</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>已选择文件：</p>
        <ul>
          <li v-for="file in targetFileList" :key="file.uid">
            {{ file.name }}
          </li>
        </ul>
      </div>
    </a-modal>

    <a-modal
      :visible="uploadUsernameModalVisible"
      title="上传用户字典"
      @cancel="uploadUsernameModalVisible = false"
      @before-ok="handleUsernameUploadOk"
    >
      <a-upload
        draggable
        :file-list="usernameFileList"
        action="https://www.mocky.io/v2/5cc8019d300000980a055e76"
      >
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">选择文件</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>已选择文件：</p>
        <ul>
          <li v-for="file in usernameFileList" :key="file.uid">
            {{ file.name }}
          </li>
        </ul>
      </div>
    </a-modal>
    <a-modal
      :visible="uploadPasswordModalVisible"
      title="上传密码字典"
      @cancel="uploadPasswordModalVisible = false"
      @before-ok="handlePasswordUploadOk"
    >
      <a-upload
        draggable
        :file-list="passwordFileList"
        action="https://www.mocky.io/v2/5cc8019d300000980a055e76"
      >
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">选择文件</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>已选择文件：</p>
        <ul>
          <li v-for="file in passwordFileList" :key="file.uid">
            {{ file.name }}
          </li>
        </ul>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from "vue";
import { Message } from "@arco-design/web-vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

interface ScanParam {
  target: string;
  username: string;
  password: string;
  threads: number;
  interval: number;
  maxRuntime: number;
  proxies: string;
}

interface ResultItem {
  target: string;
  service: string;
  auth: string;
  extrainfo: string;
}

interface LogItem {
  time: string;
  message: string;
}

const formState = reactive<ScanParam>({
  target: "",
  username: "",
  password: "",
  threads: 1024,
  interval: 100,
  maxRuntime: 0,
  proxies: "",
});

const isScanning = ref(false);
const isPaused = ref(false);
const scanLoading = ref(false);
const tableLoading = ref(false);

const progress = ref(0);
const progressStatus = ref<"normal" | "success" | "warning" | "danger">(
  "normal"
);

const results = ref<ResultItem[]>([
  {
    target: "godzilla://10.1.1.1/sdf/asdf/as/asdf/sadf/sadf/x.php",
    service: "GODZILLA",
    auth: "root:123456",
    extrainfo: "No Auth",
  },
  {
    target: "ssh://10.1.1.1:22",
    service: "ssh",
    auth: "root:123456",
    extrainfo: "No Auth",
  },
]);
const columns = computed(() => [
  { title: t("common.target"), dataIndex: "target", width: "40%" },
  { title: t("common.service"), dataIndex: "service", width: "10%" },
  { title: t("common.auth"), dataIndex: "auth", width: "20%" },
  { title: t("common.extra_info"), dataIndex: "extrainfo", width: "30%" },
]);

const logs = ref<LogItem[]>([]);
const showLogDrawer = ref(false);

const addLog = (message: string) => {
  logs.value.push({
    time: new Date().toLocaleTimeString(),
    message,
  });
};

const handleGenerateReport = () => {
  console.log("生成报告");
};

const handleDownloadCSV = () => {
  console.log("下载CSV");
};

let scanInterval: number;
let scanTimer: number;

const handleScan = () => {
  if (isScanning.value) {
    isScanning.value = false;
  } else {
    isScanning.value = true;
  }
};

const handlePause = () => {
  if (isScanning.value) {
    if (!isPaused.value) {
      isPaused.value = true;
      isScanning.value = false;
    }
  } else {
    if (isPaused.value) {
      isPaused.value = false;
      isScanning.value = true;
    }
  }
  if (isPaused.value) {
    addLog(t("log.scanPaused"));
    Message.info(t("message.scanPaused"));
  } else {
    addLog(t("log.scanResumed"));
    Message.info(t("message.scanResumed"));
  }
};

const uploadTargetModalVisible = ref<boolean>(false);
const uploadUsernameModalVisible = ref<boolean>(false);
const uploadPasswordModalVisible = ref<boolean>(false);

const targetFileList = ref<any[]>([]);
const usernameFileList = ref<any[]>([]);
const passwordFileList = ref<any[]>([]);

const openUploadTargetDialog = () => {
  uploadTargetModalVisible.value = true;
};
const openUploadUsernameDialog = () => {
  uploadUsernameModalVisible.value = true;
};
const openUploadPasswordDialog = () => {
  uploadPasswordModalVisible.value = true;
};

const handleTargetUploadOk = () => {
  if (targetFileList.value.length === 0) {
    Message.warning("请至少选择一个文件");
    return false;
  }
  Message.success("扫描目标文件上传成功");
  uploadTargetModalVisible.value = false;
  targetFileList.value = [];
  return true;
};
const handleUsernameUploadOk = () => {
  if (usernameFileList.value.length === 0) {
    Message.warning("请至少选择一个文件");
    return false;
  }
  Message.success("用户文件上传成功");
  uploadUsernameModalVisible.value = false;
  usernameFileList.value = [];
  return true;
};
const handlePasswordUploadOk = () => {
  if (passwordFileList.value.length === 0) {
    Message.warning("请至少选择一个文件");
    return false;
  }
  Message.success("密码文件上传成功");
  uploadPasswordModalVisible.value = false;
  passwordFileList.value = [];
  return true;
};
</script>

<style scoped>
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