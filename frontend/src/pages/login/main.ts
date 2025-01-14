import '@/assets/main.css';

import {createApp} from 'vue';
import LoginView from '@/pages/login/LoginView.vue';

// Vuetify
import 'vuetify/styles';
import {createVuetify} from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

const vuetify = createVuetify({
  components,
  directives,
});

const app = createApp(LoginView);

app.use(vuetify).mount('#app');
