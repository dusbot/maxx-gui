import { createI18n } from 'vue-i18n';
import en from '@/locales/en';
import zh from '@/locales/zh';

const messages = {
  en,
  zh
};

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'zh',
  messages
});

export default i18n;