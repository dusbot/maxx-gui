<template>
  <div class="service-crack-container">
    <a-collapse :active-key="collapseActiveKeys" @change="handleCollapseChange">
      <a-collapse-item :header="$t('common.scanConfigurationCollapse')" key="1">
        <a-card :title="$t('common.scanConfiguration')" class="config-card">
          <a-form :model="scanParam" layout="vertical">
            <a-row :gutter="24">
              <a-col :span="10">
                <a-form-item>
                  <template #label>
                    <label-component>
                      <span style="font-weight: bold; margin-right: 10px">{{
                        $t("common.targetAddress")
                        }}</span>
                      <span @click="openUploadTargetDialog" style="cursor: pointer; color: #1890ff">
                        {{ t("common.click2Upload") }}
                      </span>
                    </label-component>
                  </template>
                  <a-textarea v-model="scanParam.target" :placeholder="$t('placeholder.targetAddress')"
                    :auto-size="{ minRows: 5, maxRows: 5 }" :error="isTargetInvalid"
                    @blur="isTargetInvalid = !validateTarget(scanParam.target)" @input="isTargetInvalid = false"
                    :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item>
                  <template #label>
                    <label-component>
                      <span style="font-weight: bold; margin-right: 10px">{{
                        $t("common.username")
                        }}</span>
                      <span @click="openUploadUsernameDialog" style="cursor: pointer; color: #1890ff">
                        {{ t("common.click2Upload") }}
                      </span>
                    </label-component>
                  </template>
                  <a-textarea v-model="scanParam.username" :placeholder="$t('placeholder.username')"
                    :auto-size="{ minRows: 5, maxRows: 5 }" :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
              <a-col :span="7">
                <a-form-item>
                  <template #label>
                    <label-component>
                      <span style="font-weight: bold; margin-right: 10px">{{
                        $t("common.password")
                        }}</span>
                      <span @click="openUploadPasswordDialog" style="cursor: pointer; color: #1890ff">
                        {{ t("common.click2Upload") }}
                      </span>
                    </label-component>
                  </template>
                  <a-textarea v-model="scanParam.password" :placeholder="$t('placeholder.password')"
                    :auto-size="{ minRows: 5, maxRows: 5 }" :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-row :gutter="24">
              <a-col :span="4">
                <a-form-item :label="$t('common.threads')">
                  <a-input-number v-model="scanParam.threads" :min="1" :max="4096" :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
              <a-col :span="4">
                <a-form-item :label="$t('common.interval')">
                  <a-input-number v-model="scanParam.interval" :min="0" :max="3600000"
                    :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
              <a-col :span="4">
                <a-form-item :label="$t('common.maxRuntime')">
                  <a-input-number v-model="scanParam.maxRuntime" :min="0" :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
              <a-col :span="4">
                <a-form-item :label="$t('common.proxies')">
                  <a-select v-model="scanParam.proxies" :disabled="isScanning || isPaused" />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-card>
      </a-collapse-item>
    </a-collapse>
    <div class="action-buttons" style="display: flex; justify-content: space-between; width: 100%">
      <div class="left-buttons" style="display: flex; gap: 8px">
        <a-button type="outline" :status="isScanning ? 'danger' : 'normal'" :disabled="isPaused" @click="handleScan"
          :loading="scanLoading">
          {{ isScanning ? $t("common.cancel") : $t("common.scan") }}
        </a-button>
        <a-button type="outline" @click="handlePause" :disabled="!(isScanning || isPaused)"
          :status="isScanning ? 'warning' : 'normal'">
          {{ isPaused ? $t("common.continue") : $t("common.pause") }}
        </a-button>

        <a-button type="outline" @click="showLogDrawer = true">
          {{ $t("common.log") }}
        </a-button>
        <a-button type="outline" @click="results = []">
          {{ $t("common.clearData") }}
        </a-button>
      </div>

      <div class="right-buttons" style="display: flex; gap: 8px">
        <a-button type="outline" @click="handleDownloadCSV">
          <template #icon><icon-download /></template>
          {{ $t("common.downloadCSV") }}
        </a-button>
        <a-button type="outline" @click="handleGenerateReport" :loading="scanLoading">
          {{ t('common.generateReport') }}
        </a-button>
      </div>
    </div>

    <a-progress :percent="progress" :status="progressStatus" size="large" class="progress-bar" :show-text="true"
      :animation="true" />

    <a-card :title="$t('common.results')" class="results-card">
      <a-table :columns="columns" :data="displayData" :pagination="{
        current: pagination.current,
        pageSize: pagination.pageSize,
        total: results.length,
        showTotal: true,
        showPageSize: true,
        pageSizeOptions: [10, 20, 50]
      }" :loading="tableLoading" @page-change="handlePageChange" @page-size-change="handlePageSizeChange">
        <template #password="{ record }">
          <a-tooltip :content="record.password">
            <span class="password-cell">{{ record.password }}</span>
          </a-tooltip>
        </template>
      </a-table>
    </a-card>

    <a-drawer :title="$t('common.executionLog')" :visible="showLogDrawer" @cancel="showLogDrawer = false" :width="600">
      <div class="log-content">
        <a-timeline>
          <a-timeline-item v-for="(log, index) in logs" :key="index">
            {{ log.time }} - {{ log.message }}
          </a-timeline-item>
        </a-timeline>
      </div>
    </a-drawer>

    <a-modal :visible="uploadTargetModalVisible" :title="t('common.uploadTarget')"
      @cancel="uploadTargetModalVisible = false" @before-ok="handleTargetUploadOk">
      <a-upload draggable :file-list="targetFileList" :action="''" :auto-upload="false">
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">{{ t("common.selectFile") }}</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>{{ t("common.selectedFile") }}</p>
        <ul>
          <li v-for="file in targetFileList" :key="file.uid">
            {{ file.name }}
          </li>
        </ul>
      </div>
    </a-modal>

    <a-modal :visible="uploadUsernameModalVisible" :title="t('common.uploadUsername')"
      @cancel="uploadUsernameModalVisible = false" @before-ok="handleUsernameUploadOk">
      <a-upload draggable :file-list="usernameFileList" action="https://www.mocky.io/v2/5cc8019d300000980a055e76">
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">{{ t("common.selectFile") }}</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>{{ t("common.selectedFile") }}</p>
        <ul>
          <li v-for="file in usernameFileList" :key="file.uid">
            {{ file.name }}
          </li>
        </ul>
      </div>
    </a-modal>
    <a-modal :visible="uploadPasswordModalVisible" :title="t('common.uploadPassword')"
      @cancel="uploadPasswordModalVisible = false" @before-ok="handlePasswordUploadOk">
      <a-upload draggable :file-list="passwordFileList" action="https://www.mocky.io/v2/5cc8019d300000980a055e76">
        :multiple="true" :limit="5" accept="text/plain" >
        <a-button type="primary">{{ t("common.selectFile") }}</a-button>
      </a-upload>
      <div style="margin-top: 16px">
        <p>{{ t("common.selectedFile") }}</p>
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
import { parseIpPortRange } from "@/utils/target";
import { addUniqueItem } from "@/utils/utils";
import { Message } from "@arco-design/web-vue";
import { computed, reactive, ref } from "vue";
import { useI18n } from "vue-i18n";
import { CancelAll, GenerateReport, Scan } from "../../wailsjs/go/handler/CrackHandler";
import { consts, model } from "../../wailsjs/go/models";
import { EventsOn } from "../../wailsjs/runtime/runtime";

