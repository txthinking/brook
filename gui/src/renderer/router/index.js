import Vue from 'vue'
import Router from 'vue-router'
import Server from '@/components/Server'
import Builtin from '@/components/Builtin'
import Help from '@/components/Help'
import About from '@/components/About'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Server',
      component: Server,
    },
    {
      path: '/builtin',
      name: 'Builtin',
      component: Builtin,
    },
    {
      path: '/help',
      name: 'Help',
      component: Help,
    },
    {
      path: '/about',
      name: 'About',
      component: About,
    },
  ]
})
