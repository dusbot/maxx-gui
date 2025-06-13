<script lang="ts" setup>
import { ref } from 'vue';

interface BruteResult {
  username: string;
  password: string;
  status: string;
}

const target = ref('');
const service = ref('');
const results = ref<BruteResult[]>([]);
const columns = [
  { title: '用户名', dataIndex: 'username' },
  { title: '密码', dataIndex: 'password' },
  { title: '状态', dataIndex: 'status' },
];

function onBrute() {
  results.value = [
    { username: 'admin', password: '123456', status: '成功' },
    { username: 'root', password: 'toor', status: '失败' },
  ];
}
</script>