const { t, locale } = useI18n();

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
  id: string
  target: string;
  service: string;
  auth: string;
  extrainfo: string;
}

interface LogItem {
  time: string;
  message: string;
}

interface ReportData {
  summary: {
    totalTargets: number;
    services: Record<string, number>;
    authTypes: Record<string, number>;
    riskLevels: {
      high: number;
      medium: number;
      low: number;
    };
  };
  details: ResultItem[];
  riskAssessments: Record<string, {
    level: 'high' | 'medium' | 'low';
    reasons: string[];
  }>;
}

const reportGenerated = ref(false);
const reportData = ref<Blob | null>(null);
const reportFileName = ref('');

const scanParam = reactive<ScanParam>({
  target: "",
  username: "",
  password: "",
  threads: 1024,
  interval: 100,
  maxRuntime: 0,
  proxies: "",
});

const isTargetInvalid = ref(false);
const isScanning = ref(false);
const isPaused = ref(false);
const scanLoading = ref(false);
const tableLoading = ref(false);

const isUserToggled = ref(false);
const collapseActiveKeys = ref<string[]>(['1']);

const progress = ref(0);
const progressStatus = ref<"normal" | "success" | "warning" | "danger">(
  "normal"
);

const pagination = ref({
  current: 1,
  pageSize: 10,
});

