<script lang="ts" setup>
import { ref } from 'vue';

interface PortResult {
  port: number;
  status: string;
}

const target = ref('');
const results = ref<PortResult[]>([]);
const columns = [
  { title: '端口', dataIndex: 'port' },
  { title: '状态', dataIndex: 'status' },
];

function onScan() {
  // 这里调用后端API进行扫描
  results.value = [
    { port: 80, status: 'open' },
    { port: 22, status: 'closed' },
  ];
}
</script>