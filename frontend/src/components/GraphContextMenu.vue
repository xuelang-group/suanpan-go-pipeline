<template>
  <div class="graph-contextmenu" 
    :style="{ left: `${mousePosition.x}px`, top: `${mousePosition.y}px` }">
    <ul class="graph-contextmenu-list">
      <li 
        v-for="menuItem in contextMenus"
        :key="menuItem.key"
        class="graph-contextmenu-item" :class="{'disabled': menuItem.disabled}"
        :style="menuItem.style"
        @mousedown.stop
        @click.stop="menuClick(menuItem)">
        <div class="graph-contextmenu-item-inner" @mouseenter="mouseenterHandler(menuItem)">
          <template v-if="menuItem.type === 'divider'">
            <div class="graph-contextmenu-divider"></div>
          </template>
          <template v-else>
            <div class="graph-contextmenu-label-wrap">
              <span :class="['iconfont', menuItem.icon]"></span>
              <span class="graph-contextmenu-label">{{ menuItem.label }}</span>
            </div>
            <div v-if="menuItem.children && menuItem.children.length" class="graph-contextmenu-more">
              <span class="iconfont icon-node-arrow-right"></span>
            </div>
          </template>
        </div>
        <ul 
          v-if="(!menuItem.disabled) && (menuItem.children) && (menuItem.children.length) && (openedMenu === menuItem.key)" 
          class="graph-contextmenu-list submenu">
          <li 
            v-for="menuSubItem in menuItem.children"
            :key="menuSubItem.key"
            class="graph-contextmenu-item" :class="{'disabled': menuSubItem.disabled}"
            @mousedown.stop
            @click.stop="menuClick(menuSubItem)"
            >
            <div class="graph-contextmenu-item-inner">
              <div class="graph-contextmenu-label-wrap">
                <span :class="['iconfont', menuSubItem.icon]"></span>
                <span class="graph-contextmenu-label">{{ menuSubItem.label }}</span>
              </div>
            </div>
          </li>
        </ul>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'GraphContextMenu',
  props: {
    visible: {
      type: Boolean
    },
    mousePosition: {
      type: [null, Object],
      default() {
        return {x: 0, y: 0}
      }
    },
    contextMenus: {
      type: Array
    }
  },
  emits: ['update:visible', 'menu-click'],
  data() {
    return {
      openedMenu: null
    }
  },
  watch: {
    visible(val) {
      if(!val) {
        this.openedMenu = null
      }
    }
  },
  mounted() {
    window.addEventListener('mousedown', this.mousedownHandler)
  },
  beforeUnmount() {
    window.removeEventListener('mousedown', this.mousedownHandler)
  },
  methods: {
    mousedownHandler() {
      this.$emit('update:visible', false)
    },
    menuClick(menuItem) {
      if(menuItem.disabled) {
        return
      }
      if(menuItem.children && menuItem.children.length) {
        return
      }
      this.$emit('update:visible', false)
      this.$emit('menu-click', menuItem)
    },
    mouseenterHandler(menuItem) {
      this.openedMenu = menuItem.key
    }
  }
}
</script>

<style>

</style>