const results = ref<ResultItem[]>([]);

const displayData = computed(() => {
  const { current, pageSize } = pagination.value;
  const start = (current - 1) * pageSize;
  const end = start + pageSize;
  return results.value.slice(start, end);
});

const handlePageChange = (page: number) => {
  pagination.value.current = page;
};

const handlePageSizeChange = (size: number) => {
  pagination.value.pageSize = size;
  pagination.value.current = 1;
};

const validateTarget = (value: string) => {
  const lines = value.trim().split('\n').filter(line => line.trim() !== '');
  for (const line of lines) {
    if (!isValidTarget(line)) {
      return false;
    }
  }
  return true;
};

const isValidTarget = (target: string) => {
  // protocol://ip:port (e.g. ftp://192.168.1.1:21)
  const protocolRegex = /^[a-z]+:\/\/\d{1,3}(\.\d{1,3}){3}(:\d{1,5})?$/i;

  // ip:port (e.g. 192.168.1.1:22)
  const ipPortRegex = /^\d{1,3}(\.\d{1,3}){3}:\d{1,5}$/;

  // ip/mask:[portrange] (e.g. 192.168.1.0/24:[80,8000-9000])
  const ipRangeRegex = /^\d{1,3}(\.\d{1,3}){3}\/\d{1,2}:$$(\d{1,5}(-\d{1,5})?)(\|(\d{1,5}(-\d{1,5})?))*$$$/;

  return protocolRegex.test(target) ||
    ipPortRegex.test(target) ||
    ipRangeRegex.test(target) || parseIpPortRange(target)
};

EventsOn(consts.EVENT.EVENT_RESULT, (data: any) => {
  const result = data as model.CrackResult;
  let auth = result.Username + ":" + result.Password;
  let extrainfo = "";
  if (auth === ":") {
    auth = "";
    extrainfo = "No Auth";
  }
  const newItem = {
    id: result.TaskID,
    target: result.Target,
    service: result.Service,
    auth: auth,
    extrainfo: extrainfo,
  };
  addUniqueItem(results.value, newItem, [
    "id",
    "target",
    "service",
    "auth",
    "extrainfo",
  ]);
});

EventsOn(consts.EVENT.EVENT_PROGRESS, (data: any) => {
  const currProgress = data as number;
  progress.value = currProgress;
  if (currProgress >= 1) {
    isScanning.value = false;
    Message.info(t("message.scanCompleted"));
  }
});

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

const handleCollapseChange = (keys: string[]) => {
  if (isScanning.value) {
    isUserToggled.value = true;
  }
  collapseActiveKeys.value = keys
};

