import '@/assets/main.css';

import {createApp} from 'vue';
import AdminView from '@/pages/admin/AdminView.vue';

// Vuetify
import 'vuetify/styles';
import {createVuetify} from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

const vuetify = createVuetify({
  components,
  directives,
});

const app = createApp(AdminView);

app.use(vuetify).mount('#app');
