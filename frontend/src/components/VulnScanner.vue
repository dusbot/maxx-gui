<script lang="ts" setup>
import { ref } from 'vue';

interface VulnResult {
  vuln: string;
  level: string;
}

const target = ref('');
const results = ref<VulnResult[]>([]);
const columns = [
  { title: '漏洞名称', dataIndex: 'vuln' },
  { title: '风险等级', dataIndex: 'level' },
];

function onScan() {
  results.value = [
    { vuln: 'CVE-2023-1234', level: '高' },
    { vuln: 'CVE-2022-5678', level: '中' },
  ];
}
</script>