import Vue from 'vue'
import App from './App'
import router from './router'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import VueSuperagent from 'vue-superagent'

Vue.use(Vuetify, { theme: {
  primary: '#000000',
  secondary: '#919191',
  accent: '#000000',
  error: '#000000',
  info: '#000000',
  success: '#000000',
  warning: '#000000'
}})
if (!process.env.IS_WEB) Vue.use(require('vue-electron'))
Vue.use(VueSuperagent, {
})
Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})

window.getSetting = function (){
    var o = {
        Type: 'Brook',
        Server: '',
        Password: '',
        TCPTimeout: 60,
        TCPDeadline: 0,
        UDPDeadline: 60,
        UDPSessionTime: 60,
    };
    var s = localStorage.getItem('brook/server');
    if (s){
        o = JSON.parse(s);
    }
    var o1 = {
        Address: '127.0.0.1:1080',
        Mode: 'pac',
        DomainURL: 'https://www.txthinking.com/pac/white.list',
        CidrURL: 'https://www.txthinking.com/pac/white_cidr.list',
        PacURL: 'https://www.txthinking.com/pac/white.pac',
    };
    s = localStorage.getItem('brook/mode');
    if (s){
        o1 = JSON.parse(s);
    }
    var o2 = {
        UseWhiteTrayIcon: false,
    };
    s = localStorage.getItem('brook/builtin');
    if (s){
        o2 = JSON.parse(s);
    }
    return Object.assign(Object.assign(o, o1), o2);
}
