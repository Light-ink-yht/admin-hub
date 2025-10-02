<script setup lang="ts">
import { UserOutlined, LockOutlined, BulbOutlined } from '@ant-design/icons-vue';
import { ref, onMounted, onUnmounted } from 'vue';

defineOptions({
  name: 'AuthLogin',
})

// 验证码相关
const captchaCode = ref('')
const captchaSrc = ref('')

// 刷新验证码
const refreshCaptcha = () => {
  // 这里应该调用后端接口获取验证码图片
  // 示例：captchaSrc.value = '/api/captcha?t=' + new Date().getTime();
  captchaSrc.value = 'https://via.placeholder.com/100x38?text=CAPTCHA'
}

// 忘记密码
const toRePassword = () => {
  loginForm.value = false
  registerForm.value = false
  rePasswordForm.value = true
}

// 前往注册
const toRegister = () => {
  loginForm.value = false
  registerForm.value = true
  rePasswordForm.value = false
}

// 前往登录
const toLogin = () => {
  loginForm.value = true
  registerForm.value = false
  rePasswordForm.value = false
}

const loginForm = ref(true)
const registerForm = ref(false)
const rePasswordForm = ref(false)

// 粒子背景相关
type Particle = {
  x: number;
  y: number;
  radius: number;
  color: string;
  speedX: number;
  speedY: number;
};

let particles: Particle[] = []
let animationFrameId: number
let canvas2d: CanvasRenderingContext2D | null = null

// 组件挂载时获取验证码
onMounted(() => {
  refreshCaptcha();
  initParticles()
})

// 初始化粒子系统
const initParticles = () => {
  const canvas = document.createElement('canvas')
  canvas.id = 'particleCanvas'
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight
  canvas.style.position = 'fixed'
  canvas.style.top = '0'
  canvas.style.left = '0'
  canvas.style.zIndex = '-1' // 置于底层
  canvas.style.pointerEvents = 'none' // 不阻挡鼠标事件

  document.body.appendChild(canvas)

  canvas2d = canvas.getContext('2d')
  createParticles()
  animateParticles()

  // 窗口大小改变时重置画布尺寸
  window.addEventListener('resize', resizeCanvas)
}

// 创建粒子
const createParticles = () => {
  particles = []
  const particleCount = Math.floor((window.innerWidth * window.innerHeight) / 5000)

  for (let i = 0; i < particleCount; i++) {
    particles.push({
      x: Math.random() * window.innerWidth,
      y: Math.random() * window.innerHeight,
      radius: Math.random() * 2 + 1,
      color: `rgba(${Math.floor(Math.random() * 100 + 155)}, ${Math.floor(Math.random() * 100 + 155)}, ${Math.floor(Math.random() * 100 + 200)}, ${Math.random() * 0.5 + 0.3})`,
      speedX: Math.random() * 1 - 0.5,
      speedY: Math.random() * 1 - 0.5
    })
  }
}

// 重置画布尺寸
const resizeCanvas = () => {
  const canvas = document.getElementById('particleCanvas') as HTMLCanvasElement
  if (canvas) {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
    createParticles()
  }
}

// 动画粒子
const animateParticles = () => {
  if (!canvas2d) return

  const canvas = canvas2d.canvas
  canvas2d.clearRect(0, 0, canvas.width, canvas.height)

  // 绘制渐变背景
  const gradient = canvas2d.createLinearGradient(0, 0, canvas.width, canvas.height)
  gradient.addColorStop(0, '#e0f7fa')
  gradient.addColorStop(1, '#f3e5f5')
  canvas2d.fillStyle = gradient
  canvas2d.fillRect(0, 0, canvas.width, canvas.height)

  // 更新和绘制粒子
  particles.forEach((particle, index) => {
    // 更新位置
    particle.x += particle.speedX
    particle.y += particle.speedY

    // 边界检查
    if (particle.x < 0 || particle.x > canvas.width) particle.speedX *= -1
    if (particle.y < 0 || particle.y > canvas.height) particle.speedY *= -1

    // 绘制粒子
    if (canvas2d) {
      canvas2d.fillStyle = particle.color;
      canvas2d.beginPath();
      canvas2d.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2);
      canvas2d.fill();
    }

    // 绘制连线
    particles.slice(index + 1).forEach(otherParticle => {
      const distance = Math.sqrt(
        Math.pow(particle.x - otherParticle.x, 2) +
        Math.pow(particle.y - otherParticle.y, 2)
      )

      if (distance < 100 && canvas2d) {
        canvas2d.strokeStyle = `rgba(150, 150, 255, ${1 - distance / 100})`
        canvas2d.lineWidth = 0.5
        canvas2d.beginPath()
        canvas2d.moveTo(particle.x, particle.y)
        canvas2d.lineTo(otherParticle.x, otherParticle.y)
        canvas2d.stroke()
      }
    })
  })

  animationFrameId = requestAnimationFrame(animateParticles)
}