const handleGenerateReport = () => {
  if (results.value.length === 0) {
    Message.warning(t('message.noDataToExport'));
    return;
  }
  GenerateReport(results.value[0].id, locale.value.startsWith('zh')).then((content) => {
    if (content === "") {
      Message.warning(t('message.reportGenerateFailed'));
      return;
    }
    const blob = new Blob([content], { type: 'text/html;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${results.value[0].id}.html`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  })
  reportGenerated.value = true;
};

const handleDownloadCSV = () => {
  if (results.value.length === 0) {
    Message.warning(t('message.noDataToExport'));
    return;
  }
  const headers = ['Target', 'Service', 'Authentication', 'Extra Info', 'Risk Level'];
  const csvRows = results.value.map(item => {
    const risk = stubAssestRisk(item);
    return [
      `"${item.target}"`,
      `"${item.service}"`,
      `"${item.auth || 'None'}"`,
      `"${item.extrainfo || ''}"`,
      `"${risk.level}"`
    ].join(',');
  });

  const csvContent = [
    headers.join(','),
    ...csvRows
  ].join('\n');

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `scan_results_${new Date().toISOString().slice(0, 10)}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

// It will be ported to the backend in the future, 
// and more comprehensive judgment rules will be embedded.
const stubAssestRisk = (item: ResultItem) => {
  const reasons: string[] = [];
  let level: 'high' | 'medium' | 'low' = 'low';
  if (['ssh', 'telnet', 'rdp', 'ftp',].includes(item.service.toLowerCase())) {
    reasons.push(`Disposed ${item.service} service`);
    level = 'medium';
  }
  if (item.auth) {
    if (item.auth.includes('admin') || item.auth.includes('root')) {
      reasons.push('Using admin or root account');
      level = 'high';
    }
    if (item.auth.includes('123456') || item.auth.includes('password')) {
      reasons.push('Using weak password');
      level = 'high';
    }
  } else if (item.service.toLowerCase() === 'http') {
    reasons.push('HTTP service has no authentication');
    level = 'medium';
  }
  return { level, reasons };
};

const handleScan = () => {
  reportGenerated.value = false
  if (isScanning.value) {
    CancelAll().then((ok: boolean) => {
      if (ok) {
        isScanning.value = false;
        Message.success(t("message.scanCancelled"));
      } else {
        Message.error(t("message.cancelFailed"));
      }
    })
  } else {
    progress.value = 0;
    if (scanParam.target === "") {
      Message.warning(t("message.targetAddressRequirement"));
      isTargetInvalid.value = true;
      return;
    }
    if (!validateTarget(scanParam.target)) {
      isTargetInvalid.value = true;
      Message.error(t('message.invalidTargetFormat'));
      return;
    }
    isTargetInvalid.value = false;
    isTargetInvalid.value = false;
    const task = new model.CrackTask();
    task.Targets = scanParam.target;
    task.Usernames = scanParam.username;
    task.Passwords = scanParam.password;
    task.MaxRuntime = scanParam.maxRuntime;
    task.Thread = scanParam.threads;
    task.Proxies = scanParam.proxies;
    task.Interval = scanParam.interval;
    Scan(task).then((ok) => {
      isScanning.value = ok
      if (ok) {
        collapseActiveKeys.value = [];
      }
    }).catch((err) => {
      isScanning.value = false;
      Message.error(err);
    });
  }
};

const handlePause = () => {
  Message.error(t("message.notSupportedYet"));
  // if (isScanning.value) {
  //   if (!isPaused.value) {
  //     isPaused.value = true;
  //     isScanning.value = false;
  //   }
  // } else {
  //   if (isPaused.value) {
  //     isPaused.value = false;
  //     isScanning.value = true;
  //   }
  // }
  // if (isPaused.value) {
  //   addLog(t("log.scanPaused"));
  //   Message.info(t("message.scanPaused"));
  // } else {
  //   addLog(t("log.scanResumed"));
  //   Message.info(t("message.scanResumed"));
  // }
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
    Message.warning(t("message.oneFileSelectedRequirement"));
    return false;
  }
  Message.success(t("message.uploadSuccess"));
  uploadTargetModalVisible.value = false;
  targetFileList.value = [];
  return true;
};
const handleUsernameUploadOk = () => {
  if (usernameFileList.value.length === 0) {
    Message.warning(t("message.oneFileSelectedRequirement"));
    return false;
  }
  Message.success(t("message.uploadSuccess"));
  uploadUsernameModalVisible.value = false;
  usernameFileList.value = [];
  return true;
};
const handlePasswordUploadOk = () => {
  if (passwordFileList.value.length === 0) {
    Message.warning(t("message.oneFileSelectedRequirement"));
    return false;
  }
  Message.success(t("message.uploadSuccess"));
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