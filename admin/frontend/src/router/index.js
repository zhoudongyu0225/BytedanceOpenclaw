import { createRouter, createWebHashHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Layout from '../views/Layout.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: {
      noAuth: true
    }
  },
  {
    path: '/',
    name: 'layout',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: {
          title: '数据概览'
        }
      },
      {
        path: 'anchor',
        name: 'anchor',
        component: () => import('../views/AnchorList.vue'),
        meta: {
          title: '主播管理'
        }
      },
      {
        path: 'gift-simulator',
        name: 'gift-simulator',
        component: () => import('../views/GiftSimulator.vue'),
        meta: {
          title: '礼物模拟工具'
        }
      },
      {
        path: 'gm-command',
        name: 'gm-command',
        component: () => import('../views/GmCommand.vue'),
        meta: {
          title: 'GM命令'
        }
      },
      {
        path: 'fans',
        name: 'fans',
        component: () => import('../views/Fans.vue'),
        meta: {
          title: '粉丝管理'
        }
      },
      {
        path: 'announcement',
        name: 'announcement',
        component: () => import('../views/Announcement.vue'),
        meta: {
          title: '公告管理'
        }
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('../views/Settings.vue'),
        meta: {
          title: '系统设置'
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.noAuth) {
    next()
  } else {
    if (!token) {
      next('/login')
    } else {
      next()
    }
  }
})

export default router