// 组件销毁时清理资源
onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  const canvas = document.getElementById('particleCanvas')
  if (canvas) {
    document.body.removeChild(canvas)
  }
  window.removeEventListener('resize', resizeCanvas)
})

</script>

<template>
  <div class="ln-login-body">
    <div class="ln-login">
      <div class="ln-login-header">
        <h2 class="ln-login-title">系统登录</h2>
        <p class="ln-login-subtitle">欢迎使用后台管理系统</p>
      </div>

      <!-- 登录表单 -->
      <div class="ln-login-form" v-show="loginForm">
        <a-form>
          <a-form-item name="username" :rules="[{ required: true, message: '请输入用户名!' }]">
            <a-input size="large" placeholder="用户名">
              <template #prefix>
                <UserOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item name="password" :rules="[{ required: true, message: '请输入密码!' }]">
            <a-input-password size="large" placeholder="密码">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 图片验证码 -->
          <a-form-item name="captcha" :rules="[{ required: true, message: '请输入验证码!' }]">
            <a-input size="large" v-model:value="captchaCode" placeholder="验证码">
              <template #prefix>
                <BulbOutlined class="site-form-item-icon" />
              </template>
              <template #addonAfter>
                <div class="captcha-container">
                  <img :src="captchaSrc" alt="验证码" class="captcha-image" @click="refreshCaptcha" title="点击刷新验证码" />
                </div>
              </template>
            </a-input>
          </a-form-item>

          <a-form-item class="login-form-item-href">
            <div class="ln-login-form-item-href">
              <a @click="toRePassword">忘记密码</a>
              <a @click="toRegister">前往注册</a>
            </div>
          </a-form-item>
          <a-form-item class="login-form-btn">
            <a-button type="primary" html-type="submit" class="login-form-button">
              登录
            </a-button>
          </a-form-item>
        </a-form>
        <div class="ln-login-other">
          <div class="ln-login-other-label">第三方登录</div>
          <div class="ln-login-other-icons">
            <i class="iconfont icon-QQ"></i>
            <i class="iconfont icon-weixin"></i>
          </div>
        </div>
      </div>

      <!-- 注册表单 -->
      <div class="ln-login-form" v-show="registerForm">
        <a-form>
          <a-form-item name="username" :rules="[{ required: true, message: '请输入用户名!' }]">
            <a-input size="large" placeholder="用户名">
              <template #prefix>
                <UserOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item name="password" :rules="[{ required: true, message: '请输入密码!' }]">
            <a-input-password size="large" placeholder="密码">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 图片验证码 -->
          <a-form-item name="captcha" :rules="[{ required: true, message: '请输入验证码!' }]">
            <a-input size="large" v-model:value="captchaCode" placeholder="验证码">
              <template #prefix>
                <BulbOutlined class="site-form-item-icon" />
              </template>
              <template #addonAfter>
                <a-button type="primary">发送验证码</a-button>
              </template>
            </a-input>
          </a-form-item>

          <a-form-item class="login-form-item-href">
            <div class="ln-login-form-item-href">
              <a @click="toLogin">前往登录</a>
            </div>
          </a-form-item>
          <a-form-item class="login-form-btn">
            <a-button type="primary" html-type="submit" class="login-form-button">
              注册
            </a-button>
          </a-form-item>
        </a-form>
      </div>

      <!-- 找回密码表单 -->
      <div class="ln-login-form" v-show="rePasswordForm">
        <a-form>
          <a-form-item name="username" :rules="[{ required: true, message: '请输入用户名!' }]">
            <a-input size="large" placeholder="用户名">
              <template #prefix>
                <UserOutlined class="site-form-item-icon" />
              </template>
            </a-input>
          </a-form-item>
          <a-form-item name="password" :rules="[{ required: true, message: '请输入密码!' }]">
            <a-input-password size="large" placeholder="密码">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item name="password" :rules="[{ required: true, message: '请再次输入密码!' }]">
            <a-input-password size="large" placeholder="确认密码">
              <template #prefix>
                <LockOutlined class="site-form-item-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- 图片验证码 -->
          <a-form-item name="captcha" :rules="[{ required: true, message: '请输入验证码!' }]">
            <a-input size="large" v-model:value="captchaCode" placeholder="验证码">
              <template #prefix>
                <BulbOutlined class="site-form-item-icon" />
              </template>
              <template #addonAfter>
                <a-button type="primary">发送验证码</a-button>
              </template>
            </a-input>
          </a-form-item>

          <a-form-item class="login-form-item-href">
            <div class="ln-login-form-item-href">
              <a @click="toLogin">前往登录</a>
            </div>
          </a-form-item>
          <a-form-item class="login-form-btn">
            <a-button type="primary" html-type="submit" class="login-form-button">
              确认
            </a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.ln-login-body {
  width: 100%;
  height: 100vh;
  background: transparent;
  display: flex;
  justify-content: center;
  align-items: center;

  .ln-login {
    width: 500px;
    background-color: #ffffff;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 15px 30px rgba(0, 0, 0, 0.15);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.3);

    .ln-login-header {
      padding: 30px 0 20px;
      text-align: center;
      background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);

      .ln-login-title {
        color: white;
        font-size: 28px;
        font-weight: 600;
        margin: 0 0 10px;
      }

      .ln-login-subtitle {
        color: rgba(255, 255, 255, 0.9);
        font-size: 14px;
        margin: 0;
      }
    }

    .ln-login-form {
      padding: 30px 50px;

      .login-form-item-href {
        margin-bottom: 15px;
      }

      .login-form-btn {
        margin-bottom: 20px;
      }

      .ln-login-form-item-href {
        display: flex;
        align-items: center;
        justify-content: space-between;

        a {
          color: #1890ff;
          font-size: 13px;
          transition: color 0.3s;

          &:hover {
            color: #096dd9;
            text-decoration: underline;
          }
        }
      }

      .login-form-button {
        width: 100%;
        margin: 0;
        height: 42px;
        background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%);
        border: none;
        border-radius: 6px;
        font-size: 16px;
        font-weight: 500;
        color: white;

        &:hover,
        &:focus {
          background: linear-gradient(135deg, #40a9ff 0%, #1890ff 100%);
        }
      }

      .captcha-image {
        width: 100px;
        height: 38px;
        cursor: pointer;
        object-fit: cover;
        border-radius: 4px;
      }

      .captcha-container {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100px;
        height: 38px;
      }

      .ln-login-other {
        display: flex;
        flex-direction: column;
        align-items: center;


        .ln-login-other-label {
          font-size: 13px;
          color: #999;
          width: 100%;
          display: flex;
          align-items: center;
          justify-content: space-between;

          &::before,
          &::after {
            width: 35%;
            height: 1px;
            background-color: #e1e1e1;
            content: "";
            display: inline-flex;
          }
        }

        .ln-login-other-icons {
          margin-top: 15px;

          i {
            margin-right: 15px;
            margin-left: 15px;
            font-size: 32px;
            cursor: pointer;
            color: #1890ff;
            transition: transform 0.3s, color 0.3s;

            &:hover {
              color: #096dd9;
              transform: scale(1.1);
            }
          }
        }
      }
    }
  }
}

.ant-form-item-control-input {
  min-height: 0 !important;
}

.ant-input-affix-wrapper-lg {
  padding: 7px 15px;
  border-radius: 6px;
}

.ant-input-password-affix-wrapper-lg {
  padding: 7px 15px;
  border-radius: 6px;
}

.ant-input-lg {
  border-radius: 6px;
}
</style>
